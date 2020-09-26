package server

import (
	"github.com/sergeychur/avito_auto/internal/models"
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
	err = ValidateLink(link)
	if err != nil {
		WriteToResponse(w, http.StatusBadRequest, "Link is invalid")
		return
	}
	status, link := server.repo.InsertLink(link)
	DealRequestFromRepo(w, &link, status)
}

func (server *Server) GetLink(w http.ResponseWriter, r *http.Request) {
	newUrl := "https://google.com"
	http.Redirect(w, r, newUrl, http.StatusSeeOther)
}
