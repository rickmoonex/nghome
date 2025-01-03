package helper

import (
	"github.com/rickmoonex/nghome/internal/system/instanceregistry"
	"github.com/rickmoonex/nghome/internal/system/statemachine"
)

// NGContext is passed into every instance struct and provides access to system context
type NGContext struct {
	StateMachine     *statemachine.StateMachine
	InstanceRegistry *instanceregistry.InstanceRegistry
}

func NewNGContext() *NGContext {
	sm := &statemachine.StateMachine{}
	ir := &instanceregistry.InstanceRegistry{}

	return &NGContext{StateMachine: sm, InstanceRegistry: ir}
}
