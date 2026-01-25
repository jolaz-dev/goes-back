package goesback

type Request struct {
	BodySize int64
	Port     int
	Method   string
	URL      string
	Scheme   string
	Host     string
	Path     string
	RawBody  string
	JSONBody any
	Query    map[string][]string
	Headers  map[string][]string
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
	Request Request
	Client  Client
	Server  Server
}
