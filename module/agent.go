// Copyright (c) 2021 Braden Nicholson

package module

// // type Instance struct {
// //
// // 	// Name refers to the name of the instance, not the name of the module.
// // 	Name string `json:"name" gorm:"unique"`
// // 	// Description briefly describe the nature of the instance, not the module.
// // 	Description string `json:"description"`
// // 	// Module is the actual reference to the plugin module.
// // 	// Module   *module.Module `json:"module" gorm:"foreignKey:ModuleId"`
// //
// // 	// Config holds instance related environment information in JSON format.
// // 	Config      string `json:"config"`
// // 	Buffer      string
// // 	subscribers []string
// // }
// type Agent struct {
// 	instanceId string
// 	channel    chan UpdateBuffer
// }
//
// func NewAgent(instanceId string, channel chan UpdateBuffer) Agent {
// 	agent := Agent{
// 		instanceId: instanceId,
// 		channel:    channel,
// 	}
// 	return agent
// }
//
// func (h *Agent) InstanceId() string {
// 	return h.instanceId
// }
//
// func (h *Agent) Update(data interface{}) error {
// 	marshal, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
// 	proto := string(marshal)
// 	h.channel <- UpdateBuffer{InstanceId: h.instanceId, Data: proto}
// 	return nil
// }
//
// func (h *Agent) Store(data interface{}) error {
// 	marshal, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
// 	err = db.DB.Model(&types.Instance{}).Where("id = ?", h.instanceId).Update("Config", string(marshal)).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (h *Agent) Retrieve() (string, error) {
// 	var instance types.Instance
// 	err := db.DB.Model(&types.Instance{}).Where("id = ?", h.instanceId).First(&instance).Error
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return instance.Config, nil
// }
//
// func (h *Agent) Register(entity types.Entity) (err error) {
// 	// err = types.NewEntity(entity.Name, entity.Type, h.instanceId)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	return nil
// }
//
// func (h *Agent) Ready() bool {
// 	return h.instanceId == ""
// }
