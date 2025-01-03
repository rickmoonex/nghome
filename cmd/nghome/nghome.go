package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rickmoonex/nghome/internal/system/database"
	"github.com/rickmoonex/nghome/internal/system/eventbus"
	"github.com/rickmoonex/nghome/pkg/framework/helper"
	"github.com/rickmoonex/nghome/pkg/framework/instance"
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

	testSwitch := instance.SwitchInstance{}

	ctx := &helper.NGContext{}

	err = testSwitch.Init(ctx, map[string]interface{}{
		"name": "test_switch_2",
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	testSwitch.TurnOn()

	for {
	}
}
