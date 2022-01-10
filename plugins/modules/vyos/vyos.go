// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Ullaakut/nmap/v2"
	"net/http"
	"net/url"
	"strings"
	"time"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

var Module Vyos

type Vyos struct {
	plugin.Module
}

func init() {
	config := plugin.Config{
		Name:        "vyos",
		Type:        "module",
		Description: "Network interfaces controller",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (v *Vyos) Setup() (plugin.Config, error) {
	v.fetchNetworks()
	// Since ft232 is a shitty module, it panics when USB can't be found.

	return v.Config, nil
}

func (v *Vyos) Update() error {
	return nil
}

func (v *Vyos) Run() error {
	return nil
}

func (v *Vyos) scanSubnet(network models.Network) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithPrivileged(),
		nmap.WithConnectScan(),
		nmap.WithTargets("10.0.1.2-32"),
		nmap.WithContext(ctx),
	)

	if err != nil {
		return fmt.Errorf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		return fmt.Errorf("unable to run nmap scan: %v", err)
	}

	if warnings != nil {
		return fmt.Errorf("Warnings: \n %v", warnings)
	}
	// Use the results to print an example output
	for _, host := range result.Hosts {
		device := models.NewDevice()

		for _, addr := range host.Addresses {
			switch addr.AddrType {
			case "ipv4":
				device.Ipv4 = addr.String()
			case "ipv6":
				device.Ipv6 = addr.String()
			case "mac":
				device.Mac = addr.String()
			}
		}
		device.NetworkId = network.Id
		_, err = v.Send("device", "register", &device)
		if err != nil {

			return err
		}

	}
	return nil
}

type Response struct {
	Success bool `json:"success"`
	Data    Dhcp `json:"data"`
	Error   any  `json:"error"`
}

type Addr struct {
	Address     string `json:"address"`
	Description string `json:"description"`
}

type Payload struct {
	Ethernet map[string]Addr `json:"ethernet"`
}

type Dhcp struct {
	Networks map[string]Lan `json:"shared-network-name"`
}

type Lan struct {
	NameServer []string          `json:"name-server"`
	Subnets    map[string]Subnet `json:"subnet"`
}

type Range struct {
	Start string `json:"start"`
	Stop  string `json:"stop"`
}

type Subnet struct {
	DefaultRouter string        `json:"default-router"`
	Lease         string        `json:"lease"`
	Range         map[int]Range `json:"range"`
}

func (v *Vyos) fetchNetworks() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: transport,
	}

	val := url.Values{}
	val.Set("key", "toor")
	val.Set("data", "{\"op\": \"showConfig\", \"path\": [\"service\", \"dhcp-server\"]}")

	request, err := client.PostForm("https://10.0.1.1:8005/retrieve", val)
	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	payload := Response{}
	err = json.Unmarshal(buffer.Bytes(), &payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	d := payload.Data
	for name, lan := range d.Networks {
		network := models.Network{}
		network.Name = name
		network.Dns = strings.Join(lan.NameServer, ",")
		network.Dns = strings.Join(lan.NameServer, ",")
		for s, subnet := range lan.Subnets {
			network.Mask = s
			network.Router = subnet.DefaultRouter
			network.Lease = subnet.Lease
			network.Range = fmt.Sprintf("%s-%s", subnet.Range[0].Start, subnet.Range[0].Stop)
			break
		}

		_, err = v.Send("network", "register", &network)
		if err != nil {
			return
		}

		go func() {
			err = v.scanSubnet(network)
			if err != nil {
				log.Err(err)
				return
			}
		}()

	}
}
