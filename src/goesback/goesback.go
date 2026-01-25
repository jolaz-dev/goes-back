package goesback

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func GoesBack(w http.ResponseWriter, r *http.Request) {
	clientRequest := Request{
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
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
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

	server := getServerData()

	resp := EchoResponse{
		Request: clientRequest,
		Client:  client,
		Server:  server,
	}

	output, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Error marshaling JSON:", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
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

func getClientData(r *http.Request) Client {
	host, portStr, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		slog.Error("Error splitting host and port:", "error", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Error("Error converting port:", "error", err)
	}

	return Client{
		IP:        host,
		Port:      port,
		UserAgent: r.UserAgent(),
	}
}

func getServerData() Server {
	version, err := os.ReadFile("version.txt")
	if err != nil {
		slog.Error("Error reading version file:", "error", err)
		version = []byte("unknown")
	}

	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Error getting hostname:", "error", err)
		hostname = "unknown"
	}

	return Server{
		Name:    "goes-back",
		Version: string(version),
		Host:    hostname,
	}
}
