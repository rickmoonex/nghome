package statemachine

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// State represents an entry in the state machine
type State struct {
	UniqueId         int
	InstanceUniqueId int
	State            string
	Attributes       *map[string]interface{}
	LastChanged      time.Time
	LastUpdated      time.Time
	OldStateId       int
}

// FromInterface decodes the interface{} response from ThingsDB into the struct
func (s *State) FromInterface(res interface{}) error {
	resMap, ok := res.(map[string]interface{})
	if !ok {
		return errors.New("unable to cast res into map[string]interface{}")
	}

	uniqueId, ok := resMap["unique_id"].(int8)
	if !ok {
		return errors.New("unable to cast unique_id into int8")
	}

	instanceMap, ok := resMap["instance"].(map[string]interface{})
	if !ok {
		return errors.New("unable to cast instance into map[string]interface{}")
	}
	instanceUniqueId, ok := instanceMap["unique_id"].(int8)
	if !ok {
		return errors.New("unable to cast instance unique_id into int8")
	}

	state, ok := resMap["state"].(string)
	if !ok {
		return errors.New("unable to cast state into string")
	}

	attributesStr, ok := resMap["attributes"].(string)
	if !ok {
		return errors.New("unable to cast attributes into string")
	}
	fmt.Printf("attr str: %s\n", attributesStr)
	attributes := &map[string]interface{}{}
	err := json.Unmarshal([]byte(attributesStr), attributes)
	if err != nil {
		return err
	}

	lcStr, ok := resMap["last_changed"].(string)
	if ok {
		lc, err := time.Parse(time.RFC3339, lcStr)
		if err != nil {
			return err
		}
		s.LastChanged = lc
	}

	luStr, ok := resMap["last_updated"].(string)
	if !ok {
		return errors.New("unable to cast last_updated into string")
	}
	lu, err := time.Parse(time.RFC3339, luStr)
	if err != nil {
		return err
	}

	oldStateMap, ok := resMap["old_state"].(map[string]interface{})
	if !ok {
		return errors.New("unable to cast old_state to map[string]interface{}")
	}
	oldStateId, ok := oldStateMap["unique_id"].(int8)
	if !ok {
		return errors.New("unable to cast old_state unique_id into int8")
	}

	s.UniqueId = int(uniqueId)
	s.InstanceUniqueId = int(instanceUniqueId)
	s.State = state
	s.Attributes = attributes
	s.LastUpdated = lu
	s.OldStateId = int(oldStateId)

	return nil
}
