package eventbus

import (
	"errors"
	"time"
)

// Event represents an entry in the recorded events set.
type Event struct {
	UniqueId  int
	Type      string
	Data      string
	TimeFired time.Time
}

// FromInterface decoded the interface{} returned from ThingsDB into the struct
func (e *Event) FromInterface(res interface{}) error {
	resMap, ok := res.(map[string]interface{})
	if !ok {
		return errors.New("unable to cast res to map[string]interface")
	}

	uniqueId, ok := resMap["unique_id"].(int8)
	if !ok {
		return errors.New("unable to cast unique_id into int8")
	}

	eventType, ok := resMap["type"].(string)
	if !ok {
		return errors.New("unable to cast event type into string")
	}

	data, ok := resMap["data"].(string)
	if !ok {
		return errors.New("unable to cast data into string")
	}

	tfStr, ok := resMap["time_fired"].(string)
	if !ok {
		return errors.New("unable to cast time_fired into string")
	}
	timeFired, err := time.Parse(time.RFC3339, tfStr)
	if err != nil {
		return err
	}

	e.UniqueId = int(uniqueId)
	e.Type = eventType
	e.Data = data
	e.TimeFired = timeFired

	return nil
}
