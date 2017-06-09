package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
)

type Echo struct {
	Args struct {
	} `json:"args"`
	Headers Headers `json:"headers"`
	Origin  string  `json:"origin"`
	URL     string  `json:"url"`
}

type Headers struct {
	Accept                  string `json:"Accept"`
	AcceptEncoding          string `json:"Accept-Encoding"`
	AcceptLanguage          string `json:"Accept-Language"`
	Connection              string `json:"Connection"`
	Cookie                  string `json:"Cookie"`
	Host                    string `json:"Host"`
	Referer                 string `json:"Referer"`
	UpgradeInsecureRequests string `json:"Upgrade-Insecure-Requests"`
	UserAgent               string `json:"User-Agent"`
}

type HttpBin struct {
	sling *sling.Sling
}

func NewService(client *http.Client) *HttpBin {
	return &HttpBin{
		sling: sling.New().Client(client).Base("http://httpbin.org"),
	}
}

func (h *HttpBin) Get() (Echo, *http.Response, error) {
	echo := Echo{}
	res, err := h.sling.New().Get("/get").ReceiveSuccess(&echo)
	return echo, res, err
}

func main() {
	// Do a simple GET
	echo := Echo{}
	_, _ = sling.New().Get("https://httpbin.org/get").ReceiveSuccess(&echo)
	fmt.Println("User agent is", echo.Headers.UserAgent)

	httpBin := NewService(&http.Client{})
	resp, _, _ := httpBin.Get()

	fmt.Println(resp)

}
