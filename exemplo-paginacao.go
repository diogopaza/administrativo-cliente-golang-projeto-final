package main

import (
	"database/sql"
	"lupatini/config"
)

var connectDB *sql.DB

func main() {

	connectDB = config.DB()

}
