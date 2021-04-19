package controllers

import (
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
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}
