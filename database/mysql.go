package database

import (
	"database/sql"
 _ "github.com/go-sql-driver/mysql"
	"log"
)

func LoginTest() {
	db, err := sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/go_main")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from test ")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	var id int
	var name string
	for rows.Next() {
		err := rows.Scan(&id,&name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
}