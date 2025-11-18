package controllers

import (
	"bsnack/database/models"
	"bsnack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	Service *services.ProductService
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{Service: service}
}

type GetProductRequest struct {
	CreatedAt string `json:"created_at"`
}

func (c *ProductController) GetProductByDate(ctx *gin.Context) {
	var payload GetProductRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.Service.GetByProductDate(payload.CreatedAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *ProductController) Create(ctx *gin.Context) {
	var p models.Product

	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Service.CreateProduct(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, p)
}

func (c *ProductController) Update(ctx *gin.Context) {
	var p models.Product
	idParam := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p.ID = idParam
	if err := c.Service.UpdateProduct(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, p)
}
