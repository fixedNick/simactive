package sql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func MustInit() *sql.DB {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/simactive")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return db
}
