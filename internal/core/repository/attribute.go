// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/srv/store"
)

type attributeRepo struct {
	generic.Store[domain.Attribute]
	cache map[string]domain.Attribute
	mutex sync.Mutex
	store *store.Store
	db    *gorm.DB
}

func (u *attributeRepo) FindRecent() (*[]domain.Attribute, error) {
	var logs []domain.Attribute
	err := u.db.Model(&domain.Attribute{}).Order("updated desc").Limit(100).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return &logs, nil
}

func (u *attributeRepo) FindRecentLogs() (*[]domain.AttributeLog, error) {
	var logs []domain.AttributeLog
	err := u.db.Model(&domain.AttributeLog{}).Limit(100).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return &logs, nil
}

func NewAttributeRepository(db *gorm.DB, str *store.Store) ports.AttributeRepository {
	return &attributeRepo{
		db:    db,
		store: str,
		Store: generic.NewStore[domain.Attribute](db),
		mutex: sync.Mutex{},
		cache: make(map[string]domain.Attribute),
	}
}

//
//func (u *attributeRepo) Log(attribute *domain.Attribute) error {
//	if attribute.Key == "sensor" {
//		return nil
//	}
//
//	payload := map[string]json.RawMessage{}
//
//	err := json.Unmarshal([]byte(attribute.Value), &payload)
//	if err != nil {
//		return nil
//	}
//
//	for key, data := range payload {
//		fl := 0.0
//		err = json.Unmarshal(data, &fl)
//		if err != nil {
//			continue
//		}
//		rootKey := fmt.Sprintf("%s-%s", attribute.Id, key)
//		err = u.store.Push(rootKey, fl)
//		if err != nil {
//			continue
//		}
//	}
//}
func (u *attributeRepo) Summary(key string, start int64, stop int64, window int, mode string) (map[int64]float64, error) {
	return u.store.Summary(key, start, stop, window, mode)
}

func (u *attributeRepo) Log(attribute *domain.Attribute) error {

	if attribute.Key != "sensor" {
		return nil
	}

	payload := map[string]json.RawMessage{}
	//e := domain.Entity{}
	//u.db.Model(domain.Entity{}).Where("id = ?", attribute.Entity).First(&e)
	//
	//if e.Alias != "" {
	//	attribute.Serial = e.Alias
	//} else {
	//	attribute.Serial = e.Name
	//}
	//
	//u.Update(attribute)

	err := json.Unmarshal([]byte(attribute.Value), &payload)
	if err != nil {
		return nil
	}

	for key, data := range payload {
		fl := 0.0
		err = json.Unmarshal(data, &fl)
		if err != nil {
			continue
		}
		rootKey := fmt.Sprintf("%s-%s", attribute.Id, key)
		d := map[string]string{"type": "sensor", "class": key, "entity": attribute.Entity, "model": attribute.Type, "serial": attribute.Serial}
		err = u.store.Push(rootKey, d, fl)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return nil
}

func (u *attributeRepo) StateUpdate(attribute *domain.Attribute) error {
	composite := fmt.Sprintf("%s.%s", attribute.Entity, attribute.Key)

	u.mutex.Lock()
	val, ok := u.cache[composite]
	u.mutex.Unlock()

	if !ok {
		u.cache[composite] = *attribute
	}

	if val.Value != attribute.Value {
		err := u.Store.Update(attribute)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *attributeRepo) CachedValue(entity string, key string) (*domain.Attribute, error) {
	composite := fmt.Sprintf("%s.%s", entity, key)

	u.mutex.Lock()
	val, ok := u.cache[composite]
	u.mutex.Unlock()

	if ok {
		return &val, nil
	}

	var target domain.Attribute
	if err := u.db.Model(&domain.Attribute{}).Where("entity = ? AND key = ?", entity,
		key).First(&target).Error; err != nil {
		return nil, err
	}

	u.cache[composite] = target

	return &target, nil
}

func (u *attributeRepo) Register(attribute *domain.Attribute) error {
	//serial := attribute.Serial

	err := u.db.Model(&domain.Attribute{}).Where("entity = ? AND key = ?", attribute.Entity, attribute.Key).FirstOrCreate(attribute).Error
	if err != nil {
		return err
	}

	//if attribute.Serial != serial {
	//	attribute.Serial = serial
	//	u.Update(attribute)
	//}
	return nil
}

func (u *attributeRepo) FindByComposite(entity string, key string) (*domain.Attribute, error) {
	var target domain.Attribute

	if err := u.db.Model(&domain.Attribute{}).Where("entity = ? AND key = ?", entity,
		key).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (u *attributeRepo) FindAllByEntity(entity string) (*[]domain.Attribute, error) {
	var target []domain.Attribute
	if err := u.db.Where("entity = ?", entity).Find(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}
