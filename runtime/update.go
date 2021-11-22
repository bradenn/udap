// Copyright (c) 2021 Braden Nicholson

package runtime

import (
	"encoding/json"
	"udap/module"
)

func (r *Runtime) Listen() {
	for buffer := range r.updater {
		r.updateInstance(buffer)

	}
	close(r.updater)
}

type Update struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"data"`
	Id   string `json:"id"`
}

func (r *Runtime) push() {
	// for _, endpoint := range r.endpoints {
	// 	// endpoint.Push()
	// }
}

func (r *Runtime) updateInstance(buffer module.UpdateBuffer) {

	i := r.instances[buffer.InstanceId]
	update := Update{
		Type: i.Module.Path,
		Name: i.Name,
		Data: buffer.Data,
		Id:   i.Id,
	}

	marshal, err := json.Marshal(update)
	if err != nil {
		return
	}
	r.instances[buffer.InstanceId].Buffer = string(marshal)
}
