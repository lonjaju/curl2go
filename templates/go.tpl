package main

import (
	"fmt"
	"log"
	"time"
{{- if .Data.Ascii }}
	"strings"
{{ end }}
	"net/http"
	"net/http/httputil"
)

func init() {
    log.SetFlags(log.Lshortfile | log.Ltime)
}

func main() {
	client := http.Client{
		Transport: &http.Transport{
			Proxy:                 nil,
			TLSHandshakeTimeout:   time.Second * 5,
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: time.Second * time.Duration(10),
			{{- if .Insecure }}
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // 不安全模式
			{{ end }}
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * time.Duration(10),
	}

	{{ if .Data.Ascii }}
	body := strings.NewReader(`{{ .Data.Ascii }}`)
	r, _ := http.NewRequest("{{ .Method }}", "{{ .URL }}", body)
	{{ else }}
	r, _ := http.NewRequest("{{ .Method }}", "{{ .URL }}", nil)
	{{ end -}}

	{{ range $k, $v := .Headers }}
	r.Header.Set("{{ $k }}", "{{ $v }}")
    {{- end }}

	resp, err := client.Do(r)
	if err != nil {
		log.Fatal("SendRequest", err.Error())
	}
	{{ if .Host }}
	r.Host = "{{ .Host }}"
	{{ end }}
	defer resp.Body.Close()
	data, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal("DumpResponse", err)
	}

	fmt.Println(string(data))
}
