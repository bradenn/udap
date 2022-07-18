// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Ullaakut/nmap/v2"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
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
	devicesIds   map[string]string
	pingQueue    chan domain.Device
	pingResolver chan domain.Device
	pingConn     *icmp.PacketConn
	done         chan bool
}

func init() {
	config := plugin.Config{
		Name:        "vyos",
		Type:        "module",
		Description: "Network interfaces controller",
		Version:     "1.0.2",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (v *Vyos) sendPing(ip string) error {
	// Resolve any DNS (if used) and get the real IP of the target
	dst, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		return err
	}

	// Make a new ICMP message
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: 0xffff, Seq: 0xffff,
			Data: []byte(fmt.Sprintf("%024d", time.Now().UnixNano())),
		},
	}

	b, err := m.Marshal(nil)
	if err != nil {
		return err
	}

	n, err := v.pingConn.WriteTo(b, dst)
	if err != nil {
		err = v.rejectPing(ip)
		if err != nil {
			return err
		}
	} else if n != len(b) {
		return fmt.Errorf("got %v; want %v", n, len(b))
	}

	return nil
}

func (v *Vyos) rejectPing(ipv4 string) error {
	id := v.devicesIds[ipv4]
	if id == "" {
		return nil
	}

	device, err := v.Devices.FindById(id)
	if err != nil {
		return err
	}

	device.Latency = 0
	device.State = "OFFLINE"

	err = v.Devices.Update(device)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vyos) resolvePing(ipv4 string, duration time.Duration) error {
	id := v.devicesIds[ipv4]
	if id == "" {
		return nil
	}

	device, err := v.Devices.FindById(id)
	if err != nil {
		return err
	}

	device.Latency = duration
	device.State = "ONLINE"
	device.LastSeen = time.Now()

	if device.IsQueryable {
		res := domain.Utilization{}
		res, err = v.queryDevice(device.Ipv4)
		if err != nil {
			log.Err(err)
			err = nil
		} else {
			device.Utilization = res
		}
	}

	err = v.Devices.Update(device)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vyos) queryDevice(ipv4 string) (domain.Utilization, error) {
	client := http.Client{}
	defer client.CloseIdleConnections()
	client.Timeout = time.Millisecond * 400

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	client.Transport = &http.Transport{
		TLSClientConfig: config,
	}

	get, err := client.Get(fmt.Sprintf("http://%s:5050/status", ipv4))
	if err != nil {
		return domain.Utilization{}, err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(get.Body)
	if err != nil {
		return domain.Utilization{}, err
	}

	util := domain.Utilization{}
	err = json.Unmarshal(buf.Bytes(), &util)
	if err != nil {
		return domain.Utilization{}, err
	}

	_ = get.Body.Close()

	return util, nil
}

func (v *Vyos) Setup() (plugin.Config, error) {
	v.pingQueue = make(chan domain.Device, 8)
	v.pingResolver = make(chan domain.Device, 8)
	v.done = make(chan bool)
	v.devicesIds = map[string]string{}
	err := v.UpdateInterval(5000)
	if err != nil {
		return plugin.Config{}, err
	}
	return v.Config, nil
}

func (v *Vyos) listen() error {
	for {
		select {
		case target := <-v.pingQueue:
			err := v.sendPing(target.Ipv4)
			if err != nil {
				return err
			}
		case result := <-v.pingResolver:
			err := v.resolvePing(result.Ipv4, result.Latency)
			if err != nil {
				return err
			}
		case <-v.done:
			return nil
		}
	}
}

func (v *Vyos) beginListening() (err error) {
	v.pingConn, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return err
	}
	defer v.pingConn.Close()

	for {
		select {
		case <-v.done:
			return nil
		default:
		}
		reply := make([]byte, 512)
		n, peer, err := v.pingConn.ReadFrom(reply)
		if err != nil {
			return err
		}
		received := time.Now()

		// Pack it up boys, we're done here
		rm, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			fmt.Println(err)
		}

		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			marshal, err := rm.Body.Marshal(1)
			if err != nil {
				return err
			}
			b := marshal[4:]
			nsec := uint64(0)
			_, err = fmt.Sscanf(string(b), "%024d", &nsec)
			if err != nil {

			}
			t := time.Unix(int64(nsec/1000000000), int64(nsec%1000000000))
			duration := received.Sub(t)
			address := peer.String()
			v.pingResolver <- domain.Device{Ipv4: address, Latency: duration}
		default:
			fmt.Println(":( Error")
		}

	}
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

			v.pingQueue <- device
		}

	}
	return nil
}

func (v *Vyos) Run() error {
	go func() {
		err := v.beginListening()
		if err != nil {
			log.Err(err)
		}
	}()
	go func() {
		err := v.listen()
		if err != nil {
			log.Err(err)
		}
	}()
	go func() {
		err := v.fetchNetworks()
		if err != nil {
			log.Err(err)
		}
	}()
	return nil
}

func (v *Vyos) Dispose() error {
	for {
		select {
		case v.done <- true:
		default:
			return nil
		}
	}
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
			v.devicesIds[device.Ipv4] = device.Id
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
