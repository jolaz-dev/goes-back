package goesback

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jolaz-dev/goes-back/config"
	"github.com/jolaz-dev/goes-back/internal"
)

func GoesBack(r *http.Request, config *config.Config) (*EchoResponse, error) {
	clientRequest := &Request{
		BodySize: r.ContentLength,
		Method:   r.Method,
		Headers:  r.Header,
		URL:      r.RequestURI,
	}

	scheme := r.URL.Scheme
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	clientRequest.Scheme = scheme

	u, err := url.ParseRequestURI(scheme + "://" + r.Host + r.RequestURI)
	if err != nil {
		slog.Error("Error parsing URL:", "error", err)
		return nil, internal.ErrUnprocessableEntity
	}

	clientRequest.Scheme = u.Scheme
	clientRequest.Host = u.Host
	clientRequest.Path = u.Path
	clientRequest.Query = u.Query()

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		slog.Error("Error converting port:", "error", err)
	}
	clientRequest.Port = port

	rawBody := make([]byte, r.ContentLength)
	_, err = r.Body.Read(rawBody)
	if err != nil && err.Error() != "EOF" {
		slog.Error("Error reading body:", "error", err)
	}
	clientRequest.RawBody = string(rawBody)

	clientRequest.JSONBody = tryAndParseJSONBody(r, rawBody)

	client := getClientData(r)

	server := getServerData(config)

	resp := EchoResponse{
		Request: clientRequest,
		Client:  client,
		Server:  server,
	}

	return &resp, nil
}

func tryAndParseJSONBody(r *http.Request, body []byte) any {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil
	}

	var jsonBody any
	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		slog.Error("Error trying to parse JSON body", "error", err)
		return nil
	}
	return jsonBody
}

func getClientData(r *http.Request) *Client {
	host, portStr, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		slog.Error("Error splitting host and port:", "error", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Error("Error converting port:", "error", err)
	}

	return &Client{
		IP:        host,
		Port:      port,
		UserAgent: r.UserAgent(),
	}
}

func getServerData(config *config.Config) *Server {
	return &Server{
		Name:    config.AppName,
		Version: config.Version,
		Host:    config.Hostname,
	}
}
