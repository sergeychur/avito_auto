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
	/*slugPattern := "^(\\d|\\w|-|_)*(\\w|-|_)(\\d|\\w|-|_)*$"
	idPattern := "^[0-9]+$"
	nickPattern := "^[A-Za-z0-9_\\.-]+$"*/

	subRouter := chi.NewRouter()
	//subRouter.Get("/link/get")

	r.Mount("/api/", subRouter)
	server.router = r

	newConfig, err := config.NewConfig(pathToConfig)
	if err != nil {
		log.Println("Cannot create config instance because of: ", err)
		return nil, err
	}
	server.config = newConfig
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
