package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sergeychur/avito_auto/internal/config"
	"github.com/sergeychur/avito_auto/internal/middlewares"
	"github.com/sergeychur/avito_auto/internal/repository"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	Router *chi.Mux
	Repo   repository.Repository
	Config *config.Config
	Validator IValidator
}

func NewServer(pathToConfig string) (*Server, error) {
	server := new(Server)
	r := chi.NewRouter()
	newConfig, err := config.NewConfig(pathToConfig)
	if err != nil {
		log.Println("Cannot create Config instance because of: ", err)
		return nil, err
	}
	server.Config = newConfig
	r.Use(middleware.Logger,
		middleware.Recoverer,
		middlewares.CreateCorsMiddleware(server.Config.AllowedHosts))
	subRouter := chi.NewRouter()
	subRouter.Get("/link/{shortcut}", server.GetLink)
	subRouter.Post("/link", server.CreateLink)

	r.Mount("/api/", subRouter)
	r.Get("/doc/{file:.+\\..+$}", http.StripPrefix("/doc/",
		http.FileServer(http.Dir(server.Config.DocPath))).ServeHTTP)
	server.Router = r

	dbPort, err := strconv.Atoi(server.Config.DBPort)
	if err != nil {
		return nil, err
	}
	repo := repository.NewRepository(server.Config.DBUser, server.Config.DBPass,
		server.Config.DBName, server.Config.DBHost, uint16(dbPort))
	server.Repo = repo
	server.Validator = NewValidator(server.Config.ValidRequestTimeout)
	return server, nil
}

func (server *Server) Run() error {
	err := server.Repo.Start(server.Config.DBMaxConn, server.Config.DBAcquireTimeout)
	if err != nil {
		log.Printf("Failed to connect to DB: %s", err.Error())
		return err
	}
	defer server.Repo.Close()
	port := server.Config.Port
	log.SetOutput(os.Stdout)
	log.Printf("Running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Router))
	return nil
}
