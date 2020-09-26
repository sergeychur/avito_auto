package main

import (
	"github.com/sergeychur/avito_auto/internal/server"
	"log"
	"os"
)

func main() {
	pathToConfig := ""
	if len(os.Args) != 2 {
		panic("Usage: ./main <path_to_config>")
	} else {
		pathToConfig = os.Args[1]
	}
	serv, err := server.NewServer(pathToConfig)
	if err != nil {
		log.Println(err)
		return
	}
	err = serv.Run()
	if err != nil {
		panic(err.Error())
	}
}