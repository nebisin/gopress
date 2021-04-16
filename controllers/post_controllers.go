package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nebisin/gopress/models"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (handler *Handler) CreatePost(w http.ResponseWriter, r *http.Request)  {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// TODO: Some validations
	// TODO: Authentication

	db := repository.NewRepository(handler.DB)

	if err := db.SavePost(&post); err != nil {
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

	db := repository.NewRepository(handler.DB)

	post, err := db.FindPostById(uint(i))
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

	db := repository.NewRepository(handler.DB)

	// TODO: Some validations
	// TODO: Authentication

	post, err := db.FindPostById(uint(pid))
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

	if err = db.UpdatePostById(&post, models.DTOToPost(postUpdate)); err != nil {
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

	db := repository.NewRepository(handler.DB)

	// TODO: Authentication

	_, err = db.FindPostById(uint(i))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// TODO: Check if the user is authorized

	if 	err := db.DeletePostById(uint(i)); err != nil{
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusNoContent, "")
}

func (handler Handler) GetPosts(w http.ResponseWriter, r *http.Request)  {
	limitStr := r.FormValue("limit")
	var limit int
	var err error
	if len(limitStr) != 0 {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			utils.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
	}

	db := repository.NewRepository(handler.DB)

	posts, err := db.FindPosts(limit)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusOK, posts)
}