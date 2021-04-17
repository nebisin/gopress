package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nebisin/gopress/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Handler struct {
	Router *mux.Router
	DB *gorm.DB
}

func (handler *Handler) Initialize() {

	var err error

	handler.DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	handler.DB.AutoMigrate(&models.Post{}, &models.User{})

	handler.Router = mux.NewRouter()

	handler.initializeRoutes()
}

func (handler *Handler) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, handler.Router))
}