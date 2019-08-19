package main

import "database/sql"

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbProtocol := "ip4:1"
	dbAddress := "192.168.64.2"
	dbName := "be_test"
	dbUser := "ipcc"
	dbPass := "Makanan1%"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@"+dbProtocol+"("+dbAddress+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
