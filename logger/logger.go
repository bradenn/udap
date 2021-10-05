package logger

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/color"
	"net/http"
	"strings"
	"time"
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
)

func log(logType LogType, format string, args ...interface{}) {
	var prefix string

	switch logType {
	case info:
		prefix = fmt.Sprintf("%s[INFO]%s", color.BoldBlue, color.Reset)
		break
	case warn:
		prefix = fmt.Sprintf("%s[WARN]%s", color.BoldYellow, color.Reset)
		break
	case panic:
		prefix = fmt.Sprintf("%s[ERR*]%s", color.BoldRed, color.Reset+color.Red)
		break
	}
	var body string

	if format == "" {
		body = fmt.Sprint(args)
	} else {
		body = fmt.Sprintf(format, args...)
	}

	formatted := fmt.Sprintf("%s %s", prefix, body)
	strings.ReplaceAll(formatted, "Submit", "")
	fmt.Println(formatted)
}

func Middleware(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		defer func() {
			Info("%s%s %s'%s', %s", color.Green, r.Method, color.Reset, r.RequestURI, time.Since(t1))
		}()
		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Info(format string, args ...interface{}) {
	log(info, format, args...)
}

func Warn(format string, args ...interface{}) {
	log(warn, format, args...)
}

func Error(format string, args ...interface{}) {
	log(panic, format, args...)
}

func Err(err error) {
	log(panic, err.Error())
}
