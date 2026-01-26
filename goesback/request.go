package goesback

type RequestBody struct {
	SizeRaw          int64
	SizeDecompressed int64
	Raw              string
	Decompressed     string
	JSON             any
}

type Request struct {
	Port    int
	Method  string
	URL     string
	Scheme  string
	Host    string
	Path    string
	Body    *RequestBody
	Query   map[string][]string
	Headers map[string][]string
}

type Client struct {
	Port      int
	IP        string
	UserAgent string
}

type Server struct {
	Name    string
	Version string
	Host    string
}

type EchoResponse struct {
	Request *Request
	Client  *Client
	Server  *Server
}
