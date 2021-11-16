package types

import (
	"github.com/google/uuid"
)

const (
	SWITCH = "switch"
	RGB    = "rgb"
	RGBW   = "rgbw"
	RGBCCT = "rgbcct"
)

type Entity struct {
	Persistent
	InstanceId uuid.UUID `json:"instanceId"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	State      string    `json:"state"`
}

func NewEntity(name string, entityType string, instanceId string) (err error) {
	parse, err := uuid.Parse(instanceId)
	if err != nil {
		return
	}
	entity := Entity{
		Name:       name,
		Type:       entityType,
		InstanceId: parse,
	}
	var cnt int64
	db.Model(&Entity{}).Where("instance_id = ? AND name = ?", instanceId, name).Count(&cnt)
	if cnt >= 1 {
		return nil
	}

	return db.Create(&entity).Error
}

func (e *Entity) SetState(state string) (err error) {

	return nil
}

type Device struct {
	Persistent
	InstanceId uuid.UUID `json:"instanceId"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
}
