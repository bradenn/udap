// Copyright (c) 2021 Braden Nicholson

package log

import (
	"fmt"
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
)

func log(logType LogType, format string, args ...interface{}) {
	var prefix string

	switch logType {
	case module:
		prefix = fmt.Sprintf("%s[UDAP]%s", Reset+BoldBlue, Reset)
		break
	case panic:
		prefix = fmt.Sprintf("%s[UDAP]%s", Reset+BoldRed, Reset)
		break
	case sherlock:
		prefix = fmt.Sprintf("%s[SRLK]%s", Reset+BoldCyan, Reset)
		break
	}
	var body string

	if format == "" {
		body = fmt.Sprint(args)
	} else {
		body = fmt.Sprintf(format, args...)
	}

	formatted := fmt.Sprintf("%s %s", prefix, body)
	fmt.Println(formatted)
}

func Info(format string, args ...interface{}) {
	log(module, format, args...)
}

func LogF(format string, args ...interface{}) {
	log(module, format, args...)
}

func Log(format string, args ...interface{}) {
	log(module, format, args...)
}

func Sherlock(format string, args ...interface{}) {
	log(sherlock, format, args...)
}

func Warn(format string, args ...interface{}) {
	log(module, format, args...)
}

func Error(format string, args ...interface{}) {
	log(panic, format, args...)
}

func ErrF(err error, format string, args ...interface{}) {
	log(panic, fmt.Sprintf("%s\n	%s%s%s", format, Red, err.Error(), Reset), args...)
}

func Err(err error) {
	log(panic, err.Error(), "")
}
