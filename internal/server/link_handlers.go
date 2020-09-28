package server

import (
	"github.com/go-chi/chi"
	"github.com/sergeychur/avito_auto/internal/models"
	"github.com/sergeychur/avito_auto/internal/repository"
	"log"
	"net/http"
)

func (server *Server) CreateLink(w http.ResponseWriter, r *http.Request) {
	link := models.Link{}
	err := ReadFromBody(r, w, &link)
	if err != nil {
		log.Println("Cannot read from body because of: ", err)
		WriteToResponse(w, http.StatusBadRequest, nil)
		return
	}
	err = server.Validator.ValidateLink(link)
	if err != nil {
		log.Println("URL is invalid: ", err)
		WriteToResponse(w, http.StatusBadRequest, "Link is invalid")
		return
	}
	status, link := server.Repo.InsertLink(link)
	DealRequestFromRepo(w, &link, status)
}

func (server *Server) GetLink(w http.ResponseWriter, r *http.Request) {
	shortcut := chi.URLParam(r, "shortcut")
	status, newUrl := server.Repo.GetLink(shortcut)
	if status != repository.OK {
		DealRequestFromRepo(w, nil, status)
		return
	}
	http.Redirect(w, r, newUrl, http.StatusSeeOther)
}
