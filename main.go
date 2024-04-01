package main

import (
	"log"
	"os"

	"github.com/Bec4use/ecommerce-golang/controllers"
	"github.com/Bec4use/ecommerce-golang/database"
	"github.com/Bec4use/ecommerce-golang/middleware"
	_ "github.com/Bec4use/ecommerce-golang/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	router.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
