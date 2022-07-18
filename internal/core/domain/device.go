// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"time"
	"udap/internal/core/domain/common"
)

type Utilization struct {
	Memory struct {
		Total uint64  `json:"total"`
		Used  float64 `json:"used"`
	} `json:"memory"`
	Network struct {
		Hostname string `json:"hostname"`
		Ipv4     string `json:"ipv4"`
		Mac      string `json:"mac"`
	} `json:"network"`
	Cpu struct {
		Cores int       `json:"cores"`
		Usage []float64 `json:"usage"`
	} `json:"cpu"`
	Disk struct {
		Total uint64  `json:"total"`
		Used  float64 `json:"used"`
	} `json:"disk"`
	Compute []Compute `json:"compute"`
}

type Compute struct {
	FanSpeed    int `json:"fanSpeed"`
	Temperature struct {
		Current  int `json:"current"`
		Throttle int `json:"throttle"`
		Target   int `json:"target"`
		Max      int `json:"max"`
	} `json:"temperature"`
	Utilization struct {
		GPU    int `json:"gpu"`
		Memory int `json:"memory"`
	} `json:"utilization"`
	Power struct {
		Draw float64 `json:"draw"`
		Max  float64 `json:"max"`
	} `json:"power"`
	Memory struct {
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
		Total    int `json:"total"`
	} `json:"memory"`
	Clocks struct {
		Graphics struct {
			Current int `json:"current"`
			Max     int `json:"max"`
		} `json:"graphics"`
		Streaming struct {
			Current int `json:"current"`
			Max     int `json:"max"`
		} `json:"streaming"`
		Memory struct {
			Current int `json:"current"`
			Max     int `json:"max"`
		} `json:"memory"`
		Video struct {
			Current int `json:"current"`
			Max     int `json:"max"`
		} `json:"video"`
	} `json:"clocks"`
	Processes []ComputeProcess `json:"processes"`
}

type ComputeProcess struct {
	Name   string `json:"name"`
	PID    string `json:"pid"`
	Memory int    `json:"memory"`
}

type Device struct {
	common.Persistent
	LastSeen    time.Time     `json:"lastSeen"`
	Latency     time.Duration `json:"latency"`
	State       string        `json:"state"`
	NetworkId   string        `json:"networkId" gorm:"-"`
	EntityId    string        `json:"entityId" gorm:"-"`
	Name        string        `json:"name"`
	Hostname    string        `json:"hostname"`
	IsQueryable bool          `json:"isQueryable" gorm:"default:false"`
	Utilization Utilization   `json:"utilization" gorm:"-"`
	Mac         string        `json:"mac"`
	Ipv4        string        `json:"ipv4"`
	Ipv6        string        `json:"ipv6"`
}

type DeviceRepository interface {
	common.Persist[Device]
}

type DeviceService interface {
	Observable
	FindAll() (*[]Device, error)
	FindById(id string) (*Device, error)
	Create(*Device) error
	FindOrCreate(*Device) error
	Register(*Device) error
	Update(*Device) error
	Delete(*Device) error
}
