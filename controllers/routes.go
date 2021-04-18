package controllers

import "github.com/nebisin/gopress/middlewares"

func (handler *Handler) initializeRoutes() {
	handler.Router.Use(middlewares.SetMiddlewareJSON)
	handler.Router.HandleFunc("/posts/{id}", handler.GetPostById).Methods("GET")
	handler.Router.HandleFunc("/posts", middlewares.SetMiddlewareAuthentication(handler.CreatePost)).Methods("POST")
	handler.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(handler.UpdatePost)).Methods("PUT")
	handler.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(handler.DeletePost)).Methods("DELETE")
	handler.Router.HandleFunc("/posts", handler.GetPosts).Methods("GET")

	handler.Router.HandleFunc("/register", handler.Register).Methods("POST")
	handler.Router.HandleFunc("/login", handler.Login).Methods("POST")
	handler.Router.HandleFunc("/users/{id}", handler.GetUserById).Methods("GET")
}
