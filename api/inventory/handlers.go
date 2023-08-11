package inventory

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *RoutesWrapper) HandleAllProducts(c *gin.Context) {
	resp, err := r.InventoryService.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

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
