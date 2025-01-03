package instanceregistry

import (
	"github.com/rickmoonex/nghome/internal/system/database"
)

// InstanceRegistry
type InstanceRegistry struct{}

// AddInstance adds a new instance to the registry
func (r *InstanceRegistry) RegisterInstance(instanceType InstanceType, instanceId string, friendlyName string) (*Instance, error) {
	client, err := database.GetClient()
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"type":          string(instanceType),
		"instance_id":   instanceId,
		"friendly_name": friendlyName,
	}

	var resInstance Instance

	res, err := client.Query("//user_space", `wse(); .instance_registry.register_instance(enum("InstanceType", type), instance_id, friendly_name);`, vars)
	if err != nil {
		return nil, err
	}

	if err := resInstance.FromInterface(res); err != nil {
		return nil, err
	}

	return &resInstance, nil
}

// GetInstanceById uses the InstanceId to lookup an item in the Instance Registry
func (r *InstanceRegistry) GetInstanceById(instanceId string) (*Instance, error) {
	client, err := database.GetClient()
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"instance_id": instanceId,
	}

	var resInstance Instance

	res, err := client.Query("//user_space", ".instance_registry.get_instance_by_id(instance_id);", vars)
	if err != nil {
		return nil, err
	}

	if err := resInstance.FromInterface(res); err != nil {
		return nil, err
	}

	return &resInstance, nil
}

// ChangeInstanceId uses to InstanceId to lookup an instance and update it's instance_id
func (r *InstanceRegistry) ChangeInstanceId(instanceId, newInstanceId string) (*Instance, error) {
	client, err := database.GetClient()
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"instance_id":     instanceId,
		"new_instance_id": newInstanceId,
	}

	var resInstance Instance

	res, err := client.Query("//user_space", ".instance_registry.change_instance_id(instance_id, new_instance_id;", vars)
	if err != nil {
		return nil, err
	}

	if err := resInstance.FromInterface(res); err != nil {
		return nil, err
	}

	return &resInstance, nil
}
