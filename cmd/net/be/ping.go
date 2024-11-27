package be

import (
	"net"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 2 * time.Second,
}

func Ping(target string) (int, error) {
	var url string

	// Check if the target is an IP address
	if ip := net.ParseIP(target); ip != nil {
		// It's an IP address, so use it directly for pinging
		url = "http://" + target
	} else {
		// It's not an IP address, assume it's a domain name
		url = "http://" + target
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
