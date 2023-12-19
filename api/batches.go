package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/billing_system/db/sqlc"
)

type createBatchRequest struct {
	Name             string    `json:"name" binding:"required,min=1"`
	ActivationStatus string    `json:"activation_status" binding:"required,oneof=active inactive suspended canceled"`
	CustomerName     string    `json:"customer_name" binding:"required,min=1"`
	NoOfDevices      int64     `json:"no_of_devices" binding:"required,min=1"`
	DeliveryDate     time.Time `json:"delivery_date"`
	WarrantyEnd      time.Time `json:"warranty_end"`
}

// CreateBatch godoc
//
//	@Summary		Create new Batch
//	@Description	Create a new Batch specifying its name
//	@Tags			batches
//	@Produce		json
//	@Accept			json
//	@Param			arg	body		createBatchRequest	true	"Create Batch body"
//
//	@Success		200	{object}	db.Batch
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/batches [post]
func (server *Server) createBatch(c *gin.Context) {
	var req createBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var customer db.Customer
	var err error
	customer, err = server.store.GetCustomerID(c, req.CustomerName)
	if err != nil {
		if err == sql.ErrNoRows {
			customer, err = server.store.CreateCustomer(c, req.CustomerName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	arg := db.CreateBatchParams{
		Name:             req.Name,
		ActivationStatus: req.ActivationStatus,
		CustomerID:       customer.ID,
		NoOfDevices:      int32(req.NoOfDevices),
		DeliveryDate:     sql.NullTime{Time: req.DeliveryDate, Valid: true},
		WarrantyEnd:      sql.NullTime{Time: req.WarrantyEnd, Valid: true},
	}
	category, err := server.store.CreateBatch(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, category)

}

type deleteBatchRequest struct {
	BatchID int64 `uri:"batch_id" binding:"required,min=1"`
}

// DeleteBatch godoc
//
//	@Summary		Delete Batch
//	@Description	Delete batch by a batch ID
//	@Tags			batches
//
//
//	@Produce		plain
//	@Param			batch_id	path		string	true	"Batch ID"
//
//	@Success		200			true		bool
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Router			/batches/{batch_id} [delete]
func (server *Server) deleteBatch(c *gin.Context) {
	var req deleteBatchRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetBatchForUpdate(c, req.BatchID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteBatch(c, req.BatchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, true)

}

type listBatchesRequest struct {
	CustomerName string `form:"customer_name"`
	PageID       int32  `form:"page_id" binding:"required,min=1"`
	PageSize     int32  `form:"page_size" binding:"required,min=5,max=10"`
}

// ListBatches godoc
//
//	@Summary		List Batches
//	@Description	List all batches with optional filter "customer_name", pagination params are required
//
//
//	@Tags			batches
//
//
//	@Produce		json
//	@Param			customer_name	query		string	false	"Filter: customer name"
//	@Param			page_id			query		int		true	"Page ID"
//	@Param			page_size		query		int		true	"Page Size"
//
//	@Success		200				{array}		db.Batch
//	@Failure		400				{object}	error
//	@Failure		500				{object}	error
//
//	@Router			/batches [get]
func (server *Server) listBatches(c *gin.Context) {
	var req listBatchesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var customerBool bool
	var customer db.Customer
	var err error
	if req.CustomerName != "" {
		customer, err = server.store.GetCustomerID(c, req.CustomerName)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		customerBool = true
	}
	arg := db.ListAllBatchesParams{
		CustomerID: sql.NullInt64{Valid: customerBool, Int64: customer.ID},
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}
	batches, err := server.store.ListAllBatches(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, batches)

}

type updateBatchRequestJSON struct {
	CustomerName     string       `json:"customer_name"`
	ActivationStatus string       `json:"activation_status"`
	NoOfDevices      int32        `json:"no_of_devices"`
	DeliveryDate     sql.NullTime `json:"delivery_date"`
	WarrantyEnd      sql.NullTime `json:"warranty_end"`
}
type updateBatchRequestURI struct {
	BatchID int64 `uri:"batch_id" binding:"required,min=1"`
}

// UpdateBatch godoc
//
//	@Summary		Update Batch
//	@Description	Update Batch by a Batch ID
//	@Tags			batches
//
//
//	@Produce		json
//
//	@Accept			json
//
//	@Param			arg			body		updateBatchRequestJSON	true	"Update Batch body"
//	@Param			batch_id	path		int						true	"Batch ID for update"
//
//	@Success		200			{object}	db.Batch
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Router			/batches/{batch_id} [put]
func (server *Server) updateBatch(c *gin.Context) {
	var req updateBatchRequestJSON
	var reqURI updateBatchRequestURI
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var customer db.Customer
	var err error
	batch, err := server.store.GetBatchForUpdate(c, reqURI.BatchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.CustomerName != "" {
		customer, err = server.store.GetCustomerID(c, req.CustomerName)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		batch.CustomerID = customer.ID
	}

	arg := db.UpdateBatchParams{
		ID:               reqURI.BatchID,
		ActivationStatus: req.ActivationStatus,
		WarrantyEnd:      sql.NullTime{Time: req.WarrantyEnd.Time.Round(time.Second), Valid: req.WarrantyEnd.Valid},
		CustomerID:       batch.CustomerID,
		NoOfDevices:      req.NoOfDevices,
		DeliveryDate:     sql.NullTime{Time: req.DeliveryDate.Time.Round(time.Second), Valid: req.DeliveryDate.Valid},
	}

	batch, err = server.store.UpdateBatch(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, batch)

}
