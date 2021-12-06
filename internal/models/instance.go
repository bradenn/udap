// Copyright (c) 2021 Braden Nicholson

package models

// // Instance is a subclass of Module. Instance carries instance-related environment information.
// type Instance struct {
// 	store.Persistent
// 	// Name refers to the name of the instance, not the name of the module.
// 	Name string `json:"name" gorm:"unique"`
// 	// Description briefly describe the nature of the instance, not the module.
// 	Description string `json:"description"`
// 	// ModuleId refers to the instances module
// 	ModuleId string `json:"moduleId"`
// 	// Config holds instance related environment information in JSON format.
// 	Config string `json:"config"`
//
// 	buffer string `gorm:"-"`
// }
//
// func (i *Instance) watchBuffer() {
// 	// for s := range i.bufferChannel {
// 	// 	if i.listener == nil {
// 	// 		return
// 	// 	}
// 	// 	i.listener <- UpdateBuffer{
// 	// 		InstanceId: i.Id,
// 	// 		Data:       s,
// 	// 	}
// 	// }
// 	// close(i.bufferChannel)
// }
//
// func GetInstance(instanceId string) (instance *Instance, err error) {
// 	if instances[instanceId] != nil {
// 		return instances[instanceId], nil
// 	}
// 	err = store.DB.Model(&Instance{}).Where("id = ?", instanceId).First(&instance).Error
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	instances[instanceId] = instance
// 	log.Info("Instance '%s' loaded.", instance.Name)
// 	return instance, nil
// }
//
// type ModifyInstance struct {
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Config      string `json:"config"`
// }
//
// type CreateInstance struct {
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	ModuleId    string `json:"moduleId"`
// }
//
// func (i *Instance) Entities() (entities []Entity, err error) {
// 	err = store.DB.Model(&Entity{}).Where("instance_id = ?", i.Id).Find(&entities).Error
// 	return entities, err
// }
//
// func (i *Instance) Create(body json.RawMessage) (err error) {
// 	instanceBody := CreateInstance{}
// 	err = json.Unmarshal(body, &instanceBody)
// 	if err != nil {
// 		return err
// 	}
// 	i.Name = instanceBody.Name
// 	i.Description = instanceBody.Description
// 	i.ModuleId = instanceBody.ModuleId
// 	err = store.DB.Model(i).Create(i).Error
// 	if err != nil {
// 		return err
// 	}
// 	store.DB.Model(i).Preload("Module").Where("id = ?", i.Id).First(i)
//
// 	return nil
// }
//
// func (i *Instance) Modify(body json.RawMessage) (err error) {
// 	instanceBody := ModifyInstance{}
// 	err = json.Unmarshal(body, &instanceBody)
// 	if err != nil {
// 		return err
// 	}
// 	i.Name = instanceBody.Name
// 	i.Description = instanceBody.Description
// 	err = store.DB.Model(i).Save(i).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (i *Instance) Reset() (err error) {
// 	i.Config = ""
// 	err = store.DB.Model(i).Save(i).Error
// 	if err != nil {
// 		return err
// 	}
// 	store.DB.Model(i).Preload("Module").Where("id = ?", i.Id).First(i)
//
// 	// component, err := i.Module.rawComponent()
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// err = component.Create(NewAgent(i.Id.String(), nil))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// // func (i *Instance) Delete(endpoint endpoint.Endpoint) (err error) {
// // 	err = db.DB.Model(&endpoint).Association("Instances").Delete(i)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	err = db.DB.Model(i).Delete(i).Error
// // 	if err != nil {
// // 		return err
// // 	}
// // 	return nil
// // }
//
// func (i *Instance) Instantiate(instanceId string) (err error) {
// 	err = store.DB.Model(i).Preload("Module").Where("id = ?", instanceId).First(i).Error
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }
//
// func (i *Instance) Run(data string) error {
//
// 	return nil
// }
//
// func (i *Instance) AfterFind(_ *gorm.DB) error {
// 	instances[i.Id] = i
// 	return nil
// }
//
// func (i *Instance) BeforeCreate(_ *gorm.DB) error {
//
// 	return nil
// }
//
// func (i *Instance) Save() {
// 	err := store.DB.Model(&Instance{}).Where("id = ?", i.Id).Save(i).Error
// 	if err != nil {
// 		return
// 	}
// }
