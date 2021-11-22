// Copyright (c) 2021 Braden Nicholson

package types

import "gorm.io/gorm"

var db *gorm.DB

func Load(database *gorm.DB) {
	db = database
}
