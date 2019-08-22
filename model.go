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

// Floats slice
type Floats []float64

// Response untuk template json
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Exchange
}

// Exchange struct initiation #convert to json
type Exchange struct {
	ID       int    `json:"id"`
	From     string `json:"from"`
	To       string `json:"to"`
	Rates    Rates  `json:"rates"`
	Rate     string `json:"rate"`
	Average  string `json:"7-day-avg"`
	Variance string `json:"variance"`
}
