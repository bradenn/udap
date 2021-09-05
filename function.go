package main

import (
	"github.com/jinzhu/gorm"
	"udap/module"
)

type Function struct {
	Persistent
	Name       string `json:"name" gorm:"unique"`
	Identifier string `json:"identifier" gorm:"unique"`
	Module     string `json:"module"`
}

func (f *Function) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = module.Get(f.Module)
	if err != nil {
		return err
	}
	return nil
}

func (f *Function) Run(payload interface{}) interface{} {
	mod, err := module.Get(f.Module)
	if err != nil {
		return nil
	}
	return mod.Run(f.Identifier, payload)
}
