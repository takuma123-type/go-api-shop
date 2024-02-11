package main

import (
	"go-test/controllers"
	"go-test/infra"
	"go-test/middlewares"
	"go-test/repositories"
	"go-test/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	r := gin.Default()
	r.Use(cors.Default())
	itemRouter := r.Group("/items")
	itemsRouterWithAuth := r.Group("/items", middlewares.AuthMiddleware(authService))
	authRouter := r.Group("/auth")

	itemRouter.GET("", itemController.FindAll)
	itemsRouterWithAuth.GET("/:id", itemController.FindById)
	itemsRouterWithAuth.POST("", itemController.Create)
	itemsRouterWithAuth.PUT("/:id", itemController.Update)
	itemsRouterWithAuth.DELETE("/:id", itemController.Delete)

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)

	return r
}

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	r := setupRouter(db)

	// itemRepository := repositories.NewItemRepository(db)
	// itemService := services.NewItemService(itemRepository)
	// itemController := controllers.NewItemController(itemService)

	// authRepository := repositories.NewAuthRepository(db)
	// authService := services.NewAuthService(authRepository)
	// authController := controllers.NewAuthController(authService)

	// r := gin.Default()
	// r.Use(cors.Default())
	// itemRouter := r.Group("/items")
	// itemsRouterWithAuth := r.Group("/items", middlewares.AuthMiddleware(authService))
	// authRouter := r.Group("/auth")

	// itemRouter.GET("", itemController.FindAll)
	// itemsRouterWithAuth.GET("/:id", itemController.FindById)
	// itemsRouterWithAuth.POST("", itemController.Create)
	// itemsRouterWithAuth.PUT("/:id", itemController.Update)
	// itemsRouterWithAuth.DELETE("/:id", itemController.Delete)

	// authRouter.POST("/signup", authController.Signup)
	// authRouter.POST("/login", authController.Login)

	r.Run(":8080")
}
