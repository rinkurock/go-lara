package conn

import (
	"net"
	"net/http"
	"time"
)

var sampleClient *http.Client

func ConnectSample(timeout time.Duration) {
	var netTransport = &http.Transport{
		DialContext:         (&net.Dialer{Timeout: timeout, KeepAlive: time.Minute}).DialContext,
		TLSHandshakeTimeout: timeout,
		MaxIdleConnsPerHost: 5,
	}

	sampleClient = &http.Client{
		Timeout:   timeout,
		Transport: netTransport,
	}
}
func GetSampleCon() *http.Client {
	return sampleClient
}