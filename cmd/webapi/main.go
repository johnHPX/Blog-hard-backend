package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/configsAPI"
	"github.com/johnHPX/blog-hard-backend/internal/interf/routes"
)

func main() {
	log.Println("Initializing WebAPI")
	c := configsAPI.NewConfigs()
	projectConfigs, err := c.ProjectConfigs()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Settings Started...")
	log.Println(projectConfigs.Name)

	if projectConfigs.Port == "" {
		projectConfigs.Port = "40183"
	}

	log.Println("Initialized Routes")

	//init web service
	wsvc := routes.NewWebService()
	wsvc.Init()
	loggedRouter := handlers.LoggingHandler(os.Stdout, wsvc.GetRouters())
	//server setup
	srv := &http.Server{
		Handler:        loggedRouter,
		Addr:           fmt.Sprintf("0.0.0.0:%s", projectConfigs.Port),
		WriteTimeout:   800 * time.Second,
		ReadTimeout:    800 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening on Port %s", projectConfigs.Port)
	log.Fatal(srv.ListenAndServe())
}
