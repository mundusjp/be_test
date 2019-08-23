package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/now"
	"github.com/thedevsaddam/govalidator"
)

// Function to get the 7 Days Average and return string(float64(avg))
func (rts Rates) sevenDaysAverage() string {
	var avg float64
	var sum float64
	i := 0.0
	// loop through all rts rows. it will sum the rate and count the rows
	for _, rt := range rts {
		sum += rt.Rate
		i++
	}

	// if count of the rows is not 7, return string "insufficient data"
	if i != 7.0 {
		return "insufficient data"
	}
	// if not, get the average of the rates
	avg = sum / i

	// return the converted float(avg)
	return fmt.Sprintf("%g", avg)
}

// MinMax is a function to return min and max values from an array of floats
func MinMax(array []float64) (float64, float64) {
	max := array[0] // set the first value of array as max
	min := array[0] // set the first value of array as min

	// loop through all of the array rows
	for _, value := range array {

		// getting the maximum value of an array
		// if the next value have a bigger value, set the next value as max
		if max < value {
			max = value
		}

		// getting the minimum value of an array
		// if the next value have a smaller value, set the next value as min
		if min > value {
			min = value
		}
	}

	// return min and max values as float
	return min, max
}

/*
******************************************************************
		Show all of the Exchange Rate with their Rates
		for use case: User has a list of exchange rates to be tracked
******************************************************************
Firstly, the function will validate all of the request parameters.
If validation error occured, the system will show error log on console.

Next, the function will pull all of  exchange rate values from database and their
rates. It will get the rates from last 7 days including the date. Average() will then count
the rows of the rates. If count() < 7 it will return "insufficient data", else it will
return float(avg).

Input: []date date
Output: JSON{status:"", message:"", Data:"[]Exchange{}]"}
*/
func showExchangeRates(w http.ResponseWriter, r *http.Request) {
	// request fields : date
	rules := govalidator.MapData{ // fields requirement & validation rules
		"date": []string{"required", "date"},
	}

	opts := govalidator.Options{
		Request:         r,     // request object
		Rules:           rules, // rules map
		RequiredDefault: true,  //all field must be passed to the rules
	}
	// define new validator
	v := govalidator.New(opts)
	// initiate validation()
	e := v.Validate()
	// check if any rule violation or validation error
	vErr := map[string]interface{}{"validationError": e} // validation error variable

	// if there is a validation error, it will return the error in json
	if vErr != nil {
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(vErr)
		panic(vErr)
	}

	// define variables
	var excs []Exchange   // slice of type Exchange
	var response Response // struct of response message
	exc := Exchange{}     // defining exc as type Exchange

	// initiate connection to database. configuration in dbconnection.go
	db := dbConn()

	// getting the required date from r Request
	date := r.URL.Query().Get("date")

	// convert string("date") to time("date")
	startday, err := now.Parse(date) // startday is the 7th date of query date
	if err != nil {
		panic(err.Error()) // log the error if parsing failed
	}

	// get the first date of query date
	lastweek := startday.AddDate(0, 0, -6)

	// query to get all of the exchange rates from database
	exchanges, err := db.Query("SELECT * FROM exchange ORDER BY id ASC")
	if err != nil {
		panic(err.Error()) // log the error if query failed
	}

	// loop through all exchanges rows
	for exchanges.Next() {
		// initiate variables that will be used to store data from database
		var id, ids int
		var from, to, avg string
		var rts Rates
		rt := Rate{}                          // type Rate in model.go
		err = exchanges.Scan(&id, &from, &to) // assign values to id, from, and to
		if err != nil {
			panic(err.Error()) // log the error if value assignments failed
		}
		// defining exchange_id = id
		exid := id

		// query to get all of the rates from each exchange rate
		rates, err := db.Query("SELECT * from rates WHERE exchange_id=" + strconv.Itoa(exid) + " AND date BETWEEN '" + lastweek.Format("2006-01-02") + "' AND '" + startday.Format("2006-01-02") + "' ORDER BY id ASC ")
		if err != nil {
			panic(err.Error()) // log the error if query failed
		}

		// loop through all rates rows
		for rates.Next() {
			// defining variables to store data pulled from database
			var rate float64
			var date string

			// assign pulled data values to id , exchange_id, rate, and date
			err = rates.Scan(&ids, &exid, &rate, &date)
			if err != nil {
				log.Print("Error in Scanning Rates")
				panic(err.Error()) // log the error if assignment failed

			}

			// assign values to struct Rate
			rt.ID = ids
			rt.ExchangeID = exid
			rt.Date = date
			rt.Rate = rate
			// push the struct to slice of Rate struct
			rts = append(rts, rt)
		}

		// get the average from rates.
		// if rates count() < 7, return "insufficient data"
		// if rates count() == 7, return the converted float(avg)
		avg = rts.sevenDaysAverage()

		// assign values to Exchange struct
		exc.ID = id
		exc.From = from
		exc.To = to

		// get the rate of the requested date
		// if avg returns "insufficient data", assign avg value as the rate
		// if avg returns float(avg), get the rate of the date and assign avg value to Exchange struct
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
				exc.Average = avg
			}
		}
		// push the Exchange struct to slice
		excs = append(excs, exc)
	}

	// assigning values to Response struct
	response.Status = 200
	response.Message = "Success"
	response.Data = excs

	// create the json output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	// close connection to database
	defer db.Close()
}

/*
******************************************************************
		Show Exchange Rate with its Rates, Avg, and Variance
		for use case: User wants to see the exchange rate trend from the most recent 7 data points
******************************************************************
Firstly, the function will validate all of the request parameters.
If validation error occured, the system will show error log on console.

Next, the function will pull exchange rate values from database and its
rates. It will show the latest 7 rates order by id DESC. the function MinMax
will return the max and min values from the rate. It will be used to get the
variance of rates.

Input: []string from, []string to
Output: JSON{status:"", message:"", Data:"Exchange{}"}
*/
func getExchangeRate(w http.ResponseWriter, r *http.Request) {
	// request fields : from , to

	// defining the rules for validation
	rules := govalidator.MapData{ // fields requirement & validation rules
		"from": []string{"required", "between:3,4"},
		"to":   []string{"required", "between:3,4"},
	}

	// defining the options for validation
	opts := govalidator.Options{
		Request:         r,     // request object
		Rules:           rules, // rules map
		RequiredDefault: true,  //all field must be passed to the rules
	}

	// defining v as new validator
	v := govalidator.New(opts)
	// execute the validation
	e := v.Validate()

	// defining vErr as error message if validation error occured
	vErr := map[string]interface{}{"validationError": e} // validation error variable

	// if there is a validation error, it will return the error in json
	if vErr != nil {
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(vErr)
		panic(vErr)
	}

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

/*
******************************************************************
					Create new Exchange Rate
		for use case: User wants to add an exchange rate to the list
******************************************************************
It will insert new from and to data to exchange table in database.

Input: []string from, []string to
Output: JSON{status:"", message:""}, console log
*/
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
