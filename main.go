package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Message builds json message
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

//Respond builds json responses
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Declare all the templates
var tmpl = template.Must(template.ParseGlob("views/pages/*"))

// Function to get the 7 Days Average and return string(float64(avg))
func (rts Rates) sevenDaysAverage() string {
	var avg float64
	var sum float64
	i := 0.0
	for _, rt := range rts {
		sum += rt.Rate
		i++
	}
	if i != 7.0 {
		return "insufficient data"
	}
	avg = sum / i
	return fmt.Sprintf("%g", avg)
}

func showExchangeRates(w http.ResponseWriter, r *http.Request) { //Method= GET ==> output = all exchange rates from db
	db := dbConn()
	date := r.URL.Query().Get("date")
	exchanges, err := db.Query("SELECT * FROM exchange ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	exc := Exchange{}
	for exchanges.Next() {
		var id int
		var from, to string

		err = exchanges.Scan(&id, &from, &to)
		if err != nil {
			panic(err.Error())
		}
		exid := strconv.Itoa(id)
		rates, err := db.Query("SELECT * from rates WHERE exchange_id=" + exid + " ORDER BY id ASC")
		rt := Rate{}
		var rts Rates
		for rates.Next() {
			var rate int
			var date string
			err = rates.Scan(&rate, &date)
			if err != nil {
				panic(err.Error())
			}
			rt.Date = date
			rt.Rate = rate
			rts = append(rts, rt)
		}

		exc.ID = id
		exc.From = from
		exc.To = to
	}

	defer db.Close()
}

func getExchangeRate(w http.ResponseWriter, r *http.Request) {

}

func getExchangeRateOnDate(w http.ResponseWriter, r *http.Request) { //Method= GET; input = date format->("Y-m-d") ==> output = rate on the date and <=7d avg rates

}

func createExchangeRate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	from := r.FormValue("From")
	to := r.FormValue("To")
	insert, err := db.Prepare("INSERT INTO exchange(from, to) VALUES(?,?)")
	if err != nil {
		Respond(w, Message(false, "Error while INSERTING values INTO exchange"))
		return
	}
	insert.Exec(from, to)
	Respond(w, Message(true, "Successfully Added New Exchange Currencies"))
	defer db.Close()
}

func showDailyExchangeRateForm(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nID := r.URL.Query().Get("id")
	sel, err := db.Query("SELECT * FROM exchange WHERE id=?", nID)
	if err != nil {
		panic(err.Error())
	}
	exc := Exchange{}
	for sel.Next() {
		var id int
		var from, to string
		err = sel.Scan(&id, &from, &to)
		if err != nil {
			panic(err.Error())
		}
		exc.ID = id
		exc.From = from
		exc.To = to
	}
	tmpl.ExecuteTemplate(w, "Form_DailyExrate", exc)
	defer db.Close()
}

func inputDailyExchangeRate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		from := r.FormValue("From")
		to := r.FormValue("To")
		rt := r.FormValue("Rate")
		date := r.FormValue("Date")

		insert, err := db.Prepare("")
	}
}

func deleteExchangeRate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	exc := r.URL.Query().Get("id")
	del, err := db.Prepare("DELETE FROM exchange WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	del.Exec(exc)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/exchange-rates", showExchangeRates).Methods("GET")
	log.Println("Server started on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
