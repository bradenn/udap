// Copyright (c) 2022 Braden Nicholson

package server

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
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
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Ipv4        string `json:"ipv4"`
	Ipv6        string `json:"ipv6"`
	Hostname    string `json:"hostname"`
	Mac         string `json:"mac"`
	Go          string `json:"go"`
	Cores       int    `json:"cores"`
}

var SystemInfo System

func systemInfo() (System, error) {

	ipv4Obj, err := GetOutboundIP()
	if err != nil {
		return System{}, err
	}

	ipv4 := ipv4Obj.String()

	hostname, err := os.Hostname()
	if err != nil {
		return System{}, err
	}

	fromIpv4, err := MacFromIpv4(ipv4)
	if err != nil {
		return System{}, err
	}

	s := System{
		Name:        "UDAP",
		Version:     os.Getenv("version"),
		Environment: os.Getenv("environment"),
		Hostname:    hostname,
		Ipv4:        ipv4,
		Mac:         fromIpv4,
		Go:          runtime.Version(),
		Cores:       runtime.NumCPU(),
	}

	return s, nil
}
