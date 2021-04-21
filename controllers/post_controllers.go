package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/nebisin/gopress/models"
	"github.com/nebisin/gopress/repository"
	"github.com/nebisin/gopress/utils/auth"
	"github.com/nebisin/gopress/utils/responses"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// handlePostCreate method create a post that send with body.
// Only authenticated users create a post.
func (handler *Handler) handlePostCreate(w http.ResponseWriter, r *http.Request)  {
	var postDTO models.PostDTO
	if err := json.NewDecoder(r.Body).Decode(&postDTO); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.DTOToPost(postDTO)

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

// handlePostGet method get the post by given id.
// If post is published everybody can read it.
// But post is not published only the author can access it.
func (handler *Handler) handlePostGet(w http.ResponseWriter, r *http.Request)  {
	// We get the id in url and parse it as uint type
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	// We try to find the post with given id
	post, err := db.FindById(uint(i))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id " + id + " could not found"))
		} else {
			// If method is failed for another reason than "record not found"
			// We don't want to share that reason with user
			// Instead we send a generic error to the user
			// and print the actual error to the console
			responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
			log.Println(err)
		}
		return
	}

	// If post is not published only the author can access it.
	if post.IsPublished == false {
		uid, err := auth.ExtractTokenID(r)
		if err != nil {
			// If the requester not authenticated we pretend like post is not exist
			// for protection against data leak.
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id " + id + " could not found"))
			return
		}

		if uid != post.Author.ID {
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id " + id + " could not found"))
			return
		}
	}

	responses.JSON(w, http.StatusOK, post)
}

// handlePostUpdate method update the post by given id and body.
// It requires authentication and user must be the owner of the post.
func (handler Handler) handlePostUpdate(w http.ResponseWriter, r *http.Request)  {
	// We try to get the user id from auth token:
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	vars := mux.Vars(r)
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
			log.Println(err)
		}
		return
	}

	if post.Author.ID != uid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("you can not update the post who belongs to someone else"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var postUpdate models.PostDTO

	if err = json.Unmarshal(body, &postUpdate); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	newPost := models.DTOToPost(postUpdate)

	if err = db.UpdateById(&post, newPost); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

// handlePostDelete method delete a post with it's id.
// It requires authentication and user must be the owner of the post.
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
		responses.ERROR(w, http.StatusNotFound, errors.New("the post with id " + id + " could not found"))
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if post.Author.ID != uid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("you can not delete the post who belongs to someone else"))
		return
	}

	if 	err := db.DeleteById(uint(i)); err != nil{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		log.Println(err)
		return
	}

	responses.JSON(w, http.StatusNoContent, "")
}

// handlePostGetMany method find many posts within the given limit.
// It only return published posts.
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