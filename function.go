package main

import (
	"github.com/jinzhu/gorm"
)

type Function struct {
	Persistent
	Name       string      `json:"name" gorm:"unique"`
	Identifier string      `json:"identifier" gorm:"unique"`
	Module     string      `json:"module"`
	Payload    interface{} `json:"payload" gorm:"type:varchar"`
}

func (f *Function) BeforeCreate(tx *gorm.DB) (err error) {

	return nil
}
