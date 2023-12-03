package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/billing_system/db/sqlc"
)

type listChargesRequest struct {
	CustomerName string `form:"customer_name"`
}

// ListCharges godoc
//
//	@Summary		List Charges
//	@Description	List all charges on a customer (optional filter), list all charges in the system
//
//
//	@Tags			customer_charges
//
//
//	@Produce		json
//	@Param			customer_name	query		string	false	"Filter: customer name"
//
//	@Success		200				{array}		db.Customer
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//
//	@Router			/charges [get]
func (server *Server) listChargesPerCustomer(c *gin.Context) {
	var req listChargesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var customers []db.Customer
	var err error
	if req.CustomerName != "" {
		customers, err = server.store.ListAllCharges(c, sql.NullString{String: req.CustomerName, Valid: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

	} else {

		customers, err = server.store.ListAllCharges(c, sql.NullString{String: "", Valid: false})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	c.JSON(http.StatusOK, customers)
}
