package controllers

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nebisin/gopress/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Handler struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (handler *Handler) Initialize() {
	getEnv()
	handler.initializeDatabase()
	handler.initializeRoutes()
}

func getEnv() {
	log.Println("We are getting the env values...")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}
}

func (handler *Handler) initializeDatabase() {
	log.Println("We are initializing the database...")

	var err error

	handler.DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	} else {
		log.Println("🌍 Database connection is successful")
	}

	// Migrate the schema
	if err := handler.DB.AutoMigrate(&models.Post{}, &models.User{}); err != nil {
		log.Fatalf("Error auto migration: %v", err)
	}
}

func (handler *Handler) Run(addr string) {
	log.Println("🚀 Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, handler.Router))
}
