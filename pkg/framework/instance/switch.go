package instance

import (
	"errors"

	"github.com/rickmoonex/nghome/internal/system/instanceregistry"
	"github.com/rickmoonex/nghome/pkg/framework/helper"
)

const (
	SwitchInstanceStateOn  InstanceState = "on"
	SwitchInstanceStateOff InstanceState = "off"
)

type SwitchInstanceInterface interface {
	BaseInstance

	TurnOn()
	TurnOff()
	Toggle()
}

// SwitchInstance implements all logic for a switch instance
type SwitchInstance struct {
	BaseInstance

	isOn bool
}

func (s *SwitchInstance) Init(ctx *helper.NGContext, config map[string]interface{}) error {
	s.instanceType = instanceregistry.InstanceTypeSwitch

	name, ok := config["name"]
	if !ok {
		return errors.New("no name")
	}
	s.name, ok = name.(string)
	if !ok {
		return errors.New("name cast")
	}

	return s.BaseInstance.Init(ctx, config)
}

func (s *SwitchInstance) TurnOn() {
	s.state = SwitchInstanceStateOn

	go s.writeState()
}

func (s *SwitchInstance) TurnOff() {
	s.state = SwitchInstanceStateOff

	go s.writeState()
}

func (s *SwitchInstance) Toggle() {
	if s.isOn {
		go s.TurnOff()
	}
	go s.TurnOn()
}
