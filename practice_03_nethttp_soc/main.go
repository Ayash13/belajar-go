package main

import (
	"log"
	"net/http"
	"nethttp/database"
	"nethttp/handler"
	"nethttp/repository"
	"nethttp/server"
	"nethttp/service"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	personRepo := repository.NewPersonRepository(db)
	personService := service.NewPersonService(personRepo)

	mux := http.NewServeMux()
	personHandler := handler.NewPersonHandler(mux, personService)
	personHandler.MapRoutes()

	log.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", server.ApplicationMiddlewareResponse(server.HandleRoutesNotFound(mux)))
	if err != nil {
		log.Fatal(err)
	}
}
