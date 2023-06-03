// Copyright (c) 2023 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module Proxmox

type Proxmox struct {
	plugin.Module
	monitorId string
}

func init() {
	config := plugin.Config{
		Name:        "proxmox",
		Type:        "module",
		Description: "A simple proxmox API module",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

type Network struct {
	Netout int64 `json:"netout"`

	Netin int64 `json:"netin"`
}

type Utilization struct {
	Mem    int64 `json:"mem"`
	Maxmem int64 `json:"maxmem"`

	Cpu    float64 `json:"cpu"`
	Cpus   int     `json:"cpus,omitempty"`
	Maxcpu int     `json:"maxcpu"`

	Disk    int64 `json:"disk"`
	Maxdisk int64 `json:"maxdisk"`
}

type VM struct {
	Vmid   int    `json:"vmid"`
	Status string `json:"status"`

	Name string `json:"name"`
	Pid  int    `json:"pid"`

	Uptime int `json:"uptime"`

	Utilization Utilization `json:"utilization"`
	Network     Network     `json:"network"`
}

type Node struct {
	Node            string      `json:"node"`
	Status          string      `json:"status"`
	Uptime          int         `json:"uptime"`
	Utilization     Utilization `json:"utilization"`
	VirtualMachines []VM        `json:"virtualMachines"`
}

type APINodeStruct struct {
	Node   string `json:"node"`
	Status string `json:"status"`
	Uptime int    `json:"uptime"`
	Utilization
}

type APIVMStruct struct {
	Vmid   int    `json:"vmid"`
	Status string `json:"status"`

	Name string `json:"name"`
	Pid  int    `json:"pid"`

	Uptime int `json:"uptime"`
	Network
	Utilization
}

type ProxmoxAPINodesRequest struct {
	Data []APINodeStruct `json:"data"`
}

type ProxmoxAPINodeRequest struct {
	Data []APIVMStruct `json:"data"`
}

type ProxmoxData struct {
	Nodes []Node `json:"nodes"`
}

func (p *Proxmox) Setup() (plugin.Config, error) {
	err := p.UpdateInterval(1000 * 1)
	if err != nil {
		return plugin.Config{}, err
	}
	return p.Config, nil
}

const APIFormat = "https://%s/api2/json%s"
const APIListNodes = "/nodes"
const APINode = "/nodes/%s"
const APIListQemu = "/nodes/%s/qemu"
const APIQemu = "/nodes/%s/qemu/%s"
const Authentication = "PVEAPIToken=root@pam!%s=%s"

func (p *Proxmox) fetchProxmoxData() (ProxmoxData, error) {
	request, err := p.authenticatedRequest(fmt.Sprintf(APIListNodes))
	if err != nil {
		return ProxmoxData{}, err
	}

	apiData := &ProxmoxAPINodesRequest{}
	err = json.Unmarshal(request, apiData)
	if err != nil {
		return ProxmoxData{}, err
	}
	pmd := ProxmoxData{
		Nodes: []Node{},
	}
	for _, node := range apiData.Data {
		nd := Node{
			Node:            node.Node,
			Status:          node.Status,
			Uptime:          node.Uptime,
			Utilization:     node.Utilization,
			VirtualMachines: []VM{},
		}
		request, err = p.authenticatedRequest(fmt.Sprintf(APIListQemu, node.Node))
		if err != nil {
			return ProxmoxData{}, err
		}

		apiNodeData := ProxmoxAPINodeRequest{}
		err = json.Unmarshal(request, &apiNodeData)
		if err != nil {
			return ProxmoxData{}, err
		}

		for _, vm := range apiNodeData.Data {
			nvm := VM{
				Vmid:        vm.Vmid,
				Status:      vm.Status,
				Name:        vm.Name,
				Pid:         vm.Pid,
				Uptime:      vm.Uptime,
				Utilization: vm.Utilization,
				Network:     vm.Network,
			}
			nvm.Utilization.Maxcpu = vm.Utilization.Cpus
			nd.VirtualMachines = append(nd.VirtualMachines, nvm)
		}

		pmd.Nodes = append(pmd.Nodes, nd)
	}

	return pmd, nil
}

func (p *Proxmox) authenticatedRequest(endpoint string) ([]byte, error) {
	hostname := ""
	found := false
	if hostname, found = os.LookupEnv("proxmoxHostname"); !found {
		return []byte{}, fmt.Errorf("hostname env not found")
	}
	baseUrl := fmt.Sprintf(APIFormat, hostname, endpoint)
	cli, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return []byte{}, err
	}

	username := ""
	token := ""
	if username, found = os.LookupEnv("proxmoxUsername"); !found {
		return []byte{}, fmt.Errorf("hostname env not found")
	}
	if token, found = os.LookupEnv("proxmoxToken"); !found {
		return []byte{}, fmt.Errorf("hostname env not found")
	}

	cli.Header.Set("Authorization", fmt.Sprintf(Authentication, username, token))

	client := http.Client{}
	defer client.CloseIdleConnections()
	response, err := client.Do(cli)

	if err != nil {
		return []byte{}, err
	}

	var res bytes.Buffer
	_, err = res.ReadFrom(response.Body)

	if err != nil {
		return []byte{}, err
	}

	_ = response.Body.Close()
	return res.Bytes(), nil
}

func (p *Proxmox) Update() error {
	if p.Ready() {
		data, err := p.fetchProxmoxData()
		if err != nil {
			return err
		}
		marshal, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = p.Attributes.Set(p.monitorId, "nodes", string(marshal))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Proxmox) Run() error {
	var err error
	p.monitorId, err = p.RegisterEntity("nodes", "monitor")
	if err != nil {
		return err
	}

	mux := make(chan domain.Attribute, 10)

	go func() {
		for range mux {

		}
	}()

	err = p.Attributes.Register(&domain.Attribute{
		Updated:   time.Time{},
		Requested: time.Time{},
		Value:     "{}",
		Request:   "{}",
		Entity:    p.monitorId,
		Key:       "nodes",
		Type:      "media",
		Order:     0,
		Channel:   mux,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxmox) Dispose() error {
	return nil
}
