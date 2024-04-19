package cli

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	mutex                     = &sync.RWMutex{}
	clientsWithDefaultTimeout = map[string]*Client{}
)

type Client struct {
	Rest          *resty.Client
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	HTTPS         *bool  `json:"https"`
	ExternalToken string `json:"external_token"`
}

func dialTimeout(_ context.Context, network, addr string) (net.Conn, error) {
	timeout := 2 * time.Second
	return net.DialTimeout(network, addr, timeout)
}

var transport = http.Transport{
	DisableKeepAlives: true,
	DialContext:       dialTimeout,
}

func NewClientWithNormalTimeout(cli *Client) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	if cli == nil {
		return nil
	}
	baseURL := getBaseUrl(cli)
	if client, found := clientsWithDefaultTimeout[baseURL]; found {
		if composeToken(cli.ExternalToken) != client.Rest.Header.Get("Authorization") {
			client.Rest.SetHeader("Authorization", composeToken(cli.ExternalToken))
		}
		return client
	}
	rest := resty.New()
	rest.SetBaseURL(baseURL)
	rest.SetHeader("Authorization", composeToken(cli.ExternalToken))
	rest.SetTransport(&transport)
	rest.SetTimeout(10 * time.Minute)
	rest.SetAllowGetMethodPayload(true)
	cli.Rest = rest
	clientsWithDefaultTimeout[baseURL] = cli
	return cli
}

func composeToken(token string) string {
	return fmt.Sprintf("External %s", token)
}

func getBaseUrl(cli *Client) string {
	cli.Rest = resty.New()
	if cli.Ip == "" {
		cli.Ip = "0.0.0.0"
	}
	if cli.Port == 0 {
		cli.Port = 1660
	}
	var baseURL string
	if cli.HTTPS != nil && *cli.HTTPS {
		baseURL = fmt.Sprintf("https://%s:%d", cli.Ip, cli.Port)
	} else {
		baseURL = fmt.Sprintf("http://%s:%d", cli.Ip, cli.Port)
	}
	return baseURL
}
