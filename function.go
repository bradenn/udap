package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"udap/module"
)

type Function struct {
	Persistent
	Name       string      `json:"name" gorm:"unique"`
	Identifier string      `json:"identifier" gorm:"unique"`
	Module     string      `json:"module"`
	Payload    interface{} `json:"payload" gorm:"type:varchar"`
}

func (f *Function) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = module.Get(f.Module)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(f.Payload)
	if err != nil {
		return err
	}

	f.Payload = string(marshal)

	return nil
}

func (f *Function) Run(payload interface{}) interface{} {
	mod, err := module.Get(f.Module)
	if err != nil {
		return nil
	}
	return mod.Run(f.Identifier, payload)
}
