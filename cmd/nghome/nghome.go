package main

import (
	"log"

	"github.com/rickmoonex/nghome/internal/system/database"
	"github.com/thingsdb/go-thingsdb"
)

func main() {
	conn := thingsdb.NewConn("localhost", 9200, nil)

	if err := conn.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := conn.AuthPassword("admin", "pass"); err != nil {
		log.Fatal(err)
	}

	db := &database.Client{Conn: conn}
	migC := db.GetMigrationClient()

	if err := migC.AutoMigrate("./migrations"); err != nil {
		log.Fatal(err)
	}
}
