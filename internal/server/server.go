package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sergeychur/avito_auto/internal/config"
	"github.com/sergeychur/avito_auto/internal/repository"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	router *chi.Mux
	repo 	repository.Repository
	config *config.Config
}

func NewServer(pathToConfig string) (*Server, error) {
	server := new(Server)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	newConfig, err := config.NewConfig(pathToConfig)
	if err != nil {
		log.Println("Cannot create config instance because of: ", err)
		return nil, err
	}
	server.config = newConfig

	subRouter := chi.NewRouter()
	subRouter.Get("/link/get/{shortcut}", server.GetLink)
	subRouter.Post("/link/create", server.CreateLink)

	r.Mount("/api/", subRouter)
	r.Get("/doc/{file:.+\\..+$}", http.StripPrefix("/doc/",
		http.FileServer(http.Dir(server.config.DocPath))).ServeHTTP)
	server.router = r

	dbPort, err := strconv.Atoi(server.config.DBPort)
	if err != nil {
		return nil, err
	}
	repo := repository.NewRepository(server.config.DBUser, server.config.DBPass,
		server.config.DBName, server.config.DBHost, uint16(dbPort))
	server.repo = repo
	return server, nil
}

func (server *Server) Run() error {
	err := server.repo.Start(server.config.DBMaxConn, server.config.DBAcquireTimeout)
	if err != nil {
		log.Printf("Failed to connect to DB: %s", err.Error())
		return err
	}
	defer server.repo.Close()
	port := server.config.Port
	log.SetOutput(os.Stdout)
	log.Printf("Running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, server.router))
	return nil
}
