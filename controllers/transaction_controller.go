package controllers

import (
	"bsnack/database/models"
	"bsnack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	Service *services.TransactionService
}

func NewTransactionController(service *services.TransactionService) *TransactionController {
	return &TransactionController{Service: service}
}

func (c *TransactionController) GetAll(ctx *gin.Context) {
	data, err := c.Service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *TransactionController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	data, err := c.Service.GetByID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *TransactionController) Create(ctx *gin.Context) {
	var t models.Transaction

	if err := ctx.ShouldBindJSON(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Service.CreateTransaction(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, t)
}

func (c *TransactionController) Update(ctx *gin.Context) {
	var t models.Transaction
	idParam := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t.ID = idParam
	if err := c.Service.UpdateTransaction(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, t)
}
