package main

import (
	"log"

	"github.com/Sasank-V/CIMP-Golang-Backend/controllers"
	"github.com/Sasank-V/CIMP-Golang-Backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading env: ", err)
	}
	log.Printf("ENV Loaded")

	controllers.ConnectClubCollection()
	controllers.ConnectDepartmentCollection()
	controllers.ConnectContributionCollection()
	controllers.ConnectUserCollection()

	authApi := r.Group("/api/auth")
	userApi := r.Group("/api/user")
	contApi := r.Group("/api/contribution")

	routes.SetupUserRoutes(userApi)
	routes.SetupAuthRoutes(authApi)
	routes.SetupContributionRoutes(contApi)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server", err)
	}
}
