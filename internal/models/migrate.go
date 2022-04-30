// Copyright (c) 2022 Braden Nicholson

package models

import "udap/internal/store"

func MigrateModels() error {
	err := store.DB.AutoMigrate(Log{}, Endpoint{}, Entity{}, Module{}, Device{}, Network{}, Zone{})
	if err != nil {
		return err
	}
	return nil
}
