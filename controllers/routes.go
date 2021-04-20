package controllers

import (
	"github.com/gorilla/mux"
	"github.com/nebisin/gopress/middlewares"
	"log"
)

func (handler *Handler) initializeRoutes() {
	log.Println("We are initializing the routes...")

	handler.Router = mux.NewRouter()

	handler.Router.Use(middlewares.SetLoggingMiddleware)
	handler.Router.Use(middlewares.SetMiddlewareJSON)

	handler.Router.HandleFunc("/posts/{id}", handler.handlePostGet).Methods("GET")
	handler.Router.HandleFunc("/posts", middlewares.SetMiddlewareAuthentication(handler.handlePostCreate)).Methods("POST")
	handler.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(handler.handlePostUpdate)).Methods("PUT")
	handler.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(handler.handlePostDelete)).Methods("DELETE")
	handler.Router.HandleFunc("/posts", handler.handlePostGetMany).Methods("GET")

	handler.Router.HandleFunc("/register", handler.handleAuthRegister).Methods("POST")
	handler.Router.HandleFunc("/login", handler.handleAuthLogin).Methods("POST")
	handler.Router.HandleFunc("/me", middlewares.SetMiddlewareAuthentication(handler.handleMe)).Methods("GET")
	handler.Router.HandleFunc("/me/posts", middlewares.SetMiddlewareAuthentication(handler.handleMyPosts)).Methods("GET")

	handler.Router.HandleFunc("/users/{id}", handler.handleUserGet).Methods("GET")
	handler.Router.HandleFunc("/users/{id}/posts", handler.handleUserPostsGet).Methods("GET")

}
