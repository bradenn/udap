// Copyright (c) 2022 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	ics "github.com/arran4/golang-ical"
	"net/http"
	"os"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module CalDav

type CalDav struct {
	plugin.Module
	entityId string
	receiver chan domain.Attribute
}

func (c *CalDav) mux() {
	for range c.receiver {

	}
}

func (c *CalDav) Setup() (plugin.Config, error) {

	err := c.UpdateInterval(time.Minute * 5)
	if err != nil {
		return plugin.Config{}, err
	}
	return c.Config, nil
}

func (c *CalDav) Run() error {
	err := c.InitConfig("server", "https://example.com")
	if err != nil {
		c.ErrF("%s", err)
	}
	c.receiver = make(chan domain.Attribute, 1)
	entity := &domain.Entity{
		Name:   "caldav",
		Type:   "media",
		Module: c.Config.Name,
	}

	err = c.Entities.Register(entity)
	if err != nil {
		return err
	}

	c.entityId = entity.Id

	media := domain.Attribute{
		Value:     "{}",
		Updated:   time.Now(),
		Request:   "{}",
		Requested: time.Now(),
		Entity:    c.entityId,
		Key:       "calendar",
		Type:      "media",
		Order:     0,
		Channel:   c.receiver,
	}

	err = c.Attributes.Register(&media)
	if err != nil {
		return err
	}

	return nil
}

func (c *CalDav) Update() error {
	if c.Ready() {
		c.poll()
	}
	return nil
}

func init() {

	configVariables := []plugin.Variable{
		{
			Name:    "server",
			Default: "https://example.com",
			Description: "The module will poll calendar info from any CalDAV server. " +
				"Ensure the address begins with 'https://'.",
		},
	}

	config := plugin.Config{
		Name:        "caldav",
		Type:        "module",
		Description: "Get calendar events from any caldav interface",
		Version:     "0.1.0",
		Author:      "Braden Nicholson",
		Variables:   configVariables,
	}

	Module.Config = config
}

type Event struct {
	Description string    `json:"description"`
	Summary     string    `json:"summary"`
	Location    string    `json:"location"`
	Rule        string    `json:"rule"`
	Start       time.Time `json:"start"`
	Days        string    `json:"days"`
	End         time.Time `json:"end"`
}

func (c *CalDav) poll() {
	link := os.Getenv("calDav")
	cli := http.Client{}
	get, err := cli.Get(link)
	defer get.Body.Close()
	if err != nil {
		return
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(get.Body)
	if err != nil {
		return
	}
	cal, err := ics.ParseCalendar(strings.NewReader(buf.String()))
	if err != nil {
		return
	}
	var events []Event
	for _, event := range cal.Events() {
		start, err := event.GetStartAt()
		if err != nil {
			return
		}
		end, err := event.GetEndAt()
		if err != nil {
			return
		}

		desc := event.GetProperty(ics.ComponentPropertyDescription)
		descValue := ""
		if desc != nil {
			descValue = desc.Value
		}
		summary := event.GetProperty(ics.ComponentPropertySummary)
		summaryValue := ""
		if summary != nil {
			summaryValue = summary.Value
		}
		loc := event.GetProperty(ics.ComponentPropertyLocation)
		locValue := ""
		if loc != nil {
			locValue = loc.Value
		}
		rule := event.GetProperty(ics.ComponentPropertyRrule)
		ruleValue := ""
		if rule != nil {
			ruleValue = rule.Value
		}
		ev := Event{
			Description: descValue,
			Summary:     summaryValue,
			Location:    locValue,
			Rule:        ruleValue,
			Start:       start,
			End:         end,
			Days:        "",
		}
		tokens := strings.Split(ruleValue, ";")
		for _, token := range tokens {
			keys := strings.Split(token, "=")
			switch keys[0] {
			case "BYDAY":
				ev.Days = keys[1]
			default:
				continue
			}
		}
		events = append(events, ev)

	}
	marshal, err := json.Marshal(events)
	if err != nil {
		return
	}
	err = c.Attributes.Update(c.entityId, "calendar", string(marshal), time.Now())
	if err != nil {
		return
	}

}
