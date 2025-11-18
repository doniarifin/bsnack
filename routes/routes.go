package routes

import (
	"bsnack/config"
	"bsnack/controllers"
	"bsnack/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	redisClient := config.ConnectRedis()

	// apply middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//product
	productService := services.NewProductService(db, redisClient)
	productController := controllers.NewProductController(productService)
	//customer
	customerService := services.NewCustomerService(db, redisClient, productService)
	customerController := controllers.NewCustomerController(customerService)
	//transaction
	trxSrv := services.NewTransactionService(db, redisClient, productService, customerService)
	trxController := controllers.NewTransactionController(trxSrv)

	// api group
	api := r.Group("/api/v1")
	// api.Use(middleware.JWTMiddleware())

	//customer
	api.GET("/customer", customerController.GetAll)
	api.POST("/customer", customerController.Create)
	api.POST("/customer/exchangepoint", customerController.ExchangePoint)

	//transaction
	api.GET("/transaction", trxController.GetAll)
	api.GET("/transaction/:id", trxController.GetByID)
	api.POST("/transaction", trxController.Create)
	api.PUT("/transaction/:id", trxController.Update)
	// api.DELETE("/transaction/:id", productController.Delete)

	//product
	api.GET("/product", productController.GetAll)
	api.POST("/product/getbydate", productController.GetProductByDate)
	api.POST("/product", productController.Create)
	api.PUT("/product/:id", productController.Update)
	// api.DELETE("/product/:id", productController.Delete)
}
