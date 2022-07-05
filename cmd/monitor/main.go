// Copyright (c) 2022 Braden Nicholson

package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"net"
	"os"
	"runtime"
	"time"
)

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "10.0.1.1:80")
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {

		}
	}(conn)
	if err != nil {
		return nil, err
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func MacFromIpv4(ipv4 string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range interfaces {
		a, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range a {
			if addr.String() == fmt.Sprintf("%s/24", ipv4) {
				return i.HardwareAddr.String(), nil
			}
		}
	}
	return "", nil
}

type System struct {
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
}

func MonitorStats() (System, error) {
	system := System{}

	// Memory Statistics
	memory, err := mem.VirtualMemory()
	if err != nil {
		return system, err
	}

	system.Memory.Total = memory.Total
	system.Memory.Used = memory.UsedPercent

	// Network Statistics
	hostname, err := os.Hostname()
	if err != nil {
		return system, err
	}

	system.Network.Hostname = hostname

	ip, err := GetOutboundIP()
	if err != nil {
		return system, err
	}

	system.Network.Ipv4 = ip.String()

	mac, err := MacFromIpv4(ip.String())
	if err != nil {
		return system, err
	}

	system.Network.Mac = mac

	// Cpu statistics

	numCPU := runtime.NumCPU()

	system.Cpu.Cores = numCPU

	usages, err := cpu.Percent(time.Second, true)
	if err != nil {
		return System{}, err
	}

	system.Cpu.Usage = usages

	// Disk statistics

	usage, err := disk.Usage("/")
	if err != nil {
		return system, err
	}

	system.Disk.Total = usage.Total
	system.Disk.Used = usage.UsedPercent

	return system, nil

}

func main() {
	stats, err := MonitorStats()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(stats)
}
