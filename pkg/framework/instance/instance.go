package instance

import (
	"errors"
	"fmt"

	"github.com/rickmoonex/nghome/internal/system/instanceregistry"
	"github.com/rickmoonex/nghome/pkg/framework/helper"
)

type InstanceState string

const (
	InstanceStateUnknown     InstanceState = "unknown"
	InstanceStateUnavailable InstanceState = "unavailable"
)

type InstanceAttribute string

// BaseInstanceInterface needs to implemented by every instance
type BaseInstanceInterface interface {
	Init(ctx *helper.NGContext, config map[string]interface{}) error
}

// BaseInstance implements base functions of all instances
type BaseInstance struct {
	ctx    *helper.NGContext
	config map[string]interface{}

	instanceEntry *instanceregistry.Instance
	instanceType  instanceregistry.InstanceType

	instanceId string
	state      InstanceState
	attributes map[InstanceAttribute]interface{}
	name       string
}

// Init initializes the instance
func (b *BaseInstance) Init(ctx *helper.NGContext, config map[string]interface{}) error {
	b.ctx = ctx

	if b.instanceType == "" {
		return errors.New("instanceType is not set")
	}

	b.instanceId = fmt.Sprintf("%s.%s", string(b.instanceType), b.name)

	instance, err := b.ctx.InstanceRegistry.RegisterInstance(b.instanceType, b.instanceId, b.name)
	if err != nil {
		return err
	}
	b.instanceEntry = instance

	return nil
}

// writeState writes the instance's state to the state machine
func (b *BaseInstance) writeState() {
	state := string(b.state)

	attrMap := make(map[string]interface{})
	for k, v := range b.attributes {
		attrMap[string(k)] = v
	}

	go func() {
		b.ctx.StateMachine.AddEntry(b.instanceId, state, attrMap)
	}()
}
