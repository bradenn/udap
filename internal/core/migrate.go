// Copyright (c) 2022 Braden Nicholson

package core

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func MigrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(domain.Attribute{}, domain.Entity{}, domain.Module{}, domain.Device{}, domain.Endpoint{},
		domain.User{}, domain.Network{}, domain.Zone{}, domain.Notification{}, domain.Macro{}, domain.Trigger{},
		domain.SubRoutine{}, domain.Action{}, domain.AttributeLog{})
	if err != nil {
		return err
	}
	return nil
}
