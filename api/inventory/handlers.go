package inventory

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Returns all Products
// @Description Returns all Products present in the Inventory
// @Schemes http
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 200 {object} GetAllProductsResp
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/inventory/all-products [get]
func (r *RoutesWrapper) HandleAllProducts(c *gin.Context) {
	resp, err := r.InventoryService.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

// @Summary Add a Product
// @Description Add a Product to the Inventory
// @Schemes http
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param data body AddProductReq true "AddProductRequest"
// @Success 200 {string} SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/inventory/add-product [post]
func (r *RoutesWrapper) HandleAddProduct(c *gin.Context) {
	var input AddProductReq

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	err := r.InventoryService.AddProduct(c.Request.Context(), &input)

	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, "a new product added!")
	return
}

// @Summary Update a Product
// @Description Update a Product to the Inventory
// @Schemes http
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param data body UpdateProductReq true "UpdateProductRequest"
// @Success 200 {object} ProductResp
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/inventory/update-product/{id} [post]
func (r *RoutesWrapper) HandleUpdateProduct(c *gin.Context) {

	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input UpdateProductReq

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	resp, err := r.InventoryService.UpdateProduct(c.Request.Context(), int(pathParams.ID), &input)

	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

// @Summary Delete a Product
// @Description Delete a Product to the Inventory
// @Schemes http
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {string} SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/inventory/delete-product/{id} [post]
func (r *RoutesWrapper) HandleDeleteProduct(c *gin.Context) {
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.InventoryService.DeleteProduct(c.Request.Context(), int(pathParams.ID))

	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, "successfully deleted the product")
	return
}
