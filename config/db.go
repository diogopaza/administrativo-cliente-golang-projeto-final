package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" //postgresql
)

const (
	host           = "localhost"
	port           = 5432
	user           = "postgres"
	password_admin = "123321"
	dbname         = "lupatini-cursos"
)

func DB() *sql.DB {

	banco := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable\n", host, port, user, password_admin, dbname)

	db, err := sql.Open("postgres", banco)
	if err != nil {
		fmt.Println("Erro ao conectar no banco de dados")
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Sucessfully connected!!!")
	return db

}
