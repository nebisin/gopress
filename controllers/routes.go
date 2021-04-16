package controllers

import "github.com/nebisin/gopress/middlewares"

func (handler *Handler) initializeRoutes() {
	handler.Router.Use(middlewares.SetMiddlewareJSON)
	handler.Router.HandleFunc("/posts/{id}", handler.GetPostById).Methods("GET")
	handler.Router.HandleFunc("/posts", handler.CreatePost).Methods("POST")
	handler.Router.HandleFunc("/posts/{id}", handler.UpdatePost).Methods("PUT")
	handler.Router.HandleFunc("/posts/{id}", handler.DeletePost).Methods("DELETE")
	handler.Router.HandleFunc("/posts", handler.GetPosts).Methods("GET").Queries("limit", "{[0-9]*?}")
	handler.Router.HandleFunc("/posts", handler.GetPosts).Methods("GET")
}
