package mysql

import (
	"database/sql"
 _ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/go_main")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}