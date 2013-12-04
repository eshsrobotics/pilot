package main

import (
	"database/sql"
	"log"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "pilot.db")
	if err != nil {
		log.Fatalln("sql.Open failed", err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(Submission{}, "submissions").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Create tables failed", err)
	}

	return dbmap
}
