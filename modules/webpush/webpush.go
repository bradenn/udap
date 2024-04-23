// Copyright (c) 2023 Braden Nicholson

package main

import (
	"encoding/json"
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"os"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module WebPush

func init() {
	config := plugin.Config{
		Name:        "webpush",
		Type:        "module",
		Description: "Push Notifications for web apps",
		Version:     "1.0.0",
		Author:      "Braden Nicholson",
		Interval:    time.Second * 5,
	}
	Module.Config = config
}

type WebPush struct {
	plugin.Module
	entityId string
	mux      chan domain.Attribute
}

func (w *WebPush) Setup() (plugin.Config, error) {
	return w.Config, nil
}

func (w *WebPush) Update() error {
	return nil
}

func (w *WebPush) SendNotification(message string, endpoint *domain.Endpoint) {
	s := &webpush.Subscription{}

	fmt.Println(endpoint.Push)
	err := json.Unmarshal([]byte(endpoint.Push), s)
	if err != nil {
		fmt.Println("JSON:", err.Error())
		return
	}
	//https://web.push.apple.com/QJUmGfXB5CoQheyvUdpB579jfQRTCVcMEWUyuSYTO-cZtQXmbS2Tk4pOWJMVwqhFUq4r2twm23lpEJwPpSFlyahxaUEdC8M0OFS1XYepVezanOooEr917JPYMuBvvma2keQ9TNYZ3_BuYC60W8h8v5mDdyeK0WSmpPWAOKq1Pe0
	//https://web.push.apple.com/QJUmGfXB5CoQheyvUdpB579jfQRTCVcMEWUyuSYTO-cZtQXmbS2Tk4pOWJMVwqhFUq4r2twm23lpEJwPpSFlyahxaUEdC8M0OFS1XYepVezanOooEr917JPYMuBvvma2keQ9TNYZ3_BuYC60W8h8v5mDdyeK0WSmpPWAOKq1Pe0

	private, privateOk := os.LookupEnv("webpushPrivate")
	if !privateOk {
		return
	}
	public, publicOk := os.LookupEnv("webpushPublic")
	if !publicOk {
		return
	}
	// Send Notification
	resp, err := webpush.SendNotification([]byte(message), s, &webpush.Options{
		Subscriber:      "notify@udap.app", // Do not include "mailto:"
		VAPIDPublicKey:  public,
		VAPIDPrivateKey: private,
		TTL:             30,
		Urgency:         "0",
	})

	if err != nil {
		fmt.Println("VAPID:", err.Error())
		return
	}
	defer resp.Body.Close()
}

type NotificationRequest struct {
	EndpointId string `json:"endpointId"`
	Title      string `json:"title"`
	Message    string `json:"message"`
}

func (w *WebPush) handleNotifyRequest(attribute domain.Attribute) {
	nr := NotificationRequest{}
	err := json.Unmarshal([]byte(attribute.Request), &nr)
	if err != nil {
		return
	}
	entity, err := w.Endpoints.FindById(nr.EndpointId)
	if err != nil {
		return
	}

	marshal, err := json.Marshal(nr)
	if err != nil {
		return
	}

	w.SendNotification(string(marshal), entity)

}

func (w *WebPush) handleMux() {
	for {
		select {
		case attribute := <-w.mux:
			go w.handleNotifyRequest(attribute)
		}
	}
}

func (w *WebPush) Run() error {
	w.mux = make(chan domain.Attribute, 10)

	go w.handleMux()

	entity := domain.Entity{
		Name:   "webpush",
		Type:   "media",
		Module: "webpush",
	}

	err := w.Entities.Register(&entity)
	if err != nil {
		return err
	}

	w.entityId = entity.Id

	attribute := domain.Attribute{
		Entity:  entity.Id,
		Value:   "",
		Request: "",
		Key:     "notify",
		Type:    "media",
		Order:   0,
		Channel: w.mux,
	}

	err = w.Attributes.Register(&attribute)
	if err != nil {
		return err
	}

	return nil
}

func (w *WebPush) Dispose() error {
	return nil
}
