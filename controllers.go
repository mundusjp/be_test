package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/now"
)

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

// MinMax is a function to return min and max value of float array
func MinMax(array []float64) (float64, float64) {
	max := array[0]
	min := array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func showExchangeRates(w http.ResponseWriter, r *http.Request) {
	var excs []Exchange
	var response Response
	exc := Exchange{}
	db := dbConn()
	date := r.URL.Query().Get("date")
	startday, err := now.Parse(date)

	if err != nil {
		panic(err.Error())
	}
	lastweek := startday.AddDate(0, 0, -6)

	exchanges, err := db.Query("SELECT * FROM exchange ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}

	for exchanges.Next() {
		var id, ids int
		var from, to, avg string
		var rts Rates
		rt := Rate{}
		err = exchanges.Scan(&id, &from, &to)
		if err != nil {
			panic(err.Error())
		}
		exid := id
		rates, err := db.Query("SELECT * from rates WHERE exchange_id=" + strconv.Itoa(exid) + " AND date BETWEEN '" + lastweek.Format("2006-01-02") + "' AND '" + startday.Format("2006-01-02") + "' ORDER BY id ASC ")
		if err != nil {
			panic(err.Error())
		}
		for rates.Next() {
			var rate float64
			var date string
			err = rates.Scan(&ids, &exid, &rate, &date)
			if err != nil {
				log.Print("Error in Scanning Rates")
				panic(err.Error())

			}
			rt.ID = ids
			rt.ExchangeID = exid
			rt.Date = date
			rt.Rate = rate
			rts = append(rts, rt)
		}

		avg = rts.sevenDaysAverage()
		exc.Average = avg
		exc.ID = id
		exc.From = from
		exc.To = to
		if strings.TrimRight(avg, "\n") == "insufficient data" {
			exc.Rate = avg
		} else {
			rates, err := db.Query("SELECT rate from rates WHERE exchange_id=" + strconv.Itoa(exid) + " AND date = '" + startday.Format("2006-01-02") + "'")
			if err != nil {
				panic(err.Error())
			}
			for rates.Next() {
				var rate float64
				err = rates.Scan(&rate)
				if err != nil {
					log.Print("Error in Scanning Rate")
					panic(err.Error())
				}
				exc.Rate = fmt.Sprintf("%g", rate)
			}
		}

		excs = append(excs, exc)
	}
	response.Status = 200
	response.Message = "Success"
	response.Data = excs

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	defer db.Close()
}

func getExchangeRate(w http.ResponseWriter, r *http.Request) {
	var excs []Exchange
	var response Response
	exc := Exchange{}
	db := dbConn()
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	exchanges, err := db.Query("SELECT * FROM exchange WHERE `from`='" + from + "' AND `to`='" + to + "' ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}

	for exchanges.Next() {
		var id, ids int
		var from, to, avg string
		var rts Rates
		var arr []float64
		var min, max float64
		rt := Rate{}
		err = exchanges.Scan(&id, &from, &to)
		if err != nil {
			panic(err.Error())
		}
		exid := id
		rates, err := db.Query("SELECT * from rates WHERE exchange_id=" + strconv.Itoa(exid) + " ORDER BY id DESC LIMIT 7")
		if err != nil {
			panic(err.Error())
		}
		for rates.Next() {
			var rate float64
			var date string
			err = rates.Scan(&ids, &exid, &rate, &date)
			if err != nil {
				log.Print("Error in Scanning Rates")
				panic(err.Error())

			}
			rt.ID = ids
			rt.ExchangeID = exid
			rt.Date = date
			rt.Rate = rate
			arr = append(arr, rt.Rate)
			rts = append(rts, rt)
		}
		min, max = MinMax(arr)
		variance := max - min
		avg = rts.sevenDaysAverage()
		exc.Average = avg
		exc.ID = id
		exc.From = from
		exc.To = to
		exc.Rates = rts
		exc.Variance = fmt.Sprintf("%g", variance)
		excs = append(excs, exc)
	}
	response.Status = 200
	response.Message = "Success"
	response.Data = excs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	defer db.Close()
}

func createExchangeRate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var response Response
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	insert, err := db.Prepare("INSERT INTO exchange (`from` , `to`) VALUES (?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = insert.Exec(from, to)
	if err != nil {
		panic(err.Error())
	}
	response.Status = 200
	response.Message = "Success"
	log.Print("Inserted data to database")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	defer db.Close()
}

func inputDailyExchangeRate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var id int
	var response Response
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	rt := r.URL.Query().Get("rate")
	date := r.URL.Query().Get("date")
	sel, err := db.Query("SELECT id FROM exchange WHERE `from`='" + from + "' AND `to`='" + to + "'")
	if err != nil {
		panic(err.Error())
	}
	for sel.Next() {
		err = sel.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		insert, err := db.Prepare("INSERT INTO rates (exchange_id ,rate, date) VALUES (?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(id, rt, date)
	}
	response.Status = 200
	response.Message = "Success"
	log.Print("Inserted data to database")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteExchangeRate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var response Response
	id := r.URL.Query().Get("id")
	del, err := db.Prepare("DELETE FROM exchange WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	del.Exec(id)
	log.Print("Deleted the Data from Database")
	defer db.Close()
	response.Status = 200
	response.Message = "Success"
	log.Print("Inserted data to database")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
