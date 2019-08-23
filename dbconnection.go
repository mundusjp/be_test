package main

import "database/sql"

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbProtocol := "tcp"
	dbAddress := "localhost:3306"
	dbName := "be_test"
	dbUser := "root"
	db, err := sql.Open(dbDriver, dbUser+"@"+dbProtocol+"("+dbAddress+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
