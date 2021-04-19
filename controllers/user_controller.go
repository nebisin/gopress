package controllers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils/responses"
	"net/http"
	"strconv"
)

func (handler Handler) handleUserGet(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := repository.NewUserRepository(handler.DB)

	post, err := db.FindById(uint(i))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

func (handler Handler) HandleUserPostsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	posts, err := db.FindPostsByUserId(uint(i))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		fmt.Println(err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}
