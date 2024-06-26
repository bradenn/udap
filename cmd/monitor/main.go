// Copyright (c) 2022 Braden Nicholson

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type NvidiaXML struct {
	XMLName      xml.Name `xml:"nvidia_smi_log"`
	Timestamp    string   `xml:"timestamp"`
	AttachedGpus string   `xml:"attached_gpus"`
	Gpu          struct {
		Serial           string `xml:"serial"`
		Uuid             string `xml:"uuid"`
		MinorNumber      string `xml:"minor_number"`
		VbiosVersion     string `xml:"vbios_version"`
		MultigpuBoard    string `xml:"multigpu_board"`
		BoardID          string `xml:"board_id"`
		GpuPartNumber    string `xml:"gpu_part_number"`
		GpuModuleID      string `xml:"gpu_module_id"`
		FanSpeed         string `xml:"fan_speed"`
		PerformanceState string `xml:"performance_state"`
		FbMemoryUsage    struct {
			Text     string `xml:",chardata"`
			Total    string `xml:"total"`
			Reserved string `xml:"reserved"`
			Used     string `xml:"used"`
			Free     string `xml:"free"`
		} `xml:"fb_memory_usage"`
		Bar1MemoryUsage struct {
			Text  string `xml:",chardata"`
			Total string `xml:"total"`
			Used  string `xml:"used"`
			Free  string `xml:"free"`
		} `xml:"bar1_memory_usage"`
		ComputeMode string `xml:"compute_mode"`
		Utilization struct {
			Text        string `xml:",chardata"`
			GpuUtil     string `xml:"gpu_util"`
			MemoryUtil  string `xml:"memory_util"`
			EncoderUtil string `xml:"encoder_util"`
			DecoderUtil string `xml:"decoder_util"`
		} `xml:"utilization"`
		EncoderStats struct {
			Text           string `xml:",chardata"`
			SessionCount   string `xml:"session_count"`
			AverageFps     string `xml:"average_fps"`
			AverageLatency string `xml:"average_latency"`
		} `xml:"encoder_stats"`
		FbcStats struct {
			Text           string `xml:",chardata"`
			SessionCount   string `xml:"session_count"`
			AverageFps     string `xml:"average_fps"`
			AverageLatency string `xml:"average_latency"`
		} `xml:"fbc_stats"`

		Temperature struct {
			GpuTemp                string `xml:"gpu_temp"`
			GpuTempMaxThreshold    string `xml:"gpu_temp_max_threshold"`
			GpuTempSlowThreshold   string `xml:"gpu_temp_slow_threshold"`
			GpuTempMaxGpuThreshold string `xml:"gpu_temp_max_gpu_threshold"`
			GpuTargetTemperature   string `xml:"gpu_target_temperature"`
			MemoryTemp             string `xml:"memory_temp"`
			GpuTempMaxMemThreshold string `xml:"gpu_temp_max_mem_threshold"`
		} `xml:"temperature"`

		PowerReadings struct {
			PowerState         string `xml:"power_state"`
			PowerManagement    string `xml:"power_management"`
			PowerDraw          string `xml:"power_draw"`
			PowerLimit         string `xml:"power_limit"`
			DefaultPowerLimit  string `xml:"default_power_limit"`
			EnforcedPowerLimit string `xml:"enforced_power_limit"`
			MinPowerLimit      string `xml:"min_power_limit"`
			MaxPowerLimit      string `xml:"max_power_limit"`
		} `xml:"power_readings"`
		Clocks struct {
			Text          string `xml:",chardata"`
			GraphicsClock string `xml:"graphics_clock"`
			SmClock       string `xml:"sm_clock"`
			MemClock      string `xml:"mem_clock"`
			VideoClock    string `xml:"video_clock"`
		} `xml:"clocks"`
		ApplicationsClocks struct {
			Text          string `xml:",chardata"`
			GraphicsClock string `xml:"graphics_clock"`
			MemClock      string `xml:"mem_clock"`
		} `xml:"applications_clocks"`
		DefaultApplicationsClocks struct {
			Text          string `xml:",chardata"`
			GraphicsClock string `xml:"graphics_clock"`
			MemClock      string `xml:"mem_clock"`
		} `xml:"default_applications_clocks"`
		MaxClocks struct {
			Text          string `xml:",chardata"`
			GraphicsClock string `xml:"graphics_clock"`
			SmClock       string `xml:"sm_clock"`
			MemClock      string `xml:"mem_clock"`
			VideoClock    string `xml:"video_clock"`
		} `xml:"max_clocks"`
		Processes []struct {
			Text        string `xml:",chardata"`
			ProcessInfo struct {
				Text              string `xml:",chardata"`
				GpuInstanceID     string `xml:"gpu_instance_id"`
				ComputeInstanceID string `xml:"compute_instance_id"`
				Pid               string `xml:"pid"`
				Type              string `xml:"type"`
				ProcessName       string `xml:"process_name"`
				UsedMemory        string `xml:"used_memory"`
			} `xml:"process_info"`
		} `xml:"processes"`
		AccountedProcesses string `xml:"accounted_processes"`
	} `xml:"gpu"`
}

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
	Compute []domain.Compute `json:"compute"`
}

func queryGPU() (NvidiaXML, error) {
	cmd := exec.Command("nvidia-smi", "-q", "-x")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return NvidiaXML{}, err
	}

	nvm := NvidiaXML{}
	err = xml.Unmarshal(output, &nvm)
	if err != nil {
		return NvidiaXML{}, err
	}

	return nvm, nil
}

func parseCompute(xml NvidiaXML) ([]domain.Compute, error) {

	x := xml.Gpu

	c := domain.Compute{
		FanSpeed: 0,
		Temperature: struct {
			Current  int `json:"current"`
			Throttle int `json:"throttle"`
			Target   int `json:"target"`
			Max      int `json:"max"`
		}{
			Current:  0,
			Throttle: 0,
			Target:   0,
			Max:      0,
		},
		Utilization: struct {
			GPU    int `json:"gpu"`
			Memory int `json:"memory"`
		}{
			GPU:    0,
			Memory: 0,
		},
		Power: struct {
			Draw float64 `json:"draw"`
			Max  float64 `json:"max"`
		}{
			Draw: 0,
			Max:  0,
		},
		Memory: struct {
			Used     int `json:"used"`
			Reserved int `json:"reserved"`
			Total    int `json:"total"`
		}{},
		Clocks: struct {
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
		}{},
		Processes: []domain.ComputeProcess{},
	}

	_, err := fmt.Sscanf(x.FanSpeed, "%d %%", &c.FanSpeed)
	if err != nil {
		return nil, err
	}

	// GPU Power
	_, err = fmt.Sscanf(x.PowerReadings.PowerLimit, "%f W", &c.Power.Max)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Sscanf(x.PowerReadings.PowerDraw, "%f W", &c.Power.Draw)
	if err != nil {
		return nil, err
	}

	// GPU Temperature
	_, err = fmt.Sscanf(x.Temperature.GpuTemp, "%d C", &c.Temperature.Current)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Sscanf(x.Temperature.GpuTargetTemperature, "%d C", &c.Temperature.Target)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Sscanf(x.Temperature.GpuTempSlowThreshold, "%d C", &c.Temperature.Throttle)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Sscanf(x.Temperature.GpuTempMaxThreshold, "%d C", &c.Temperature.Max)
	if err != nil {
		return nil, err
	}

	// Utilization

	_, err = fmt.Sscanf(x.Utilization.GpuUtil, "%d %%", &c.Utilization.GPU)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Sscanf(x.Utilization.MemoryUtil, "%d %%", &c.Utilization.Memory)
	if err != nil {
		return nil, err
	}

	// Memory Usage

	_, err = fmt.Sscanf(x.FbMemoryUsage.Free, "%d MiB", &c.Memory.Used)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.FbMemoryUsage.Total, "%d MiB", &c.Memory.Total)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.FbMemoryUsage.Reserved, "%d MiB", &c.Memory.Reserved)
	if err != nil {
		return nil, err
	}

	// Clocks

	_, err = fmt.Sscanf(x.Clocks.GraphicsClock, "%d MHz", &c.Clocks.Graphics.Current)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.MaxClocks.GraphicsClock, "%d MHz", &c.Clocks.Graphics.Max)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.Clocks.SmClock, "%d MHz", &c.Clocks.Streaming.Current)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.MaxClocks.SmClock, "%d MHz", &c.Clocks.Streaming.Max)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.Clocks.MemClock, "%d MHz", &c.Clocks.Memory.Current)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.MaxClocks.MemClock, "%d MHz", &c.Clocks.Memory.Max)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.Clocks.VideoClock, "%d MHz", &c.Clocks.Video.Current)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Sscanf(x.MaxClocks.VideoClock, "%d MHz", &c.Clocks.Video.Max)
	if err != nil {
		return nil, err
	}

	if len(x.Processes) > 0 {
		for _, process := range x.Processes {
			proc := domain.ComputeProcess{}
			proc.Name = process.ProcessInfo.ProcessName
			proc.PID = process.ProcessInfo.Pid
			_, err = fmt.Sscanf(process.ProcessInfo.UsedMemory, "%d MiB", &proc.Memory)
			if err != nil {
				return nil, err
			}
			c.Processes = append(c.Processes, proc)
		}
	}

	return []domain.Compute{c}, nil

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

	usages, err := cpu.Percent(0, true)
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

	// Compute

	if isGPU {
		gpu, err := queryGPU()
		if err != nil {
			return system, err
		}
		compute, err := parseCompute(gpu)
		if err != nil {
			return system, err
		}

		system.Compute = compute
	}

	return system, nil

}

func Status(writer http.ResponseWriter, request *http.Request) {
	stats, err := MonitorStats()
	if err != nil {
		return
	}

	marshal, err := json.Marshal(stats)
	if err != nil {
		return
	}

	_, err = writer.Write(marshal)
	if err != nil {
		return
	}
}

var isGPU bool

func SendUpdate(conn *websocket.Conn) error {
	stats, err := MonitorStats()
	if err != nil {
		log.Err(err)
	}
	err = conn.WriteJSON(stats)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	isGPU = false

	if len(os.Args) > 1 {
		if os.Args[1] == "gpu" {
			fmt.Println("Using GPU mode")
			isGPU = true
		}
	}

	for {
		time.Sleep(time.Second * 5)
		remote := url.URL{Scheme: "ws", Host: "10.0.1.2" + ":" + "3021", Path: "/connect"}
		conn, _, err := websocket.DefaultDialer.Dial(remote.String(), nil)
		if err != nil {
			log.Event("Connection failed: Trying again in 5 seconds")
			continue
		}
		log.Event("Connected")
		go func() {
			_, _, err = conn.ReadMessage()
			if err != nil {
				return
			}
		}()

		for {
			err = SendUpdate(conn)
			if err != nil {
				break
			}
			time.Sleep(time.Second * 4)
		}

		log.Event("Connection Lost: Attempting to reconnect in 5 seconds")
	}

	// server := http.Server{}
	//
	// http.HandleFunc("/status", Status)
	//
	// server.Addr = ":5050"
	// fmt.Println("Running on port :5050")
	// err = server.ListenAndServe()
	// if err != nil {
	// 	return
	// }
	// err := server.ListenAndServeTLS("./certs/monitor.crt", "./certs/monitor.key")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
