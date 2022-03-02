// Copyright (c) 2021 Braden Nicholson

package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	Normal       = ""
	Reset        = "\033[m"
	Bold         = "\033[1m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Magenta      = "\033[35m"
	Cyan         = "\033[36m"
	BoldRed      = "\033[1;31m"
	BoldGreen    = "\033[1;32m"
	BoldYellow   = "\033[1;33m"
	BoldBlue     = "\033[1;34m"
	BoldMagenta  = "\033[1;35m"
	BoldCyan     = "\033[1;36m"
	FaintRed     = "\033[2;31m"
	FaintGreen   = "\033[2;32m"
	FaintYellow  = "\033[2;33m"
	FaintBlue    = "\033[2;34m"
	FaintMagenta = "\033[2;35m"
	FaintCyan    = "\033[2;36m"
	BgRed        = "\033[41m"
	BgGreen      = "\033[42m"
	BgYellow     = "\033[43m"
	BgBlue       = "\033[44m"
	BgMagenta    = "\033[45m"
	BgCyan       = "\033[46m"
	Faint        = "\033[2m"
	FaintItalic  = "\033[2;3m"
	Reverse      = "\033[7m"
)

type LogType int

const (
	info = iota
	warn
	panic
	module
	sherlock
	in
	out
	partial
	process
	event
	critical
)

func log(logType LogType, format string, args ...interface{}) {
	var prefix string
	ln := true
	switch logType {
	case module:
		prefix = fmt.Sprintf("%s[UDAP]%s", Reset+BoldBlue, Reset)
		break
	case process:
		prefix = fmt.Sprintf("%s[UDAP]%s", Reset+BoldBlue, Reset)
		ln = false
		break
	case event:
		prefix = fmt.Sprintf("%s[EVENT]%s", Reset+BoldGreen, Reset+Faint)
		ln = true
		break
	case critical:
		prefix = fmt.Sprintf("%s[ACTION REQUIRED]%s", Reset+BoldMagenta, Reset)
		break
	case partial:
		prefix = fmt.Sprintf("%s", Green)
		ln = false
		break
	case panic:
		prefix = fmt.Sprintf("%s[UDAP]%s", Reset+BoldRed, Reset)
		break
	case sherlock:
		prefix = fmt.Sprintf("%s[SRLK]%s", Reset+BoldCyan, Reset)
		break
	case in:
		prefix = fmt.Sprintf("%s[<IN ]%s", Reset+BoldCyan, Reset)
		break
	case out:
		prefix = fmt.Sprintf("%s[OUT>]%s", Reset+BoldCyan, Reset)
		break
	}
	var body string

	if format == "" {
		body = fmt.Sprint(args)
	} else {
		body = fmt.Sprintf(format, args...)
	}

	formatted := fmt.Sprintf("%s %s", prefix, body)
	if ln {
		fmt.Println(formatted)
	} else {
		fmt.Printf(formatted)
	}

}

func Critical(format string, args ...interface{}) {
	log(critical, format, args...)
}

func Log(format string, args ...interface{}) {
	log(module, format, args...)
}

func Event(format string, args ...interface{}) {
	// _, ln, _, ok := runtime.Caller(1)
	tag := fmt.Sprintf("%s%s", Reset+BoldGreen, "[EVNT]")
	color := Reset + Faint
	switch t := strings.Split(format, " "); strings.ToLower(t[0]) {
	case "entity":
		color = Reset
		break
	case "module":
		color = Reset
		break
	}
	fmt.Printf("%s%s %s%s\n", tag, color, fmt.Sprintf(format, args...), Reset)
}

func ErrF(err error, format string, args ...interface{}) {
	log(panic, fmt.Sprintf("%s\n	%s%s%s", format, Red, err.Error(), Reset), args...)
}

func Err(err error) {
	_, file, ln, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("%s%s%s %s\n", Reset+BoldRed, fmt.Sprintf("Error (%s:%d)", filepath.Base(file), ln), Reset+FaintRed,
			fmt.Sprintf(err.Error()))
	}
}
