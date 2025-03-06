package main

import (
	"database/sql"
	"gin-api/database"
	"gin-api/models"
	"gin-api/routes"
	"gin-api/util"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var settings models.Settings

func main() {
	log.Println("Starting the application...")
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Printf("Failed to open database: %v", err)
	}
	defer db.Close()

	database.SetDB(db)
	log.Println("Database connection established.")
	


	r := gin.New()

	// CORS middleware setup before routes
	config := cors.DefaultConfig()
    config.AllowOrigins = []string{
        "http://localhost:5173", 
    }
    config.AllowHeaders = []string{
        "Origin",
        "Content-Type",
        "Accept",
        "Authorization",
    }
    config.ExposeHeaders = []string{"Content-Length"}
    config.AllowCredentials = true

    r.Use(cors.New(config))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(util.GetRateLimiter())

	routes.SetupRouter(r)


	r.Run(":3000")
}
