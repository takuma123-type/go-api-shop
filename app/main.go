package main

import (
	"go-test/controllers"
	"go-test/models"
	"go-test/repositories"
	"go-test/services"

	"github.com/gin-gonic/gin"
)

func main() {
	items := []models.Item{
		{ID: 1, Name: "Laptop", Price: 1000, Description: "A laptop", SoldOut: false},
		{ID: 2, Name: "Mouse", Price: 10, Description: "A mouse", SoldOut: false},
		{ID: 3, Name: "Keyboard", Price: 20, Description: "A keyboard", SoldOut: false},
		{ID: 4, Name: "Monitor", Price: 200, Description: "A monitor", SoldOut: false},
		{ID: 5, Name: "Headset", Price: 50, Description: "A headset", SoldOut: false},
	}

	itemRepository := repositories.NewItemMemoryRepository(items)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	r := gin.Default()
	r.GET("/items", itemController.FindAll)
	r.GET("items/:id", itemController.FindById)
	r.POST("/items", itemController.Create)
	r.PUT("/items/:id", itemController.Update)
	r.DELETE("/items/:id", itemController.Delete)
	r.Run(":8080")
}
