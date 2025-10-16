package tools

import (
	"net/http"
	"net/url"
	"strings"
)

func RedirectURL(url string) (r string) {
	defer func() {
		if r == "" {
			r = url
		}
	}()

	client := http.Client{
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			r = req.URL.String()
			return nil
		},
	}
	rsp, err := client.Head(url)
	if err != nil {
		r = ""
		return
	}
	rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		r = ""
		return
	}
	return
}

func GetURLQuery(u string, kay string, unescape bool) string {
	pu, e := url.Parse(u)
	if e != nil {
		return ""
	}
	val := pu.Query().Get(kay)
	if val != "" && unescape {
		val, e = url.QueryUnescape(val)
		if e != nil {
			return ""
		}
	}
	return val
}

func GetSchemeHost(u string) string {
	pu, e := url.Parse(u)
	if e != nil {
		return ""
	}
	var buf strings.Builder
	if pu.Scheme != "" {
		buf.WriteString(pu.Scheme)
		buf.WriteByte(':')
	}
	buf.WriteString("//")
	buf.WriteString(pu.Host)
	return buf.String()
}
