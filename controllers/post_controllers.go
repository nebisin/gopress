package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nebisin/gopress/models"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils/auth"
	"github.com/nebisin/gopress/utils/responses"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (handler *Handler) handlePostCreate(w http.ResponseWriter, r *http.Request)  {
	var postDTO models.PostDTO
	if err := json.NewDecoder(r.Body).Decode(&postDTO); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.DTOToPost(postDTO)

	if len(post.Title) < 3 {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("title must be at least 3 characters long"))
		return
	}
	if len(post.Body) < 3 {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("content must be at least 3 characters long"))
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	post.AuthorID = &uid

	db := repository.NewPostRepository(handler.DB)

	if err := db.Save(&post); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func (handler *Handler) handlePostGet(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	post, err := db.FindById(uint(i))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id " + id + " could not found"))
		} else {
			responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
			fmt.Println(err)
		}
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

func (handler Handler) handlePostUpdate(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	post, err := db.FindById(uint(pid))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id " + vars["id"] + " could not found"))
		} else {
			responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
			fmt.Println(err)
		}
		return
	}

	if post.Author.ID != uint(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("you can not update the post who belongs to someone else"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	var postUpdate models.PostDTO

	if err = json.Unmarshal(body, &postUpdate); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if len(postUpdate.Title) < 3 && postUpdate.Title != "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("title must be at least 3 characters long"))
		return
	}
	if len(postUpdate.Body) < 3 && postUpdate.Body != "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("content must be at least 3 characters long"))
		return
	}

	if err = db.UpdateById(&post, models.DTOToPost(postUpdate)); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func (handler *Handler) handlePostDelete(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	post, err := db.FindById(uint(i))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if post.Author.ID != uint(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if 	err := db.DeleteById(uint(i)); err != nil{
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, "")
}

func (handler Handler) handlePostGetMany(w http.ResponseWriter, r *http.Request)  {
	keys := r.URL.Query()
	limitStr := keys.Get("limit")

	var limit int
	var err error
	if len(limitStr) != 0 {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
	}

	db := repository.NewPostRepository(handler.DB)

	posts, err := db.FindMany(limit)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}