package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rickmoonex/nghome/internal/system/database"
	"github.com/rickmoonex/nghome/internal/system/eventbus"
)

func onMessage(args []interface{}) {
	fmt.Println("changed")
}

func onUpdated(args []interface{}) {
	fmt.Println("updated")
}

func main() {
	db, err := database.InitializeClient("localhost", 9200, "K5N2NM2hcfD0FtPIQg5ATb")
	if err != nil {
		log.Fatal(err)
	}

	migC := db.GetMigrationClient()

	if err := migC.AutoMigrate("./migrations"); err != nil {
		log.Fatal(err)
	}

	eb, err := eventbus.InitEventBus()
	if err != nil {
		log.Fatal(err)
	}

	eb.Listen("state_changed", onMessage)
	eb.Listen("state_updated", onUpdated)

	time.Sleep(time.Second * 600)
}
