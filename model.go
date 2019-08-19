package main

// Rate struct initiation #convert to json
type Rate struct {
	ID         int     `json:"id"`
	ExchangeID int     `json:"exchange_id"`
	Rate       float64 `json:"rate"`
	Date       string  `json:"date"`
}

// Rates slices
type Rates []Rate

// Average struct initiation #convert to json
type Average struct {
	Avg string `json:"7-day-avg"`
}

// Exchange struct initiation #convert to json
type Exchange struct {
	ID      int     `json:"id"`
	From    string  `json:"from"`
	To      string  `json:"to"`
	Rates   Rates   `json:"rates"`
	Average Average `json:"7-day-avg"`
}
