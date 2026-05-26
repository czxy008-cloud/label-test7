package handler

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	api := r.Group("/api/v1")
	{
		productHandler := NewProductHandler()
		products := api.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.GetByID)
		}

		cartHandler := NewCartHandler()
		carts := api.Group("/cart")
		{
			carts.GET("", cartHandler.List)
			carts.POST("", cartHandler.Add)
			carts.PUT("/:product_id", cartHandler.Update)
			carts.DELETE("/:product_id", cartHandler.Remove)
			carts.DELETE("", cartHandler.Clear)
		}

		groupBuyHandler := NewGroupBuyHandler()
		groupBuys := api.Group("/group-buys")
		{
			groupBuys.GET("", groupBuyHandler.List)
			groupBuys.POST("", groupBuyHandler.Create)
			groupBuys.POST("/join", groupBuyHandler.Join)
			groupBuys.GET("/:code", groupBuyHandler.GetByCode)
		}

		orderHandler := NewOrderHandler()
		orders := api.Group("/orders")
		{
			orders.GET("", orderHandler.List)
			orders.GET("/:id", orderHandler.GetByID)
			orders.POST("/:id/pay", orderHandler.Pay)
			orders.POST("/:id/ship", orderHandler.Ship)
		}
	}

	return r
}
