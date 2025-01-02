package eventbus

import (
	"errors"
	"fmt"

	"github.com/rickmoonex/nghome/internal/system/database"
	"github.com/thingsdb/go-thingsdb"
)

// globalEventBus is a globally useable instance of the event bus
var globalEventBus *EventBus

// EventBus holds the utilities and data for using the event bus
type EventBus struct {
	dbClient *database.Client
	handlers map[string]func(args []interface{})
}

// GetEventBus returns the global event bus
func GetEventBus() (*EventBus, error) {
	if globalEventBus == nil {
		return nil, errors.New("global event bus not initialized")
	}
	return globalEventBus, nil
}

// InitEventBus initiates the global event bus
func InitEventBus() (*EventBus, error) {
	client, err := database.GetClient()
	if err != nil {
		return nil, err
	}

	eventBus := EventBus{dbClient: client, handlers: map[string]func(args []interface{}){}}

	globalEventBus = &eventBus

	roomId, err := globalEventBus.getRoomId()
	if err != nil {
		return nil, err
	}

	room := thingsdb.NewRoomFromId("//user_space", uint64(roomId))
	room.OnEmit = globalEventBus.onEmit
	if err := room.Join(globalEventBus.dbClient.Conn, thingsdb.DefaultWait); err != nil {
		return nil, err
	}

	return globalEventBus, nil
}

// onEmit is called when a new event is received
func (b *EventBus) onEmit(_ *thingsdb.Room, event string, args []interface{}) {
	callback, ok := b.handlers[event]
	if !ok {
		fmt.Printf("event `%s` received, no handler present\n", event)
		return
	}

	go callback(args)
}

// Listen will listen for a event type and fire the callback function
func (b *EventBus) Listen(eventType string, callback func(args []interface{})) {
	b.handlers[eventType] = callback
}

// FireEvent fires a new event over the event bus
func (b *EventBus) FireEvent(eventType, data string) (*Event, error) {
	client, err := database.GetClient()
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"type": eventType,
		"data": data,
	}

	var newEvent Event

	res, err := client.Query("//user_space", "wse(); .event_bus.fire_event(type, data);", vars)
	if err != nil {
		return nil, err
	}

	if err := newEvent.FromInterface(res); err != nil {
		return nil, err
	}

	return &newEvent, nil
}

// GetRoomId returns the room id of the event bus
func (b *EventBus) getRoomId() (int, error) {
	client, err := database.GetClient()
	if err != nil {
		return 0, err
	}

	res, err := client.Query("//user_space", ".event_bus.get_room_id();", nil)
	if err != nil {
		return 0, nil
	}

	roomId, ok := res.(int8)
	if !ok {
		return 0, errors.New("unable to cast room id to int8")
	}

	return int(roomId), nil
}
