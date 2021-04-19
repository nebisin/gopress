package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/nebisin/gopress/models"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils/auth"
	"github.com/nebisin/gopress/utils/responses"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

var validate = validator.New()

func (handler *Handler) handleAuthRegister(w http.ResponseWriter, r *http.Request)  {
	var userPayload models.UserPayload
	if err := json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := validate.Var(userPayload.Email, "required,email"); err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("you have to provide a valid email"))
		return
	}

	if len(userPayload.Password) < 8 {
		responses.ERROR(w, http.StatusBadRequest, errors.New("password must be at least 8 characters"))
		return
	}

	user := models.PayloadToUser(userPayload)

	db := repository.NewUserRepository(handler.DB)

	if err := db.Save(&user); err != nil {
		if strings.Contains(err.Error(), "email") {
			responses.ERROR(w, http.StatusBadRequest, errors.New("email is already taken"))
			return
		}
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		fmt.Println(err)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		fmt.Println(err)
		return
	}

	responses.JSON(w, http.StatusCreated, token)
}

func (handler Handler) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	var userPayload models.UserPayload
	if err := json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := repository.NewUserRepository(handler.DB)

	user, err := db.FindByEmail(userPayload.Email)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("email or password is wrong"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("email or password is wrong"))
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, token)
}
