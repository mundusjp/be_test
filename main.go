package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// initiate a new router
	r := mux.NewRouter()

	//defining all the routes and serve at port:8080
	r.HandleFunc("/api/exchange-rate", getExchangeRate).Methods("GET")
	r.HandleFunc("/api/exchange-currency/delete", deleteExchangeCurrency).Methods("DELETE")
	r.HandleFunc("/api/exchange-currency/insert", createExchangeCurrency).Methods("POST")
	r.HandleFunc("/api/exchange-rates", showExchangeRates).Methods("GET")
	r.HandleFunc("/api/daily-exchange-rates/insert", inputDailyExchangeRate).Methods("POST")
	r.Handle("/", r)
	log.Println("Server started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
