package controllers

import (
	"encoding/json"
	"errors"
	"github.com/nebisin/gopress/models"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils/auth"
	"github.com/nebisin/gopress/utils/responses"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (handler *Handler) handleAuthRegister(w http.ResponseWriter, r *http.Request)  {
	var userPayload models.UserPayload
	if err := json.NewDecoder(r.Body).Decode(&userPayload); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.PayloadToUser(userPayload)

	db := repository.NewUserRepository(handler.DB)

	if err := db.Save(&user); err != nil {
		if strings.Contains(err.Error(), "users.email") {
			responses.ERROR(w, http.StatusBadRequest, errors.New("email is already taken"))
			return
		}
		if strings.Contains(err.Error(), "users.username") {
			responses.ERROR(w, http.StatusBadRequest, errors.New("username is already taken"))
			return
		}
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		log.Println(err)
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

	user, err := db.FindByEmailOrUsername(userPayload.Email, userPayload.Username)
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

// handleMyPosts method gets users own posts
// including both published and unpublished ones.
func (handler Handler) handleMyPosts(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db := repository.NewPostRepository(handler.DB)
	posts, err := db.FindMyPosts(uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

// handleMe method return the authenticated user info.
func (handler Handler) handleMe(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db := repository.NewUserRepository(handler.DB)
	user, err := db.FindById(uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func (handler Handler) handleUpdateMe(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db := repository.NewUserRepository(handler.DB)

	user, err := db.FindById(uid)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var userUpdate models.UserDTO

	if err = json.Unmarshal(body, &userUpdate); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	newUser := models.DTOToUser(userUpdate)

	if err := db.UpdateById(&user, &newUser); err != nil {
		if strings.Contains(err.Error(), "users.email") {
			responses.ERROR(w, http.StatusBadRequest, errors.New("email is already taken"))
			return
		}
		if strings.Contains(err.Error(), "users.username") {
			responses.ERROR(w, http.StatusBadRequest, errors.New("username is already taken"))
			return
		}
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}