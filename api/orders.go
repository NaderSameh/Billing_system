package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/naderSameh/billing_system/worker"
)

type createOrderRequest struct {
	BatchName string `json:"batch_name" binding:"required,min=1"`
}

// CreateOrder godoc
//
//	@Summary		Create new Order
//	@Description	Create a new order specifying its batch name, it will create with no NRC, default MRC for 1 year
//	@Tags			orders
//	@Produce		json
//	@Accept			json
//	@Param			arg	body		createOrderRequest	true	"Create Order body"
//
//	@Success		200	{object}	db.Order
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/orders [post]
func (server *Server) createOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var batch db.Batch
	var err error
	batch, err = server.store.GetBatchByName(c, req.BatchName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return

		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	bundles, err := server.store.ListBundlesByCustomerID(c, batch.CustomerID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}
	arg := db.CreateOrderParams{
		StartDate: time.Now().Round(time.Second),
		EndDate:   time.Now().AddDate(1, 0, 0).Round(time.Second),
		BatchID:   batch.ID,
	}
	if len(bundles) == 0 {
		arg.BundleID = 1
	} else {
		arg.BundleID = bundles[0].ID
	}
	order, err := server.store.CreateOrder(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	bundle, err := server.store.GetBundleByID(c, order.BundleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg2 := db.CreatePaymentParams{
		Confirmed:  false,
		OrderID:    order.ID,
		DueDate:    time.Now().AddDate(0, 1, 0),
		Payment:    bundle.Mrc * float64(batch.NoOfDevices),
		CustomerID: batch.CustomerID,
	}

	_, err = server.store.CreatePayment(c, arg2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	charges, err := server.store.AddToDue(c, db.AddToDueParams{Amount: arg2.Payment, ID: batch.CustomerID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	taskPayload := &worker.PayloadSendEmail{
		OrderID:      order.ID,
		BatchName:    batch.Name,
		BatchID:      arg.BatchID,
		CustomerName: charges.Customer,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err = taskDistributor.NewEmailDeliveryTask(taskPayload, opts...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, order)

}

type updateOrderRequestJSON struct {
	EndDate   time.Time       `json:"end_date" binding:"required"`
	StartDate sql.NullTime    `json:"start_date"`
	BundleID  int64           `json:"bundle_id" binding:"required,min=1"`
	Nrc       sql.NullFloat64 `json:"nrc"`
}

type updateOrderRequestPATH struct {
	OrderID int64 `uri:"order_id" binding:"required,min=1"`
}

// UpdateOrder godoc
//
//	@Summary		Update Order with actual params
//	@Description	Update an order specifying its end date, bundle mrc, nrc flag
//	@Tags			orders
//	@Produce		json
//	@Accept			json
//	@Param			arg			body		updateOrderRequestJSON	true	"Create Order body"
//	@Param			order_id	path		int						true	"order ID for update"
//
//	@Success		200			{object}	db.Order
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Router			/orders/{order_id} [put]
func (server *Server) updateOrder(c *gin.Context) {
	var req updateOrderRequestJSON
	var reqPath updateOrderRequestPATH

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.ShouldBindUri(&reqPath); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateOrdersParams{
		ID:        reqPath.OrderID,
		Nrc:       req.Nrc,
		BundleID:  req.BundleID,
		EndDate:   req.EndDate,
		StartDate: req.StartDate,
	}
	order, err := server.store.UpdateOrders(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return

		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	bundle, err := server.store.GetBundleByID(c, req.BundleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	batch, err := server.store.GetBatchForUpdate(c, order.BatchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg2 := db.CreatePaymentParams{
		Confirmed:  false,
		OrderID:    order.ID,
		DueDate:    time.Now().AddDate(0, 1, 0),
		Payment:    bundle.Mrc * float64(batch.NoOfDevices),
		CustomerID: batch.CustomerID,
	}

	_, err = server.store.CreatePayment(c, arg2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.AddToDue(c, db.AddToDueParams{Amount: arg2.Payment, ID: batch.CustomerID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if req.Nrc.Valid == true {
		arg2 := db.CreatePaymentParams{
			Confirmed:  false,
			OrderID:    order.ID,
			DueDate:    time.Now().AddDate(0, 1, 0),
			Payment:    req.Nrc.Float64,
			CustomerID: batch.CustomerID,
		}
		_, err = server.store.CreatePayment(c, arg2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		arg3 := db.AddToDueParams{Amount: req.Nrc.Float64, ID: batch.CustomerID}
		_, err = server.store.AddToDue(c, arg3)
		if err != nil {

			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

	}
	c.JSON(http.StatusOK, order)

}

// GetOrder godoc
//
//	@Summary		Get all orders
//	@Description	Get all placed order with details
//	@Tags			orders
//	@Produce		json
//
//	@Success		200	{object}	db.Order
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/orders [get]
func (server *Server) getOrder(c *gin.Context) {

	orders, err := server.store.ListAllOrders(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, orders)

}
