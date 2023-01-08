package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"./config"
	"./handler/bookHandler"
	"./handler/clientHandler"
	"./service"
	"./storage/bookStorage"
	"./storage/clientStorage"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), config.GetUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	clientStorage := clientStorage.NewClientStorage(dbpool)
	bookStorage := bookStorage.NewBookStorage(dbpool)
	service := service.NewServiceBook(bookStorage, clientStorage)
	bookHandler := bookHandler.NewBookHandler(service)
	clientHandler := clientHandler.NewClientHandler(service)

	route := mux.NewRouter()
	
	route.HandleFunc("/book/list", bookHandler.GetBooks).Methods("GET")
	route.HandleFunc("/book/{id}", bookHandler.GetBook).Methods("GET")
	route.HandleFunc("/book/create", bookHandler.CreateBook).Methods("POST")
	route.HandleFunc("/book/{id}", bookHandler.UpdateBook).Methods("PUT")
	route.HandleFunc("/book/{id}", bookHandler.DeleteBook).Methods("DELETE")
	route.HandleFunc("/book/take/{clientId}", bookHandler.TakeABook).Methods("POST")
	route.HandleFunc("/book/return/{id}", bookHandler.ReturnABook).Methods("GET")

	route.HandleFunc("/client/create", clientHandler.CreateClient).Methods("POST")
	route.HandleFunc("/client/list", clientHandler.GetClients).Methods("GET")
	route.HandleFunc("/client/{id}", clientHandler.DeleteClient).Methods("DELETE")
	route.HandleFunc("/client/{id}", clientHandler.GetClient).Methods("GET")
	route.HandleFunc("/client/{id}", clientHandler.UpdateClient).Methods("PUT")
	route.HandleFunc("/client/booklist/{id}", bookHandler.GetBooksByClientId).Methods("GET")

	http.Handle("/", route)

	srv := &http.Server{
		Handler: route,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
