package instanceregistry

import (
	"errors"
	"time"
)

// InstanceType represent the types of instances available
type InstanceType string

const (
	InstanceTypeSwitch InstanceType = "switch"
	InstanceTypeLight  InstanceType = "light"
	InstanceTypeSensor InstanceType = "sensor"
)

// Instance repesents an entry into the instance registry
type Instance struct {
	UniqueId     int
	Type         InstanceType
	InstanceId   string
	CreatedAt    time.Time
	ModifiedAt   time.Time
	FriendlyName string
}

// FromInterface decodes the interface{} response from ThingsDB into the struct
func (i *Instance) FromInterface(res interface{}) error {
	resMap, ok := res.(map[string]interface{})
	if !ok {
		return errors.New("unable to cast res to map[string]interface")
	}

	uniqueId, ok := resMap["unique_id"].(int8)
	if !ok {
		return errors.New("unable to cast unique_id into int8")
	}

	typeStr, ok := resMap["type"].(string)
	if !ok {
		return errors.New("unable to cast type into string")
	}
	instanceType := InstanceType(typeStr)

	instanceId, ok := resMap["instance_id"].(string)
	if !ok {
		return errors.New("unable to cast instance_id to string")
	}

	caStr, ok := resMap["created_at"].(string)
	if !ok {
		return errors.New("unable to cast created_at to string")
	}
	createdAt, err := time.Parse(time.RFC3339, caStr)
	if err != nil {
		return err
	}

	maStr, ok := resMap["modified_at"].(string)
	if !ok {
		return errors.New("unable to cast modified_at to string")
	}
	modifiedAt, err := time.Parse(time.RFC3339, maStr)
	if err != nil {
		return err
	}

	friendlyName, ok := resMap["friendly_name"].(string)
	if !ok {
		return errors.New("unable to cast friendly_name to string")
	}

	i.UniqueId = int(uniqueId)
	i.Type = instanceType
	i.InstanceId = instanceId
	i.CreatedAt = createdAt
	i.ModifiedAt = modifiedAt
	i.FriendlyName = friendlyName

	return nil
}
