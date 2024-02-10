package main

import (
	"go-test/infra"
	"go-test/models"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.Item{}); err != nil {
		panic("failed to migrate")
	}
}
