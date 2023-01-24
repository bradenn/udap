// Copyright (c) 2023 Braden Nicholson

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
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

type Light struct {
	Ip         string `json:"ip"`
	GwId       string `json:"gwId"`
	DevId      string `json:"devId"`
	LocalId    string `json:"localId"`
	Active     int    `json:"active"`
	UUID       string `json:"uid"`
	Encrypt    bool   `json:"encrypt"`
	ProductKey string `json:"productKey"`
	Version    string `json:"version"`
	Seq        int
	Latest     DpsPayload `json:"latest"`
	mutex      sync.Mutex
	socket     net.Conn
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
)

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

func Pad(data []byte) []byte {
	p := aes.BlockSize - len(data)%aes.BlockSize
	return []byte(strings.Repeat("0", len(data)+p))
}

func (l *Light) sendModeCommandDps(mode string, key int, value any) error {
	dps := make(map[string]any)
	dps[fmt.Sprintf("%d", Mode)] = mode
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

func (l *Light) sendCommandDps(key int, value any) error {
	dps := make(map[string]any)
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

func (l *Light) setDim(dim int) error {

	if dim > 100 || dim < 0 {
		return fmt.Errorf("dim value must be a value betweeen 0 and 100 inclusive")
	}

	return l.sendCommandDps(Brightness, dim*10)
}

// setCCT sets the bulb's color corrected value on a scale from 2700 to 6500.
func (l *Light) setCCT(cct int) error {
	// The color temperature is set from 0-1000, where 0 is the lowest the bulb can handler, and 1000 is the highest color temp,
	// for example, 500 would represent a cct value of 4600
	target := convertRange(math.Min(math.Max(float64(cct), 2700), 6500), 2700, 6500, 0, 1000)
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
func (l *Light) setHue(r int) error {
	// The light expects an HSV value in the mode Hue: 0-360, Saturation: 0-1000, and Value: 0-1000

	values := []int{int(math.Min(float64(r), 360)), int(1000), int(1000)}
	hexVal := ""
	for _, value := range values {
		local := fmt.Sprintf("%04x", value)
		hexVal += local
	}

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
	// Compute the size of the version and protocol spacer from the packet
	if resp {
		packet := append([]byte("3.3"), PROTOCOL_3x_HEADER...)
		// Remove the packet by slicing the bytes
		payload = payload[len(packet):]
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

func (l *Light) parsePayload(payload []byte) error {
	data := DpsPayload{}
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	l.Latest = data

	return nil
}

func (l *Light) sendStatusPayload() error {
	// Get the socket if it is still open, or open a new one
	socket, err := l.getSocket()
	if err != nil {
		return err
	}
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
	// Write the packet to the socket
	_, err = socket.Write(packet)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Craft a buffer to contain the incoming data
	response := make([]byte, 512)
	err = socket.SetDeadline(time.Now().Add(time.Millisecond * 250))
	if err != nil {
		return err
	}
	// Wait for and read the response from the device
	n, err := socket.Read(response)
	if err != nil {
		return err
	}
	// Cut off the unused bits
	response = response[:n]
	// Decode/Decrypt the data
	processed, err := l.decodePayload(response[:n], false)
	if err != nil {
		return err
	}
	// Parse the data and update the local struct
	err = l.parsePayload(processed)
	if err != nil {
		return err
	}
	socket.Close()
	return nil
}

func (l *Light) sendCommandPayload(data []byte) error {
	// Get the socket if it is still open, or open a new one
	socket, err := l.getSocket()
	if err != nil {
		return err
	}
	// Encrypt the payload
	payload, err := EncryptPacket([]byte(l.LocalId), data)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", data)
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
	// Increment the sequence number
	l.Seq++
	// Write the packet to the socket
	_, err = socket.Write(packet)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Craft a buffer to contain the incoming data
	re := make([]byte, 256)
	err = socket.SetReadDeadline(time.Now().Add(time.Millisecond * 250))
	if err != nil {
		return err
	}
	// Wait for and read the response from the device
	n, err := socket.Read(re)
	if err != nil {
		return err
	}
	_, err = l.decodePayload(re[:n], true)
	if err != nil {
		return err
	}

	socket.Close()
	// Disconnect from the device

	return nil
}

func (l *Light) getSocket() (net.Conn, error) {

	// Check if the socket is already connected
	//if l.socket != nil {
	//	l.Seq++
	//	return l.socket, nil
	//}
	//res := make([]byte, 0)
	//read, err := l.socket.Read(res)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(read)

	l.socket = nil
	l.Seq = 0
	var err error
	// Try to connect to the device
	l.socket, err = net.Dial("tcp", fmt.Sprintf("%s:%d", l.Ip, 6668))
	if err != nil {
		return nil, err
	}

	err = l.socket.SetWriteDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return nil, err
	}

	// Return opened socket
	return l.socket, nil
}

type System struct {
	Lights map[string]*Light
}

func (s *System) Scan(timeout time.Duration) error {
	pc, err := net.ListenPacket("udp4", ":6667")
	if err != nil {
		return err
	}
	defer pc.Close()
	//to := time.NewTimer(timeout)
	buf := make([]byte, 512)

	err = pc.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}

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

		if s.Lights[l.GwId] == nil {
			l.mutex = sync.Mutex{}
			s.Lights[l.GwId] = &l
		}
	}

}

type Tuya struct {
	plugin.Module
	entityId string
	devices  map[string]*Light
	entities map[string]string
	receiver chan domain.Attribute
	ready    bool
}

func init() {
	Module = Tuya{
		devices:  map[string]*Light{},
		entities: map[string]string{},
	}
	Module.Config = plugin.Config{
		Name:        "tuya",
		Type:        "daemon",
		Description: "Tuya Device interface",
		Version:     "0.1.1",
		Author:      "Braden Nicholson",
	}
}

func (t *Tuya) Setup() (plugin.Config, error) {
	err := t.UpdateInterval(1000 * 10)
	if err != nil {
		return plugin.Config{}, err
	}
	return t.Config, nil
}

func (t *Tuya) Update() error {
	if !t.ready {
		return nil
	}
	if t.Ready() {

		for _, light := range t.devices {

			eid, _ := t.findEntityIdByLightId(light.GwId)
			err := light.sendStatusPayload()
			if err != nil {
				return nil
			}
			if light.Latest.Dps.On {
				err = t.Attributes.Set(eid, "on", "true")
				if err != nil {
					log.Err(err)
				}
			} else {
				err = t.Attributes.Set(eid, "on", "false")
				if err != nil {
					log.Err(err)
				}
			}
			err = t.Attributes.Set(eid, "dim", fmt.Sprintf("%d", light.Latest.Dps.Dim/10))
			if err != nil {
				log.Err(err)
			}
			err = t.Attributes.Set(eid, "cct", fmt.Sprintf("%d", int(convertRange(float64(light.Latest.Dps.Cct), 0, 1000, 2700, 6500))))
			if err != nil {
				log.Err(err)
			}
			h := int(0)
			o := int(0)
			_, err = fmt.Sscanf(light.Latest.Dps.Color, "%04x%04x%04x", &h, &o, &o)
			if err != nil {
				return nil
			}
			_ = t.Attributes.Set(eid, "hue", fmt.Sprintf("%d", h))
			if err != nil {
				log.Err(err)
			}

		}

	}
	return nil
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
	light := t.devices[lightId]

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
	default:
		return nil
	}

	err = t.Attributes.Set(attribute.Entity, attribute.Key, attribute.Request)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tuya) mux() error {
	for {
		if !t.ready {
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
		default:

		}
	}
}

func (t *Tuya) Run() error {
	t.ready = false

	system := System{
		Lights: map[string]*Light{},
	}
	t.receiver = make(chan domain.Attribute, 8)
	go func() {
		err := t.mux()
		if err != nil {
			log.Err(fmt.Errorf("odd close of mux in tuya"))
			return
		}
	}()

	go func() {

		err := system.Scan(time.Second * 5)
		if err != nil {
			return
		}

		for _, light := range system.Lights {

			entity := &domain.Entity{
				Name:   light.GwId,
				Type:   "spectrum",
				Module: t.Config.Name,
			}

			light.LocalId = os.Getenv(light.GwId)

			err = t.Entities.Register(entity)
			if err != nil {
				log.Err(err)
				continue
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

		}
		t.devices = system.Lights
		t.ready = true
	}()

	return nil
}
