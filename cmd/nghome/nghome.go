package main

import (
	"log"

	"github.com/rickmoonex/nghome/internal/system/database"
)

func main() {
	db, err := database.InitializeClient("localhost", 9200, "K5N2NM2hcfD0FtPIQg5ATb")
	if err != nil {
		log.Fatal(err)
	}

	migC := db.GetMigrationClient()

	if err := migC.AutoMigrate("./migrations"); err != nil {
		log.Fatal(err)
	}
}
