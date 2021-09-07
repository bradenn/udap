package main

type Entity struct {
	Persistent
	Name        string     `json:"name" gorm:"unique"`
	Description string     `json:"description"`
	Functions   []Function `json:"functions" gorm:"many2many:entityFunction;"`
}

type Group struct {
	Persistent
	Name       string   `json:"name"  gorm:"unique"`
	Entities   []Entity `json:"entities" gorm:"many2many:entityGroup;"`
	Identifier string   `json:"identifier"  gorm:"unique"`
}
