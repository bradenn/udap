// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"time"
	"udap/internal/log"
	"udap/internal/pkg/dmx"
	"udap/internal/pkg/dmx/ft232"
	"udap/pkg/plugin"
)

var Plugin Lor

type Lor struct {
	plugin.SDK
	requests chan plugin.Request
	resolver chan plugin.Event
	metadata plugin.Metadata
}

var morseSymbols = map[string]string{
	"A":  ".-",
	"B":  "-...",
	"C":  "-.-.",
	"D":  "-..",
	"E":  ".",
	"F":  "..-.",
	"G":  "--.",
	"H":  "....",
	"I":  "..",
	"J":  ".---",
	"K":  "-.-",
	"L":  ".-..",
	"M":  "--",
	"N":  "-.",
	"O":  "---",
	"P":  ".--.",
	"Q":  "--.-",
	"R":  ".-.",
	"S":  "...",
	"T":  "-",
	"U":  "..-",
	"V":  "...-",
	"W":  ".--",
	"X":  "-..-",
	"Y":  "-.--",
	"Z":  "--..",
	"1":  ".----",
	"2":  "..---",
	"3":  "...--",
	"4":  "....-",
	"5":  ".....",
	"6":  "-....",
	"7":  "--...",
	"8":  "---..",
	"9":  "----.",
	"0":  "-----",
	".":  ".-.-.-",  // period
	":":  "---...",  // colon
	",":  "--..--",  // comma
	";":  "-.-.-",   // semicolon
	"?":  "..--..",  // question
	"=":  "-...-",   // equals
	"'":  ".----.",  // apostrophe
	"/":  "-..-.",   // slash
	"!":  "-.-.--",  // exclamation
	"-":  "-....-",  // dash
	"_":  "..--.-",  // underline
	"\"": ".-..-.",  // quotation marks
	"(":  "-.--.",   // parenthesis (open)
	")":  "-.--.-",  // parenthesis (close)
	"()": "-.--.-",  // parentheses
	"$":  "...-..-", // dollar
	"&":  ".-...",   // ampersand
	"@":  ".--.-.",  // at
	"+":  ".-.-.",   // plus
	"Á":  ".--.-",   // A with acute accent
	"Ä":  ".-.-",    // A with diaeresis
	"É":  "..-..",   // E with acute accent
	"Ñ":  "--.--",   // N with tilde
	"Ö":  "---.",    // O with diaeresis
	"Ü":  "..--",    // U with diaeresis
	" ":  " ",       // word interval
}

func EntityHandler(payload string) error {
	fmt.Println(payload)
	return nil
}

func (s *Lor) Startup() (plugin.Metadata, error) {
	s.SDK = plugin.SDK{}
	Plugin = Lor{
		metadata: plugin.Metadata{
			Name:        "Lor",
			Type:        "module",
			Description: "Lor 16-Channel 120VAC 15A/15A 10-Bit Dimmer",
			Version:     "0.0.1",
			Author:      "Braden Nicholson",
		},
		requests: make(chan plugin.Request),
	}

	config := dmx.NewConfig(0x02)

	config.GetUSBContext()
	// Since ft232 is a shitty module, it panics when USB can't be found.
	defer func() {
		recover()
	}()
	// for i := 0; i < 16; i++ {
	// 	err := s.SDK.CreateOrInitEntity(fmt.Sprintf("Outlet %d", i+1), plugin.DIMMER, EntityHandler)
	// 	if err != nil {
	// 		return plugin.Metadata{}, err
	// 	}
	// }
	controller := ft232.NewDMXController(config)
	if err := controller.Connect(); err != nil {
		fmt.Printf("failed to connect DMX Controller: %s\n", err)
		log.Err(err)
	}

	go func() {
		defer func() {
			err := controller.Close()
			if err != nil {
				return
			}
		}()
		var sequence []rune

		morseStr := "I LOVE YOU"
		// delay := 0
		log.Sherlock("Running Morse: '%s'", morseStr)
		// For each letter of the string
		for i := range morseStr {
			morse := morseSymbols[string(morseStr[i])]
			for _, symbol := range morse {
				sequence = append(sequence, symbol)
			}
		}

		fmt.Println(sequence)

		// place := 0
		for {

			var err error

			err = controller.SetChannel(1, byte(255))
			if err != nil {
				fmt.Println(err)
			}

			err = controller.Render()
			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(time.Millisecond * 250)
			// _ = controller.Render()
			// switch r := sequence[place]; r {
			// case '.':
			// 	time.Sleep(time.Millisecond * 250)
			// 	break
			// case '-':
			// 	time.Sleep(time.Millisecond * 500)
			// 	break
			// case ' ':
			// 	time.Sleep(time.Millisecond * 250)
			// 	break
			// }

			// time.Sleep(time.Millisecond * 100)

		}
	}()

	return s.metadata, nil
}

func (s *Lor) Metadata() plugin.Metadata {
	return s.metadata
}

func (s *Lor) Listen() {
	for request := range s.requests {
		s.handle(request)
	}
	close(s.requests)
}

func (s *Lor) handle(req plugin.Request) {

}

func (s *Lor) Connect(events chan plugin.Event) chan plugin.Request {
	s.resolver = events
	return s.requests
}

func (s *Lor) resolve(event plugin.Event) {
	s.resolver <- event
}

func (s *Lor) Cleanup() {
	close(s.requests)
}
