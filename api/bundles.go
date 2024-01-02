package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/billing_system/db/sqlc"
)

type createBundleRequest struct {
	MRC         float64 `json:"mrc" binding:"required,min=1"`
	Description string  `json:"description" binding:"required"`
}

// CreateBundle godoc
//
//	@Summary		Create new Bundle
//	@Description	Create a new Bundle specifying its name
//	@Tags			bundles
//	@Produce		json
//	@Accept			json
//	@Param			arg	body		createBundleRequest	true	"Create bundle body"
//
//	@Success		200	{object}	db.Bundle
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/bundles [post]
func (server *Server) createBundle(c *gin.Context) {
	var req createBundleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateBundleParams{
		Mrc:         req.MRC,
		Description: req.Description,
	}
	category, err := server.store.CreateBundle(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, category)

}

type assignBundleRequest struct {
	AssignedCusomers json.RawMessage `json:"assigned_customers" binding:"required"`
	BundleID         int64           `json:"bundle_id" binding:"required,min=1"`
}

// AssignBundle godoc
//
//	@Summary		Assign bundle to customer
//	@Description	Assign a bundle to a specific customer using bundle id and customer name
//	@Tags			bundles
//	@Produce		plain
//	@Accept			json
//	@Param			arg	body		assignBundleRequest	true	"Assign bundle body"
//
//	@Success		200	bool		true
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/bundles/assign [post]
func (server *Server) assignBundleToCustomer(c *gin.Context) {
	var req assignBundleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteOldBundleCustomers(c, req.BundleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.InsertNewBundleCustomersParams{
		Column1:   req.AssignedCusomers,
		BundlesID: req.BundleID,
	}
	err = server.store.InsertNewBundleCustomers(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, true)

}

type deleteBundleRequest struct {
	BundleID int64 `uri:"bundle_id" binding:"required,min=1"`
}

// DeleteBundle godoc
//
//	@Summary		Delete  Bundle
//	@Description	Delete  Bundle
//	@Tags			bundles
//	@Produce		json
//	@Param			bundle_id	path		int	true	"Bundle ID"
//
//	@Success		200			true		bool
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Router			/bundles/{bundle_id} [delete]
func (server *Server) deleteBundle(c *gin.Context) {
	var req deleteBundleRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetBundleByID(c, req.BundleID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteBundle(c, req.BundleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, true)

}

type getBundlesRequest struct {
	CustomerName string `form:"customer_name"`
}

type getBundleResponse struct {
	ID          int64   `json:"id"`
	MRC         float64 `json:"mrc"`
	Description string  `json:"description"`
}

// GetBundles godoc
//
//	@Summary		Get Bundles
//	@Description	Get bundles for a specific customer - get all system bundles if no customer specified
//	@Tags			bundles
//	@Produce		json
//	@Accept			json
//	@Param			customer_name	query		string	false	"Get bundle body"
//
//	@Success		200				{array}		db.ListBundlesWithCustomerRow
//	@Failure		400				{object}	error
//	@Failure		500				{object}	error
//	@Router			/bundles [get]
func (server *Server) getBundles(c *gin.Context) {
	var req getBundlesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var bundles []db.Bundle
	if req.CustomerName != "" {
		customer, err := server.store.GetCustomerID(c, req.CustomerName)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		bundles, err = server.store.ListBundlesByCustomerID(c, customer.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {

		bundles, err := server.store.ListBundlesWithCustomer(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, bundles)
		return
	}

	c.JSON(http.StatusOK, bundles)

}
