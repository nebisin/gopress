package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nebisin/gopress/models"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils"
	"github.com/nebisin/gopress/utils/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request)  {
	var userPayload models.UserPayload
	if err := json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user := models.PayloadToUser(userPayload)

	db := repository.NewUserRepository(handler.DB)

	if err := db.Save(&user); err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusCreated, token)
}

func (handler Handler) Login(w http.ResponseWriter, r *http.Request) {
	var userPayload models.UserPayload
	if err := json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	db := repository.NewUserRepository(handler.DB)

	user, err := db.FindByEmail(userPayload.Email)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusCreated, token)
}

func (handler Handler) GetUserById(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	db := repository.NewUserRepository(handler.DB)

	post, err := db.FindById(uint(i))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusOK, post)
}
