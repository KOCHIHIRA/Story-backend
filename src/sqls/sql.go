package sqls

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	User     string
	Password string
	DBname   string
}

func NewDB(DBname string) *sql.DB {
	DB_KEYWORD := "root:K4143568k@tcp(localhost:3306)/" + DBname
	db, err := sql.Open("mysql", DB_KEYWORD)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
