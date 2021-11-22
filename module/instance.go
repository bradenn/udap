// Copyright (c) 2021 Braden Nicholson

package module

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"udap/types"
	"udap/udap"
	"udap/udap/db"
)

// Instance is a subclass of Module. Instance carries instance-related environment information.
type Instance struct {
	udap.Persistent
	// Name refers to the name of the instance, not the name of the module.
	Name string `json:"name" gorm:"unique"`
	// Description briefly describe the nature of the instance, not the module.
	Description string `json:"description"`
	// Module is the actual reference to the plugin module.
	Module   *Module   `json:"module" gorm:"foreignKey:ModuleId"`
	ModuleId uuid.UUID `json:"moduleId"`
	// Config holds instance related environment information in JSON format.
	Config      string `json:"config"`
	Buffer      string
	subscribers []string
}

func LoadInstance(id string) (instance *Instance, err error) {
	err = db.DB.Model(&Instance{}).Preload("Module").Where("id = ?", id).First(instance).Error
	udap.Info("Instance %s loaded.", id)
	return instance, err
}

type ModifyInstance struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Config      string `json:"config"`
}

type CreateInstance struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ModuleId    string `json:"moduleId"`
}

func (i *Instance) Entities() (entities []types.Entity, err error) {
	err = db.DB.Model(&types.Entity{}).Where("instance_id = ?", i.Id).Find(&entities).Error
	return entities, err
}

func (i *Instance) Create(body json.RawMessage) (err error) {
	instanceBody := CreateInstance{}
	err = json.Unmarshal(body, &instanceBody)
	if err != nil {
		return err
	}
	i.Name = instanceBody.Name
	i.Description = instanceBody.Description
	parse, err := uuid.Parse(instanceBody.ModuleId)
	if err != nil {
		return err
	}
	i.ModuleId = parse
	err = db.DB.Model(i).Create(i).Error
	if err != nil {
		return err
	}

	db.DB.Model(i).Preload("Module").Where("id = ?", i.Id).First(i)

	return nil
}

func (i *Instance) Modify(body json.RawMessage) (err error) {
	instanceBody := ModifyInstance{}
	err = json.Unmarshal(body, &instanceBody)
	if err != nil {
		return err
	}
	i.Name = instanceBody.Name
	i.Description = instanceBody.Description
	err = db.DB.Model(i).Save(i).Error
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Reset() (err error) {
	i.Config = ""
	err = db.DB.Model(i).Save(i).Error
	if err != nil {
		return err
	}
	db.DB.Model(i).Preload("Module").Where("id = ?", i.Id).First(i)

	// component, err := i.Module.rawComponent()
	// if err != nil {
	// 	return err
	// }
	// err = component.Create(NewAgent(i.Id.String(), nil))
	if err != nil {
		return err
	}
	return nil
}

// func (i *Instance) Delete(endpoint endpoint.Endpoint) (err error) {
// 	err = db.DB.Model(&endpoint).Association("Instances").Delete(i)
// 	if err != nil {
// 		return err
// 	}
// 	err = db.DB.Model(i).Delete(i).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (i *Instance) Subscribers() []string {
	return i.subscribers
}

func (i *Instance) Subscribe(endpointId string) {
	for _, subscriber := range i.subscribers {
		if endpointId == subscriber {
			return
		}
	}
	i.subscribers = append(i.subscribers, endpointId)
}

func (i *Instance) Unsubscribe(endpointId string) {
	for in, subscriber := range i.subscribers {
		if endpointId == subscriber {
			i.subscribers = append(i.subscribers[:in], i.subscribers[in+1:]...)
		}
	}
}

func (i *Instance) Instantiate(instanceId string) (err error) {
	err = db.DB.Model(i).Preload("Module").Where("id = ?", instanceId).First(i).Error
	if err != nil {
		return err
	}
	return err
}

func (i *Instance) Run(data string) error {
	err := i.Module.Run(data)
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) AfterFind(_ *gorm.DB) error {
	return nil
}

func (i *Instance) BeforeCreate(_ *gorm.DB) error {

	return nil
}

func (i *Instance) Save() {
	err := db.DB.Model(&Instance{}).Where("id = ?", i.Id).Save(i).Error
	if err != nil {
		return
	}
}
