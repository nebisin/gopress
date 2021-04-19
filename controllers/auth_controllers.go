package controllers

import (
	"encoding/json"
	"github.com/nebisin/gopress/models"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils/auth"
	"github.com/nebisin/gopress/utils/responses"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (handler *Handler) handleAuthRegister(w http.ResponseWriter, r *http.Request)  {
	var userPayload models.UserPayload
	if err := json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user := models.PayloadToUser(userPayload)

	db := repository.NewUserRepository(handler.DB)

	if err := db.Save(&user); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, token)
}

func (handler Handler) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	var userPayload models.UserPayload
	if err := json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	db := repository.NewUserRepository(handler.DB)

	user, err := db.FindByEmail(userPayload.Email)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, token)
}
