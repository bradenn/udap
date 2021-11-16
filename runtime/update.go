package runtime

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"udap/types"
)

func (r *Runtime) Listen() {
	for buffer := range r.updater {
		r.updateInstance(buffer)
		r.push()
	}
	close(r.updater)
}

type Update struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"data"`
	Id   string `json:"id"`
}

func (r *Runtime) updateInstance(buffer types.UpdateBuffer) {

	i := r.instances[buffer.InstanceId]
	update := Update{
		Type: i.Module.Path,
		Name: i.Name,
		Data: buffer.Data,
		Id:   i.Id.String(),
	}

	marshal, err := json.Marshal(update)
	if err != nil {
		return
	}
	r.instances[buffer.InstanceId].Buffer = string(marshal)
}

// push sends data to all enrolled endpoints
func (r *Runtime) push() {
	for _, endpoint := range r.endpoints {
		response := Response{
			Status:    SUCCESS,
			Operation: UPDATE,
			Body:      map[string]interface{}{},
		}

		for _, instance := range endpoint.Subscriptions {
			response.Body[instance] = r.instances[instance].Buffer
		}

		marshal, err := json.Marshal(response)
		if err != nil {
			return
		}

		err = endpoint.Connection.WriteMessage(websocket.TextMessage, marshal)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
