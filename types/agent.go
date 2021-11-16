package types

import (
	"encoding/json"
)

type Agent struct {
	instanceId string
	channel    chan UpdateBuffer
}

func NewAgent(instanceId string, channel chan UpdateBuffer) Agent {
	agent := Agent{
		instanceId: instanceId,
		channel:    channel,
	}
	return agent
}

func (h *Agent) InstanceId() string {
	return h.instanceId
}

func (h *Agent) Update(data interface{}) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	proto := string(marshal)
	h.channel <- UpdateBuffer{InstanceId: h.instanceId, Data: proto}
	return nil
}

func (h *Agent) Store(data interface{}) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = db.Model(&Instance{}).Where("id = ?", h.instanceId).Update("Config", string(marshal)).Error
	if err != nil {
		return err
	}
	return nil
}

func (h *Agent) Retrieve() (string, error) {
	var instance Instance
	err := db.Model(&Instance{}).Where("id = ?", h.instanceId).First(&instance).Error
	if err != nil {
		return "", err
	}

	return instance.Config, nil
}

func (h *Agent) Register(entity Entity) (err error) {
	err = NewEntity(entity.Name, entity.Type, h.instanceId)
	if err != nil {
		return err
	}
	return nil
}

func (h *Agent) Ready() bool {
	return h.instanceId == ""
}
