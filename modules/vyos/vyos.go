// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
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

type RemoteSocket struct {
	server  *http.Server
	router  chi.Router
	running bool

	done   chan bool
	export chan domain.Utilization
}

func NewRemoteSocket(export chan domain.Utilization) RemoteSocket {
	srv := http.Server{}
	srv.Addr = "0.0.0.0:3021"
	router := chi.NewRouter()
	srv.Handler = router
	r := RemoteSocket{
		running: false,
		server:  &srv,
		router:  router,

		done: make(chan bool),
	}
	r.export = export
	router.Get("/connect", r.connect)
	return r
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (r *RemoteSocket) connect(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Err(err)
		return
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
	}()
	res := domain.Utilization{}

	for {
		err = conn.ReadJSON(&res)
		if err != nil {
			return
		}
		select {
		case r.export <- res:
			continue
		case <-time.After(time.Millisecond * 100):
			log.Err(fmt.Errorf("utilization export timed out"))
			continue
		}

	}
}

func (r *RemoteSocket) Run() error {
	go func() {
		err := r.server.ListenAndServe()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case <-r.done:
			err := r.server.Close()
			if err != nil {
				return err
			}
			return nil
		}
	}
}

type Vyos struct {
	plugin.Module
	devicesIds          map[string]string
	pingQueue           chan domain.Device
	pingResolver        chan domain.Device
	utilizationResolver chan domain.Utilization
	networks            map[string]domain.Network
	sockets             RemoteSocket
	pingConn            *icmp.PacketConn
	done                chan bool
}

func init() {
	configVariables := []plugin.Variable{
		{
			Name:        "router",
			Default:     "10.0.1.1",
			Description: "The address to a vyOS router with api access enabled.",
		},
	}
	config := plugin.Config{
		Name:        "vyos",
		Type:        "module",
		Description: "Network interfaces controller",
		Version:     "1.0.5",
		Author:      "Braden Nicholson",
		Variables:   configVariables,
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
	if duration < 0 {
		err := v.rejectPing(ipv4)
		if err != nil {
			return err
		}
		return nil
	}

	err := v.Devices.Ping(id, duration)
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
	v.networks = map[string]domain.Network{}
	v.pingQueue = make(chan domain.Device, 8)
	v.pingResolver = make(chan domain.Device, 8)
	v.done = make(chan bool)
	v.utilizationResolver = make(chan domain.Utilization)
	v.devicesIds = map[string]string{}
	v.sockets = NewRemoteSocket(v.utilizationResolver)
	err := v.UpdateInterval(1000 * 30)
	if err != nil {
		return plugin.Config{}, err
	}
	return v.Config, nil
}

func (v *Vyos) updateUtilization(utilization domain.Utilization) error {
	deviceId := v.devicesIds[utilization.Network.Ipv4]
	if deviceId == "" {
		return nil
	}
	err := v.Devices.Utilization(deviceId, utilization)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vyos) listen() error {
	for {
		select {
		case util := <-v.utilizationResolver:
			err := v.updateUtilization(util)
			if err != nil {
				log.Err(err)
			}
		// case target := <-v.pingQueue:
		// 	err := v.sendPing(target.Ipv4)
		// 	if err != nil {
		// 		return err
		// 	}
		// case result := <-v.pingResolver:
		// 	err := v.resolvePing(result.Ipv4, result.Latency)
		// 	if err != nil {
		// 		return err
		// 	}
		case <-v.done:
			v.sockets.done <- true
			return nil
		}
	}
}

func (v *Vyos) beginListening() (err error) {
	// v.pingConn, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	// if err != nil {
	// 	return err
	// }
	// defer v.Err(v.pingConn.Close())

	// for {
	// 	select {
	// 	case <-v.done:
	// 		return nil
	// 	default:
	// 	}
	// 	// reply := make([]byte, 512)
	// 	// n, peer, err := v.pingConn.ReadFrom(reply)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	// received := time.Now()
	//
	// 	// Pack it up boys, we're done here
	// 	// rm, err := icmp.ParseMessage(1, reply[:n])
	// 	// if err != nil {
	// 	// 	fmt.Println(err)
	// 	// }
	//
	// 	// switch rm.Type {
	// 	//
	// 	// case ipv4.ICMPTypeEchoReply:
	// 	// 	marshal, err := rm.Body.Marshal(1)
	// 	// 	if err != nil {
	// 	// 		return err
	// 	// 	}
	// 	// 	b := marshal[4:]
	// 	// 	nsec := uint64(0)
	// 	// 	_, err = fmt.Sscanf(string(b), "%024d", &nsec)
	// 	// 	if err != nil {
	// 	//
	// 	// 	}
	// 	// 	t := time.Unix(int64(nsec/1000000000), int64(nsec%1000000000))
	// 	// 	duration := received.Sub(t)
	// 	// 	address := peer.String()
	// 	// 	v.pingResolver <- domain.Device{Ipv4: address, Latency: duration}
	// 	// case ipv4.ICMPTypeDestinationUnreachable:
	// 	// 	address := peer.String()
	// 	// 	v.pingResolver <- domain.Device{Ipv4: address, Latency: -1}
	// 	// default:
	// 	// 	fmt.Println(":( Error")
	// 	// }
	//
	// }
	return nil
}

func (v *Vyos) Update() error {
	if v.Ready() {
		for _, network := range v.networks {
			err := v.arpScan(network)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (v *Vyos) Run() error {
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
		err = v.sockets.Run()
		if err != nil {
			log.Err(err)
			return
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

	router, err := v.GetConfig("router")
	if err != nil {
		return nil
	}
	routerApi := fmt.Sprintf("https://%s:8005/retrieve", router)
	request, err := client.PostForm(routerApi, val)
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
		v.networks[network.Id] = network
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
