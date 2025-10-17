package ylog

import (
	js "github.com/bitly/go-simplejson"
)

const (
	ostd = 1 << iota
	odbg
	otrace
)

func otoPrefix(f int) string {
	var s string
	switch f {
	case odbg:
		s = "[debug] "
	case otrace:
		s = "[trace] "
	default:
		return ""
	}
	return s
}

func createJSONLogWithJs(code int, message string, data *js.Json) *js.Json {
	j := js.New()
	j.Set("code", code)
	j.Set("message", message)
	if data != nil {
		j.Set("data", data)
	} else {
		j.Set("data", "")
	}
	return j
}

func createJSONLog(code int, message string, data interface{}) string {
	log := js.New()
	log.Set("code", code)
	log.Set("message", message)
	if data != nil {
		log.Set("data", data)
	} else {
		log.Set("data", "")
	}
	b, _ := log.MarshalJSON()
	return string(b)
}

func correctProg100(prog string) string {
	if prog == "100.00" || prog == "100.0" {
		prog = "100"
	}
	if prog == "100.00%" || prog == "100.0%" {
		prog = "100%"
	}
	if prog == "0.00" || prog == "0.0" {
		prog = "0"
	}
	if prog == "0.00%" || prog == "0.0%" {
		prog = "0%"
	}
	return prog
}
