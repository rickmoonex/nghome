package statemachine

import (
	"encoding/json"
	"fmt"

	"github.com/rickmoonex/nghome/internal/system/database"
)

// StateMachine
type StateMachine struct{}

// GetLastState takes in an instance id and returns the last known state of that instance
func (s *StateMachine) GetLastState(instanceId string) (*State, error) {
	vars := map[string]interface{}{
		"instance_id": instanceId,
	}

	query := `
    instance = .instance_registry.get_instance_by_id(instance_id);
    .state_machine.get_last_state(instance);
  `

	client, err := database.GetClient()
	if err != nil {
		return nil, err
	}

	var state State

	res, err := client.Query("//user_space", query, vars)
	if err != nil {
		return nil, err
	}

	if err := state.FromInterface(res); err != nil {
		return nil, err
	}

	return &state, nil
}

// AddEntry adds a new entry to the state machine
func (s *StateMachine) AddEntry(instanceId, state string, attributes map[string]interface{}) (*State, error) {
	if attributes == nil {
		attributes = map[string]interface{}{}
	}
	attrJson, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"instance_id": instanceId,
		"state":       state,
		"attributes":  string(attrJson),
	}

	query := `
    wse();
    instance = .instance_registry.get_instance_by_id(instance_id);
    .state_machine.add_entry(instance, state, attributes);
  `

	client, err := database.GetClient()
	if err != nil {
		return nil, err
	}

	var newState State

	res, err := client.Query("//user_space", query, vars)
	if err != nil {
		return nil, err
	}

	if err := newState.FromInterface(res); err != nil {
		return nil, fmt.Errorf("error in from interface: %v", err)
	}

	return &newState, nil
}
