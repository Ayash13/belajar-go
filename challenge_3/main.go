package main

import (
	"belajar-go/challenge_3/database"
	"belajar-go/challenge_3/handler"
	"belajar-go/challenge_3/repository"
	"belajar-go/challenge_3/server"
	"belajar-go/challenge_3/service"
	"log"
	"net/http"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	accountService := service.NewAccountService(accountRepo, transactionRepo)

	mux := http.NewServeMux()
	accountHandler := handler.NewAccountHandler(mux, accountService)
	accountHandler.MapRoutes()

	log.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", server.ApplicationMiddlewareResponse(server.HandleRoutesNotFound(mux)))
	if err != nil {
		log.Fatal(err)
	}
}
