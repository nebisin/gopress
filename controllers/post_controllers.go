package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/nebisin/gopress/models"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils"
	"github.com/nebisin/gopress/utils/auth"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (handler *Handler) CreatePost(w http.ResponseWriter, r *http.Request)  {
	var postDTO models.PostDTO
	if err := json.NewDecoder(r.Body).Decode(&postDTO); err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	post := models.DTOToPost(postDTO)

	// TODO: Some validations

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	post.AuthorID = &uid

	db := repository.NewPostRepository(handler.DB)

	if err := db.Save(&post); err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusCreated, post)
}

func (handler *Handler) GetPostById(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	post, err := db.FindById(uint(i))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusOK, post)
}

func (handler Handler) UpdatePost(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	// TODO: Some validations

	post, err := db.FindById(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// TODO: Check if the user is authorized

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	var postUpdate models.PostDTO

	if err = json.Unmarshal(body, &postUpdate); err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = db.UpdateById(&post, models.DTOToPost(postUpdate)); err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusCreated, post)
}

func (handler *Handler) DeletePost(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	_, err = db.FindById(uint(i))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// TODO: Check if the user is authorized

	if 	err := db.DeleteById(uint(i)); err != nil{
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusNoContent, "")
}

func (handler Handler) GetPosts(w http.ResponseWriter, r *http.Request)  {
	keys := r.URL.Query()
	limitStr := keys.Get("limit")

	var limit int
	var err error
	if len(limitStr) != 0 {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			utils.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
	}

	db := repository.NewPostRepository(handler.DB)

	posts, err := db.FindMany(limit)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusOK, posts)
}