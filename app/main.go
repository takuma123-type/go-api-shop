package main

import (
	"go-test/controllers"
	"go-test/infra"
	"go-test/repositories"

	// "go-test/models"

	"go-test/services"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	// items := []models.Item{
	// 	{ID: 1, Name: "Laptop", Price: 1000, Description: "A laptop", SoldOut: false},
	// 	{ID: 2, Name: "Mouse", Price: 10, Description: "A mouse", SoldOut: false},
	// 	{ID: 3, Name: "Keyboard", Price: 20, Description: "A keyboard", SoldOut: false},
	// 	{ID: 4, Name: "Monitor", Price: 200, Description: "A monitor", SoldOut: false},
	// 	{ID: 5, Name: "Headset", Price: 50, Description: "A headset", SoldOut: false},
	// }

	// itemRepository := repositories.NewItemMemoryRepository(items)
	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	r := gin.Default()
	itemRouter := r.Group("/items")
	authRouter := r.Group("/auth")

	itemRouter.GET("", itemController.FindAll)
	itemRouter.GET("/:id", itemController.FindById)
	itemRouter.POST("", itemController.Create)
	itemRouter.PUT("/:id", itemController.Update)
	itemRouter.DELETE("/:id", itemController.Delete)

	authRouter.POST("/signup", authController.Signup)

	r.Run(":8080")
}
