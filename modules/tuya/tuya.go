// Copyright (c) 2023 Braden Nicholson

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
)

var Module Tuya

const tuyaKey = "yGAdlopoPVldABfn"

func DecodeUDP(data []byte) ([]byte, error) {
	res := md5.Sum([]byte(tuyaKey))

	data = data[20 : len(data)-8]
	decrypted := make([]byte, len(data))

	c, _ := aes.NewCipher([]byte(fmt.Sprintf("%s", res)))

	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		c.Decrypt(decrypted[bs:be], data[bs:be])
	}

	// remove the padding. The last character in the byte array is the number of padding chars
	paddingSize := int(decrypted[len(decrypted)-1])
	return decrypted[0 : len(decrypted)-paddingSize], nil

}

func unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("Cannot unpad an empty slice")
	}
	paddingSize := int(data[len(data)-1])
	if paddingSize > aes.BlockSize || paddingSize > len(data) {
		return nil, fmt.Errorf("Invalid padding size")
	}
	return data[:len(data)-paddingSize], nil
}

func padF(data []byte) []byte {
	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func EncryptPacket(key []byte, data []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	size := cipher.BlockSize()
	data = padF(data)
	out := make([]byte, len(data))
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Encrypt(out[bs:be], data[bs:be])
	}
	return out, nil
}

func DecodePacket(key []byte, data []byte) ([]byte, error) {

	if len(key) < 16 {
		return nil, fmt.Errorf("Key is too short (minimum length is 16 bytes)")
	}

	if len(data)%16 != 0 {
		return nil, fmt.Errorf("Data length is not a multiple of the block size (16 bytes)")
	}

	decrypted := make([]byte, len(data))
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	size := 16
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}
	// remove the padding. The last character in the byte array is the number of padding chars
	return unpad(decrypted)
}

var HEADER = 0x000055AA

var PROTOCOL_3x_HEADER = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type controlCmd struct {
	DevId string          `json:"devId"`
	Uid   string          `json:"uid"`
	T     string          `json:"t"`
	Dps   json.RawMessage `json:"dps,omitempty"`
}

type statusCmd struct {
	GwId  string `json:"gwId"`
	DevId string `json:"devId"`
	Uid   string `json:"uid"`
	T     string `json:"t"`
}

type Light struct {
	Ip          string `json:"ip"`
	GwId        string `json:"gwId"`
	DevId       string `json:"devId"`
	LocalId     string `json:"localId"`
	Active      int    `json:"active"`
	Mode        string
	Hue         int
	Dim         int
	UUID        string `json:"uid"`
	Encrypt     bool   `json:"encrypt"`
	ProductKey  string `json:"productKey"`
	Version     string `json:"version"`
	Seq         int
	Latest      DpsPayload `json:"latest"`
	Reconnects  int
	lastStatus  time.Time
	lastRequest time.Time
	mutex       sync.Mutex
	socket      net.Conn
}

const (
	ModeColorTemp = "white"
	ModeColor     = "colour"
)

const (
	Control    = 7
	Status     = 0x0a
	HeartBeat  = 9
	IsOn       = 20
	Mode       = 21
	Brightness = 22
	ColorTemp  = 23
	Color      = 24
	RealTime   = 28
)

func Pad(data []byte) []byte {
	p := aes.BlockSize - len(data)%aes.BlockSize
	return []byte(strings.Repeat("0", len(data)+p))
}

func (l *Light) sendModeCommandDps(mode string, key int, value any) error {
	dps := make(map[string]any)
	dps[fmt.Sprintf("%d", Mode)] = mode
	dps[fmt.Sprintf("%d", key)] = value

	l.Mode = mode

	marshal, err := json.Marshal(dps)
	if err != nil {
		return err
	}

	payload, err := l.createPayload(marshal)
	if err != nil {
		return err
	}

	err = l.sendCommandPayload(payload)
	if err != nil {
		return err
	}

	l.lastRequest = time.Now()
	//countdown := time.NewTimer(time.Second)
	//
	//select {
	//case l.requestQueue <- payload:
	//	countdown.Stop()
	//	break
	//case <-countdown.C:
	//	break
	//}

	return nil
}

func (l *Light) sendCommandDps(key int, value any) error {
	dps := make(map[string]any)
	dps[fmt.Sprintf("%d", Mode)] = l.Mode
	dps[fmt.Sprintf("%d", key)] = value

	marshal, err := json.Marshal(dps)
	if err != nil {
		return err
	}

	payload, err := l.createPayload(marshal)
	if err != nil {
		return err
	}

	err = l.sendCommandPayload(payload)
	if err != nil {
		return err
	}

	return nil
}

func (l *Light) createPayload(dps json.RawMessage) ([]byte, error) {
	c := controlCmd{
		DevId: l.GwId,
		Uid:   l.GwId,
		T:     fmt.Sprintf("%d", time.Now().Unix()),
		Dps:   dps,
	}

	payload, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func convertRange(value, fromMin, fromMax, toMin, toMax float64) float64 {
	fromRange := fromMax - fromMin
	toRange := toMax - toMin
	return (((value - fromMin) * toRange) / fromRange) + toMin
}

// RGBToHSV converts an RGB value (0-255), to a HSV value (0-1)
func RGBToHSV(r, g, b int) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255
	max := math.Max(math.Max(rf, gf), bf)
	min := math.Min(math.Min(rf, gf), bf)
	h := 0.0
	s := 0.0
	v := max
	d := max - min
	if max != 0 {
		s = d / max
	}
	if max == min {
		h = 0
	} else {
		if max == rf {
			h = (gf - bf) / d
			if gf < bf {
				h += 6
			}
		} else if max == gf {
			h = (bf-rf)/d + 2
		} else if max == bf {
			h = (rf-gf)/d + 4
		}
		h /= 6
	}
	return h, s, v
}

func (l *Light) setOn(on bool) error {
	return l.sendCommandDps(IsOn, on)
}
func (l *Light) setDimFine(dim int) error {

	if dim > 1000 || dim < 0 {
		return fmt.Errorf("dim.fine value must be a value betweeen 0 and 1000 inclusive")
	}
	if l.Mode == "colour" {
		l.Dim = dim
		return l.setHueDim(l.Hue, dim)
	}
	return l.sendCommandDps(Brightness, dim)
}

func (l *Light) setDim(dim int) error {

	if dim > 100 || dim < 0 {
		return fmt.Errorf("dim value must be a value betweeen 0 and 100 inclusive")
	}
	if l.Mode == "colour" {
		l.Dim = dim * 10
		return l.setHueDim(l.Hue, dim*10)
	}
	return l.sendCommandDps(Brightness, dim*10)
}

// setCCT sets the bulb's color corrected value on a scale from 2700 to 6500.
func (l *Light) setCCT(cct int) error {
	// The color temperature is set from 0-1000, where 0 is the lowest the bulb can handler, and 1000 is the highest color temp,
	// for example, 500 would represent a cct value of 4600
	target := convertRange(math.Min(math.Max(float64(cct), 2700), 6500), 2700, 6500, 0, 1000)
	target = math.Max(0, target)
	target = math.Min(1000, target)
	return l.sendModeCommandDps(ModeColorTemp, ColorTemp, target)
}

func (l *Light) setRGB(r, g, b int) error {
	// The light expects an HSV value in the mode Hue: 0-360, Saturation: 0-1000, and Value: 0-1000
	h, s, v := RGBToHSV(r, g, b)
	values := []int{int(h * 360), int(s * 1000), int(v * 1000)}
	hexVal := ""
	for _, value := range values {
		local := fmt.Sprintf("%04x", value)
		hexVal += local
	}

	return l.sendModeCommandDps(ModeColor, Color, hexVal)
}

func (l *Light) setRGBW(r, g, b, w int) error {
	// The light expects an HSV value in the mode Hue: 0-360, Saturation: 0-1000, and Value: 0-1000
	h, s, v := RGBToHSV(r, g, b)
	values := []int{int(h * 360), int(s * 1000), int(v * math.Min(float64(w*100), 1000))}
	hexVal := ""
	for _, value := range values {
		local := fmt.Sprintf("%04x", value)
		hexVal += local
	}

	return l.sendModeCommandDps(ModeColor, Color, hexVal)
}

func (l *Light) setHueDim(r int, dim int) error {
	// The light expects an HSV value in the mode Hue: 0-360, Saturation: 0-1000, and Value: 0-1000

	values := []int{int(math.Min(float64(r), 360)), int(1000), int(math.Min(float64(dim), 1000))}
	hexVal := ""
	for _, value := range values {
		local := fmt.Sprintf("%04x", value)
		hexVal += local
	}

	l.Hue = r

	return l.sendModeCommandDps(ModeColor, Color, hexVal)
}

func (l *Light) setHue(r int) error {
	// The light expects an HSV value in the mode Hue: 0-360, Saturation: 0-1000, and Value: 0-1000

	values := []int{int(math.Min(float64(r), 360)), int(1000), int(l.Dim)}
	hexVal := ""
	for _, value := range values {
		local := fmt.Sprintf("%04x", value)
		hexVal += local
	}
	l.Hue = r
	return l.sendModeCommandDps(ModeColor, Color, hexVal)
}

func (l *Light) poll() error {
	err := l.sendStatusPayload()
	if err != nil {
		return err
	}
	return nil
}

func (l *Light) decodePayload(data []byte, resp bool) ([]byte, error) {
	// Create a reader to process the packer from
	buf := bytes.NewReader(data)

	// Read the prefix from the front of the packet
	var header uint32
	_ = binary.Read(buf, binary.BigEndian, &header)
	// Read the sequence number from the packet header
	var sequence uint32
	_ = binary.Read(buf, binary.BigEndian, &sequence)
	// Read the control code from the packet header
	var control uint32
	_ = binary.Read(buf, binary.BigEndian, &control)
	// Read the packet data length from the header
	var length uint32
	_ = binary.Read(buf, binary.BigEndian, &length)
	// Read the control code from the packet header
	var ret uint32
	_ = binary.Read(buf, binary.BigEndian, &ret)
	// Read the remaining packet
	payload := make([]byte, length)
	_, err := buf.Read(payload)
	if err != nil {
		return nil, err
	}
	packet := append([]byte("3.3"), PROTOCOL_3x_HEADER...)
	// Compute the size of the version and protocol spacer from the packet
	if resp && len(payload) > len(packet) {
		// Remove the packet by slicing the bytes
		payload = payload[len(packet):]
	}

	if len(payload) <= 12 {
		return nil, fmt.Errorf("invalid packet")
	}

	// Cut off the footer, the suffix and crc
	payload = payload[:len(payload)-12]
	// Decrypt the packet with the LocalId
	i, err := DecodePacket([]byte(l.LocalId), payload)
	if err != nil {
		return nil, err
	}

	return i, nil
}

type DpsPayload struct {
	Dps struct {
		On    bool   `json:"20,omitempty"`
		Mode  string `json:"21,omitempty"`
		Dim   int    `json:"22,omitempty"`
		Cct   int    `json:"23,omitempty"`
		Color string `json:"24,omitempty"`
		T25   string `json:"25,omitempty"`
		T26   int    `json:"26,omitempty"`
	} `json:"dps"`
}

//{"dps":{"20":false,"21":"white","22":1000,"23":0,"24":"000003e803e8","25":"000e0d0000000000000000c80000","26":0}}
func (l *Light) parsePayload(payload []byte) error {
	data := DpsPayload{}
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	l.Mode = data.Dps.Mode
	l.Latest = data
	l.lastStatus = time.Now()

	return nil
}

func (l *Light) sendStatusPayload() error {

	// Generate the status command structure
	c := statusCmd{
		GwId:  l.GwId,
		DevId: l.GwId,
		Uid:   l.GwId,
		T:     fmt.Sprintf("%d", time.Now().Unix()),
	}

	// Convert the structure into json
	payload, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Encrypt the payload with the LocalId
	payload, err = EncryptPacket([]byte(l.LocalId), payload)
	if err != nil {
		return err
	}
	// Prefix the version and spacing to the encrypted payload
	//packet := append([]byte("3.3"), PROTOCOL_3x_HEADER...)
	//packet = append(packet, payload...)
	packet := payload
	// Define a buffer to write the metadata to
	var buf bytes.Buffer
	// Write the header prefix
	_ = binary.Write(&buf, binary.BigEndian, uint32(0x000055AA))
	// Write the sequence number
	_ = binary.Write(&buf, binary.BigEndian, int32(l.Seq))
	// Write the Status code
	_ = binary.Write(&buf, binary.BigEndian, int32(Status))
	// Write the packet length
	_ = binary.Write(&buf, binary.BigEndian, int32(len(packet)+8))
	// Write the metadata buffer to the front of the packet
	packet = append(buf.Bytes(), packet...)
	// Generate the crc for the packet
	crc := crc32.ChecksumIEEE(packet) & 0xFFFFFFFF
	// Reset the buffer to use again
	buf = bytes.Buffer{}
	// Write the crc to the end of the buffer
	_ = binary.Write(&buf, binary.BigEndian, crc)
	// Write the suffix to the end of the buffer
	_ = binary.Write(&buf, binary.BigEndian, uint32(0x0000AA55))
	// Add the new buffer to the end of the packet
	packet = append(packet, buf.Bytes()...)

	//s := time.Now()
	err = l.WritePacket(packet, true)
	//log.Log("%s = %s", l.Ip, time.Since(s))

	return err

}

func (l *Light) sendCommandPayload(data []byte) error {

	// Encrypt the payload
	payload, err := EncryptPacket([]byte(l.LocalId), data)
	if err != nil {
		return err
	}
	//fmt.Printf("%s\n", data)
	// Prefix the version and spacing to the encrypted payload
	packet := append([]byte("3.3"), PROTOCOL_3x_HEADER...)
	packet = append(packet, payload...)
	// Define a buffer to write the metadata to
	var buf bytes.Buffer
	// Write the header prefix
	_ = binary.Write(&buf, binary.BigEndian, uint32(0x000055AA))
	// Write the sequence number
	_ = binary.Write(&buf, binary.BigEndian, int32(l.Seq))
	// Write the control code
	_ = binary.Write(&buf, binary.BigEndian, int32(Control))
	// Write the packet length
	_ = binary.Write(&buf, binary.BigEndian, int32(len(packet)+8))
	// Write the metadata buffer to the front of the packet
	packet = append(buf.Bytes(), packet...)
	// Generate the crc for the packet
	crc := crc32.ChecksumIEEE(packet) & 0xFFFFFFFF
	// Reset the buffer to use again
	buf = bytes.Buffer{}
	// Write the crc to the end of the buffer
	_ = binary.Write(&buf, binary.BigEndian, crc)
	// Write the suffix to the end of the buffer
	_ = binary.Write(&buf, binary.BigEndian, uint32(0x0000AA55))
	// Add the new buffer to the end of the packet
	packet = append(packet, buf.Bytes()...)

	return l.WritePacket(packet, false)

}

func (l *Light) ConnectSocket(read bool) (net.Conn, error) {
	if l.socket != nil {
		err := l.socket.SetWriteDeadline(time.Now().Add(time.Millisecond * 200))
		if err != nil {
			return nil, err
		}
		_, err = l.socket.Write([]byte{})
		if err == nil {
			return l.socket, nil
		}
		l.socket = nil
	}

	attempts := 0
	for attempts < 3 {
		socket, err := net.Dial("tcp", fmt.Sprintf("%s:%d", l.Ip, 6668))
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			attempts++
			continue
		}
		l.socket = socket
		return socket, err
	}

	return nil, fmt.Errorf("reconnect failed")
}

func (l *Light) WritePacket(packet []byte, read bool) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	var err error
	socket, err := l.ConnectSocket(read)
	//if !read {
	//	defer socket.Close()
	//}

	if err != nil {
		return err
	}

	//if l.socket == nil {
	//	l.socket, err = l.ConnectSocket()
	//	if err != nil {
	//		l.socket = nil
	//		return err
	//	}
	//}

	err = socket.SetDeadline(time.Now().Add(500 * time.Millisecond))
	if err != nil {
		return err
	}

	expected := len(packet)
	n := 0
	for n < expected {
		var sent int
		sent, err = socket.Write(packet[n:])
		if err != nil {
			return err
		}
		if sent == 0 {
			break
		}
		n += sent
	}

	// Craft a buffer to contain the incoming data
	response := make([]byte, 512)
	if socket == nil {
		return nil
	}
	// Wait for and read the response from the device
	n, err = socket.Read(response)
	if err != nil {
		return err
	}
	// Cut off the unused bits
	response = response[:n]

	// Decode/Decrypt the data
	processed, err := l.decodePayload(response, !read)
	if err != nil {
		return nil
	}

	if read {
		err = l.parsePayload(processed)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Tuya) updateLight(light *Light) error {

	val, ok := t.Lights[light.GwId]
	if !ok || val == nil {
		light.mutex = sync.Mutex{}
		light.Reconnects = 0

		err := t.registerLight(light)
		if err != nil {
			return err
		}

		t.Lights[light.GwId] = light
		val = t.Lights[light.GwId]
		val.lastStatus = time.Now()
	} else {
		val.lastStatus = time.Now()
	}

	return nil
}

func (t *Tuya) lightMux() error {
	for {
		select {
		case light := <-t.ping:
			err := t.updateLight(light)
			if err != nil {
				log.Err(err)
			}
		}
	}
}

func (t *Tuya) Scan() error {

	pc, err := net.ListenPacket("udp", ":6667")
	if err != nil {
		return err
	}

	defer pc.Close()

	go func() {
		err = t.lightMux()
		if err != nil {
			return
		}
	}()

	buf := make([]byte, 256)
	go func() {
		<-t.done
		err = pc.Close()
		if err != nil {
			return
		}
		close(t.ping)
	}()
	/*	err = pc.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			return err
		}
	*/
	for {
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			return nil
		}

		i, err := DecodeUDP(buf[:n])
		if err != nil {
			return err
		}

		l := Light{}
		err = json.Unmarshal(i, &l)
		if err != nil {
			return err
		}

		if l.GwId == "" {
			fmt.Println("Um incorrect light id...")
			continue
		}

		t.ping <- &l

	}

}

type Tuya struct {
	plugin.Module
	entityId string

	entities map[string]string
	receiver chan domain.Attribute
	ping     chan *Light
	interval int
	ready    bool
	done     chan bool
	Lights   map[string]*Light
}

func init() {
	Module = Tuya{
		entities: map[string]string{},
		ping:     make(chan *Light, 8),
		done:     make(chan bool),
		Lights:   map[string]*Light{},
		interval: 0,
	}
	Module.Config = plugin.Config{
		Name:        "tuya",
		Type:        "daemon",
		Description: "Tuya Device interface",
		Version:     "0.1.2",
		Author:      "Braden Nicholson",
	}
}

func (t *Tuya) Setup() (plugin.Config, error) {
	err := t.UpdateInterval(1000 * 5)
	if err != nil {
		return plugin.Config{}, err
	}
	return t.Config, nil
}

func (t *Tuya) updateOne(light *Light) error {

	stamp := time.Now()
	err := light.sendStatusPayload()
	if err != nil {
		return err
	}

	eid, _ := t.findEntityIdByLightId(light.GwId)
	h := int(0)
	s := int(0)
	v := int(0)
	_, err = fmt.Sscanf(light.Latest.Dps.Color, "%04x%04x%04x", &h, &s, &v)
	if err != nil {
		//log.Event("Failed to parse color: '%s'", light.Latest.Dps.Color)
		return err
	}

	//if light.lastStatus.Before(light.lastRequest) {
	//	return err
	//}

	_ = t.Attributes.Update(eid, "hue", fmt.Sprintf("%d", h), stamp)
	light.Hue = h
	light.Dim = v
	if err != nil {
		return err
	}

	if light.Latest.Dps.On {
		err = t.Attributes.Update(eid, "on", "true", stamp)
		if err != nil {
			return err
		}
	} else {
		err = t.Attributes.Update(eid, "on", "false", stamp)
		if err != nil {
			return err
		}
	}

	if light.Mode == "colour" {
		err = t.Attributes.Update(eid, "dim", fmt.Sprintf("%d", light.Dim/10), stamp)
		if err != nil {
			return err
		}
	} else {
		err = t.Attributes.Update(eid, "dim", fmt.Sprintf("%d", light.Latest.Dps.Dim/10), stamp)
		if err != nil {
			return err
		}
	}

	err = t.Attributes.Update(eid, "cct", fmt.Sprintf("%d", int(convertRange(float64(light.Latest.Dps.Cct), 0, 1000, 2700, 6500))), stamp)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tuya) updateAll() error {
	wg := sync.WaitGroup{}
	wg.Add(len(t.Lights))
	for _, l := range t.Lights {
		go func(light *Light) {
			defer wg.Done()
			err := t.updateOne(light)
			if err != nil {
				return
			}
		}(l)
	}

	wg.Wait()
	//t.LastUpdate = time.Now()
	return nil
}

func (t *Tuya) Update() error {
	if !t.ready {
		return nil
	}

	if t.Ready() {
		err := t.updateAll()
		if err != nil {
			return err
		}
	}
	return nil
}

func mapRange(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func (t *Tuya) deriveAmbientColorTemperature() int {
	delta := math.Abs(float64(time.Now().Hour() - 11))

	return int(mapRange(delta, 0, 12, 6000, 4100))
}

func (t *Tuya) findEntityIdByLightId(lightId string) (string, error) {
	for id, s := range t.entities {
		if s == lightId {
			return id, nil
		}
	}
	return "", fmt.Errorf("invalid")
}

func (t *Tuya) handleRequest(attribute domain.Attribute) error {
	var err error

	lightId := t.entities[attribute.Entity]
	light := t.Lights[lightId]

	switch attribute.Key {
	case "on":
		err = light.setOn(attribute.Request == "true")
		if err != nil {
			return err
		}
		break
	case "dim":
		atoi, err2 := strconv.Atoi(attribute.Request)
		if err2 != nil {
			return err2
		}
		err = light.setDim(atoi)
		if err != nil {
			return err
		}
		break
	case "dim.fine":
		atoi, err2 := strconv.Atoi(attribute.Request)
		if err2 != nil {
			return err2
		}
		err = light.setDimFine(atoi)
		if err != nil {
			return err
		}
		break
	case "cct":
		atoi, err2 := strconv.Atoi(attribute.Request)
		if err2 != nil {
			return err2
		}
		err = light.setCCT(atoi)
		if err != nil {
			return err
		}
		break
	case "hue":
		atoi, err2 := strconv.Atoi(attribute.Request)
		if err2 != nil {
			return err2
		}
		err = light.setHue(atoi)
		if err != nil {
			return err
		}
		break
	case "spectral":
		power := "on"
		mode := "cct"
		value := 5500
		dim := 0

		chunks := strings.Split(attribute.Request, ";")
		if len(chunks) == 3 {
			power = chunks[0]
			params := strings.Split(chunks[1], ":")
			mode = params[0]
			v, err := strconv.ParseInt(params[1], 10, 32)
			if err != nil {
				return err
			}
			value = int(v)

			d, err := strconv.ParseInt(chunks[2], 10, 32)
			if err != nil {
				return err
			}

			dim = int(d)
		} else {
			break
		}

		if power == "off" {
			err = light.setOn(false)
			if err != nil {
				return err
			}
		}

		if mode == "cct" {
			err = light.setDim(dim)
			if err != nil {
				return err
			}

			if value == -1 {
				value = t.deriveAmbientColorTemperature()
			}

			err = light.setCCT(value)
			if err != nil {
				return err
			}

		} else if mode == "hue" {
			err = light.setHueDim(value, dim)
			if err != nil {
				return err
			}
		}

		if power == "on" {
			err = light.setOn(true)
			if err != nil {
				return err
			}
		}
		break
	default:
		return fmt.Errorf("invalid attribute key")
	}

	err = t.Attributes.Set(attribute.Entity, attribute.Key, attribute.Request)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tuya) mux(index int) error {
	for {
		if !t.ready {
			time.Sleep(time.Second)
			continue
		}
		select {
		case attr := <-t.receiver:
			go func() {
				err := t.handleRequest(attr)
				if err != nil {
					log.Err(err)
				}
			}()
			break
		}
	}
}

func (t *Tuya) registerLight(light *Light) error {
	entity := &domain.Entity{
		Name:   light.GwId,
		Type:   "spectrum",
		Module: t.Config.Name,
	}

	data, err := base64.StdEncoding.DecodeString(os.Getenv(light.GwId))
	if err != nil {
		return err
	}

	light.LocalId = string(data)

	err = t.Entities.Register(entity)
	if err != nil {
		return err
	}

	t.entities[entity.Id] = light.GwId

	err = t.Attributes.Register(&domain.Attribute{
		Key:     "on",
		Value:   "false",
		Request: "false",
		Type:    "toggle",
		Order:   0,
		Entity:  entity.Id,
		Channel: t.receiver,
	})

	err = t.Attributes.Register(&domain.Attribute{
		Key:     "dim",
		Value:   "0",
		Request: "0",
		Type:    "range",
		Order:   1,
		Entity:  entity.Id,
		Channel: t.receiver,
	})

	err = t.Attributes.Register(&domain.Attribute{
		Key:     "dim.fine",
		Value:   "0",
		Request: "0",
		Type:    "range",
		Order:   1,
		Entity:  entity.Id,
		Channel: t.receiver,
	})

	err = t.Attributes.Register(&domain.Attribute{
		Key:     "cct",
		Value:   "2000",
		Request: "2000",
		Type:    "range",
		Order:   2,
		Entity:  entity.Id,
		Channel: t.receiver,
	})

	err = t.Attributes.Register(&domain.Attribute{
		Key:     "hue",
		Value:   "0",
		Request: "0",
		Type:    "range",
		Order:   3,
		Entity:  entity.Id,
		Channel: t.receiver,
	})

	err = t.Attributes.Register(&domain.Attribute{
		Key:     "spectral",
		Value:   "off;cct:5500;100",
		Request: "off;cct:5500;100",
		Type:    "media",
		Order:   4,
		Entity:  entity.Id,
		Channel: t.receiver,
	})

	return nil
}

func (t *Tuya) runScan() error {

	go func() {
		err := t.Scan()
		if err != nil {
			log.Err(err)
		}
	}()

	t.ready = true
	return nil
}

func (t *Tuya) Run() error {
	t.ready = false

	t.receiver = make(chan domain.Attribute, 8)

	go func() {
		for {
			err := t.mux(0)
			if err != nil {
				fmt.Printf("The tuya mux had to close unexpectedly... %s\n", err.Error())
				time.Sleep(time.Second)
				continue
			}
		}
	}()

	go func() {
		err := t.runScan()
		if err != nil {
			return
		}
	}()

	return nil
}
