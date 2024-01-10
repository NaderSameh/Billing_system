package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/billing_system/db/sqlc"
)

type updatePaymentLogRequestJSON struct {
	DueDate   sql.NullTime `json:"due_date"`
	Confirmed *bool        `json:"confirmed" binding:"required"`
}
type updatePaymentLogRequestURI struct {
	LogID int64 `uri:"log_id" binding:"required,min=1"`
}

// UpdatePaymentLog godoc
//
//	@Summary		Update Payment
//	@Description	Update Payment log's due date & confirmation using the payment log ID
//	@Tags			payments logs
//
//
//	@Produce		json
//
//	@Accept			json
//
//	@Param			arg		body		updatePaymentLogRequestJSON	true	"Update log body"
//	@Param			log_id	path		int							true	"payment log ID for update"
//
//	@Success		200		{object}	db.PaymentLog
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/payments_logs/{log_id} [put]
func (server *Server) updatePaymentLog(c *gin.Context) {
	var req updatePaymentLogRequestJSON
	var reqURI updatePaymentLogRequestURI
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payment, err := server.store.GetPaymentForUpdate(c, reqURI.LogID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	arg := db.UpdatePaymentParams{
		ID:        reqURI.LogID,
		DueDate:   sql.NullTime{Time: req.DueDate.Time.Round(time.Second), Valid: req.DueDate.Valid},
		Confirmed: *req.Confirmed,
	}
	if payment.Confirmed == false && *req.Confirmed == true {
		arg.ConfirmationDate = sql.NullTime{Time: time.Now().Round(time.Second), Valid: true}
		_, err = server.store.AddToPaid(c, db.AddToPaidParams{Amount: payment.Payment, ID: payment.CustomerID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else if *req.Confirmed == true {
		arg.ConfirmationDate = sql.NullTime{Time: time.Now().Round(time.Second), Valid: true}
	}
	payment, err = server.store.UpdatePayment(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, payment)
}

type deletePaymentRequest struct {
	PaymentID int64 `uri:"log_id" binding:"required,min=1"`
}

// DeletePaymentLog godoc
//
//	@Summary		Delete Payment Log
//	@Description	Delete a payment log, removing its corresponding charges from the customer total charges
//	@Tags			payments logs
//
//
//	@Produce		plain
//	@Param			log_id	path		string	true	"Log ID"
//
//	@Success		200		true		bool
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/payments_logs/{log_id} [delete]
func (server *Server) deletePayment(c *gin.Context) {
	var req deletePaymentRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payment, err := server.store.GetPaymentForUpdate(c, req.PaymentID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeletePayment(c, req.PaymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.AddToDue(c, db.AddToDueParams{Amount: -float64(payment.Payment), ID: payment.CustomerID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, true)

}

type listPaymentLogsQuery struct {
	Confirmed    *bool   `form:"confirmed"`
	CustomerName *string `form:"customer_name"`
	PageID       int32   `form:"page_id" binding:"required,min=1"`
	PageSize     int32   `form:"page_size" binding:"required,min=5,max=10"`
}

type listPaymentsResponse struct {
	Payments []db.PaymentLog `json:"payments"`
	Pages    int32           `json:"pages"`
}

// ListPayments godoc
//
//	@Summary		List Payments
//	@Description	List payments, filtering by confirmation & customer id are optional params, pagination params is required
//	@Tags			payments logs
//
//
//	@Produce		json
//	@Param			page_id			query		int		true	"Page ID"
//	@Param			page_size		query		int		true	"Page Size"
//	@Param			customer_name	query		string	false	"customer_name"
//	@Success		200				{array}		db.PaymentLog
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Router			/payments_logs [get]
func (server *Server) listPaymentLogs(c *gin.Context) {

	var req listPaymentLogsQuery
	var res listPaymentsResponse
	var err error
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var confirmed sql.NullBool
	if req.Confirmed != nil {
		confirmed = sql.NullBool{Bool: *req.Confirmed, Valid: true}
	} else {
		confirmed = sql.NullBool{Bool: false, Valid: false}
	}

	var customer_ID sql.NullInt64
	if req.CustomerName != nil {
		customer, _ := server.store.GetCustomerID(c, *req.CustomerName)
		customer_ID = sql.NullInt64{Int64: customer.ID, Valid: true}
	} else {
		customer_ID = sql.NullInt64{Int64: 0, Valid: false}
	}
	arg := db.ListPaymentsParams{
		Confirmed:  confirmed,
		CustomerID: customer_ID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}
	res.Payments, err = server.store.ListPayments(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg2 := db.ListAllPaymentsCountParams{
		Confirmed:  confirmed,
		CustomerID: customer_ID,
	}
	count, err := server.store.ListAllPaymentsCount(c, arg2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res.Pages = int32(count) / req.PageSize
	if int32(count)%req.PageSize != 0 {
		res.Pages++
	}

	c.JSON(http.StatusOK, res)

}
