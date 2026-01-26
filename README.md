# goes-back

**goes-back** is a minimal HTTP echo server written in Go, partially inspired on Ealenn's [Echo-Server](https://github.com/Ealenn/Echo-Server). It accepts incoming HTTP requests and replies with a structured JSON response that reflects the original request — useful for testing clients and debugging requests.

## 🚀 Features

- Echoes request data back to the client in JSON format
- Returns method, path, query, headers, body, and client info
- Allows you to customize the response status, body and headers on a per-request basis
- Lightweight and minimal by design
- Easy to run locally or in a container

## 📦 Example JSON Response

```json
{
  "Request": {
    "BodySize": 20,
    "Port": 8080,
    "Method": "POST",
    "URL": "/hahaha?uno=due\u0026uno=tre",
    "Scheme": "http",
    "Host": "localhost:8080",
    "Path": "/hahaha",
    "RawBody": "{\n  \"key\": \"value\"\n}",
    "JSONBody": {
      "key": "value"
    },
    "Query": {
      "uno": [
        "due",
        "tre"
      ]
    },
    "Headers": {
      "Accept-Encoding": [
        "gzip, deflate"
      ],
      "Connection": [
        "close"
      ],
      "Content-Length": [
        "20"
      ],
      "Content-Type": [
        "application/json"
      ],
      "User-Agent": [
        "vscode-restclient"
      ]
    }
  },
  "Client": {
    "Port": 41218,
    "IP": "127.0.0.1",
    "UserAgent": "vscode-restclient"
  },
  "Server": {
    "Name": "goes-back",
    "Version": "0.0.1\n",
    "Host": "my-host"
  }
}
```

## 🧪 Getting Started

If you just want to use the server, there's a [Docker image](https://hub.docker.com/r/jonathanlazaro1/goes-back) that you can use with:
```bash
docker run -p 8080:8080 jonathanlazaro1/goes-back
```

If you want to contribute or prefer to run it locally, follow the instructions below.

### Prerequisites

- Go 1.24 or newer

### Run Locally
```bash
git clone https://github.com/jolaz-dev/goes-back
cd goes-back
go get
go run main.go
```


The server will start to listen on `:8080`.

### Using Docker

This project uses base [Docker Hardened Images](https://www.docker.com/products/hardened-images/), so you'll need a free Docker account in order to build its image locally. After [logging in](https://docs.docker.com/reference/cli/docker/login/) in both default and dhi.io registries,  you can build the image using:

`docker build -t goes-back .`


, then run the container using:

`docker run -p 8080:8080 goes-back`

## 📡 Example Requests
### Simple GET
`curl http://localhost:8080/hello?foo=bar`

### POST with Body
```bash
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -d '{ "msg": "hello" }'
  ```

### Customizing status code
```bash
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -H "X-GoesBack-Status: 202" \
  -d '{ "msg": "hello" }'
```

### Customizing headers
```bash
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -H "X-GoesBack-Header-Foo: Bar" \
  -d '{ "msg": "hello" }'
```

Each header starting with `X-GoesBack-Header-` will be added to the response headers, with the name coming after the dash becoming the name of the response header. In the example above, the response would contain a header `Foo` with the value `Bar`. Multiple headers with the same name will be treated as multiple values for the same header.

### Customizing response body
```bash
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -H "X-GoesBack-Body: {\"status\": \"accepted\"}" \
  -d '{ "msg": "hello" }'
```

For convenience, whenever you customize the response body, the original goes-back response will be returned as a header: `X-GoesBack-Response`.

## 🧑‍💻 Contributing

Contributions are welcome! Please fork the repo and open a pull request.

## 📄 License

This project is licensed under the MIT License.
