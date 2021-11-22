// Copyright (c) 2021 Braden Nicholson

package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"udap/module"
)

func (e *Endpoint) EndpointRequest(request Request, c *websocket.Conn) (err error) {

	var reqBody requestBody
	err = json.Unmarshal(request.Body, &reqBody)
	if err != nil {
		return err
	}
	switch t := request.Operation; t {
	case METADATA:
		err = e.Metadata()
		if err != nil {
			return err
		}
	case ENROLL:
		err = e.Enroll(*e, c, request.Body)
		if err != nil {
			return err
		}
	case SUBSCRIBE:
		err = e.Subscribe(reqBody)
		if err != nil {
			return err
		}
	}
	return nil
}

type IdentifierBody struct {
	Id string `json:"id"`
}

func (e *Endpoint) InstanceRequest(request Request) (err error) {
	body := IdentifierBody{}
	instance := module.Instance{}

	err = json.Unmarshal(request.Body, &body)
	if err != nil {
		return err
	}

	if body.Id == "" {
		switch t := request.Operation; t {
		case CREATE:
			err = instance.Create(request.Body)
			if err != nil {
				return err
			}
			err = request.Sender.GrantInstance(instance)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid operation '%s'; you may have forgotten to provice an Id", t)
		}
		return nil
	}

	instance, err = request.Sender.GetInstance(body.Id)
	if instance.ModuleId.String() == "" {
		panic("Wtf? No Module ID on this bad boy")
	}
	instance.Module, err = module.Get(instance.ModuleId.String())
	if err != nil {
		fmt.Println(err)
	}
	switch t := request.Operation; t {
	case MODIFY:
		err = instance.Modify(request.Body)
		if err != nil {
			return err
		}
	case RUN:
		err = instance.Run(string(request.Body))
		if err != nil {
			return err
		}
	case RESET:
		err = instance.Reset()
		if err != nil {
			return err
		}
	case DELETE:
		err = instance.Reset()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid operation '%s'", t)
	}

	return nil
}
