package main

import (
	"fmt"
	"log"

	"github.com/Sasank-V/CIMP-Golang-Backend/routes"
	"github.com/Sasank-V/CIMP-Golang-Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading env: ", err)
	}
	fmt.Printf("ENV Loaded\n")
	utils.InitDB()

	userApi := r.Group("/api/user")
	routes.SetupUserRoutes(userApi)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server", err)
	}
}
