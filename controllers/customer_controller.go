package controllers

import (
	"bsnack/database/models"
	"bsnack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	Service *services.CustomerService
}

func NewCustomerController(service *services.CustomerService) *CustomerController {
	return &CustomerController{Service: service}
}

func (c *CustomerController) GetAll(ctx *gin.Context) {
	data, err := c.Service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *CustomerController) Create(ctx *gin.Context) {
	var cst models.Customer

	if err := ctx.ShouldBindJSON(&cst); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Service.CreateCustomer(&cst); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cst)
}

type ExchangePointRequest struct {
	CustomerID string `json:"customer_id"`
	Point      int    `json:"point"`
}

func (c *CustomerController) ExchangePoint(ctx *gin.Context) {
	var payload ExchangePointRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Service.ExchangePoint(payload.CustomerID, payload.Point); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "exchange point success"})
}
