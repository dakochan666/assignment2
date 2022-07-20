package router

import (
	"assignment2/controllers"
	"assignment2/database"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()

	db := database.ConnectDB()
	// itemController := controllers.NewItemController(db)
	orderController := controllers.NewOrderController(db)

	// itemGroup := router.Group("/items")
	// {
	// 	itemGroup.POST("/", itemController.CreateItem)
	// }

	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("/", orderController.CreateOrder)
		orderGroup.GET("/", orderController.GetOrders)
		orderGroup.PUT("/:orderId", orderController.UpdateOrder)
		orderGroup.DELETE("/:orderId", orderController.DeleteOrder)
	}

	return router
}
