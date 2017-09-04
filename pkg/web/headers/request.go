package headers

import (
	"fmt"
	"net/http"
	"strings"
)

func ParseBaseUrl(req *http.Request) (string, error) {
	// Default to the Host header for the hostname
	// Use the X-Forwarded-Host header value if present
	splitHost := strings.Split(req.Host, ":")
	if ho := req.Header.Get("X-Forwarded-Host"); ho != "" {
		splitHost = strings.Split(ho, ":")
	}
	host := splitHost[0]

	// Default to https
	// Use the X-Forwarded-Proto header value if present
	proto := "https"
	if pr := req.Header.Get("X-Forwarded-Proto"); pr != "" {
		proto = pr
	}

	// Default to port 443 for https, 80 for http
	// Use the port from the Host header if present
	// Use the X-Forwarded-Port header value if present
	port := "443"
	if proto == "http" {
		port = "80"
	}
	if len(splitHost) > 1 {
		port = splitHost[1]
	}
	if po := req.Header.Get("X-Forwarded-Port"); po != "" {
		port = po
	}

	// Build up a base url in the format: <protocol>://<host>[:port]
	baseUrl := fmt.Sprintf("%s://%s", proto, host)
	if port != "443" && port != "80" {
		baseUrl += ":" + port
	}
	return baseUrl, nil
}
