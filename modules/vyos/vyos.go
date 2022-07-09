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
	"os/exec"
	"strings"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
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
	err := v.UpdateInterval(30000)
	if err != nil {
		return plugin.Config{}, err
	}
	return v.Config, nil
}

func (v *Vyos) Update() error {
	if v.Ready() {
		devices, err := v.Devices.FindAll()
		if err != nil {
			return err
		}
		if len(*devices) < 0 {
			return nil
		}
		for _, device := range *devices {
			go func(device domain.Device) {
				ctx := context.Background()
				timeout, cancelFunc := context.WithTimeout(ctx, time.Second*1)
				defer cancelFunc()
				cmd := exec.CommandContext(timeout, "/bin/zsh", "-c",
					fmt.Sprintf("ping -c 1 %s > /dev/null && echo true || echo false", device.Ipv4))
				start := time.Now()
				output, _ := cmd.Output()
				if string(output) == "true\n" {
					device.Latency = time.Since(start)
					device.LastSeen = time.Now()
					device.State = "ONLINE"
					err = v.Devices.Update(&device)
					if err != nil {
						fmt.Println(err)
						return
					}
				} else {
					device.Latency = 0
					device.State = "OFFLINE"
					err = v.Devices.Update(&device)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

			}(device)

		}
	}
	return nil
}

func (v *Vyos) Run() error {
	go func() {
		err := v.fetchNetworks()
		if err != nil {
			log.Err(err)
		}
	}()
	return nil
}

func (v *Vyos) arpScan(network domain.Network) error {
	cmd := exec.Command("/bin/zsh", "-c", "arp -an | awk -F'[ ()]' '{OFS=\";\"; print $1,$3,$6}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	out := string(output)
	devices := strings.Split(out, "\n")
	for _, str := range devices {
		args := strings.Split(str, ";")
		if len(args) < 3 {
			continue
		}
		if args[2] != "" {
			device := domain.Device{}
			device.Ipv4 = args[1]
			device.Mac = strings.ToLower(args[2])
			device.NetworkId = network.Id
			device.LastSeen = time.Now()
			err = v.Devices.Register(&device)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Vyos) scanSubnet(network domain.Network) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithConnectScan(),
		nmap.WithTargets("10.0.1.0/24"),
		nmap.WithInterface("en1"),
		nmap.WithContext(ctx),
	)

	if err != nil {
		return fmt.Errorf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		return fmt.Errorf("network scan failed: %v", err)
	}

	if warnings != nil {
		return fmt.Errorf("Warnings: \n %v", warnings)
	}
	// Use the results to print an example output
	for _, host := range result.Hosts {

		device := domain.Device{}

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

	}
	return nil
}

type Response struct {
	Success bool        `json:"success"`
	Data    Dhcp        `json:"data"`
	Error   interface{} `json:"error"`
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
	DefaultRouter string        `json:"default-index"`
	Lease         string        `json:"lease"`
	Range         map[int]Range `json:"range"`
}

func (v *Vyos) fetchNetworks() error {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	val := url.Values{}
	val.Set("key", "toor")
	val.Set("data", "{\"op\": \"showConfig\", \"path\": [\"service\", \"dhcp-server\"]}")

	request, err := client.PostForm("https://10.0.1.1:8005/retrieve", val)
	if err != nil {
		return fmt.Errorf("vyos index is non-existance")
	}

	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(request.Body)
	if err != nil {
		return err
	}

	payload := Response{}
	err = json.Unmarshal(buffer.Bytes(), &payload)
	if err != nil {
		return err
	}
	d := payload.Data

	wg := sync.WaitGroup{}

	for name, lan := range d.Networks {
		network := domain.Network{}
		network.Name = name
		network.Dns = strings.Join(lan.NameServer, ",")
		for s, subnet := range lan.Subnets {
			network.Mask = s
			network.Router = subnet.DefaultRouter
			network.Lease = subnet.Lease
			network.Range = fmt.Sprintf("%s-%s", subnet.Range[0].Start, subnet.Range[0].Stop)
			break
		}

		err = v.Networks.Register(&network)
		if err != nil {
			log.Err(err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			err = v.arpScan(network)
			if err != nil {
				log.Err(err)
			}
		}()

	}

	wg.Wait()
	return nil
}
