# goes-back

**goes-back** is a minimal HTTP echo server written in Go, partially inspired on Ealenn's [Echo-Server](https://github.com/Ealenn/Echo-Server). It accepts incoming HTTP requests and replies with a structured JSON response that reflects the original request — useful for testing clients and debugging requests.

## 🚀 Features

- Echoes request data back to the client in JSON format
- Returns method, path, query, headers, body, and client info
- Allows you to customize the response status, body and headers on a per-request basis
- Supports request/response gzip (de)compression
- Lightweight and minimal by design
- Easy to run locally or in a container

## ❓ Why use it instead of Echo-Server?

First things first: I didn't write `goes-back` to "compete" with `Echo-Server`. I was inspired by it and just wanted to make a coding exercise by writing the same thing in Go. At the end of the day, both are open source projects and you can use one or even both of them the best way it suits you.


With that said, here goes some info in case you need to choose. For one side `goes-back` is currently not in feature-parity with `Echo-Server`: we don't support (yet) request delays, for example. Also, `Echo-Server` allows you to enable/disabling response customization using env vars, a feature that I particularly don't find that useful, but some folks out there may be used to have.

For other side, there are some potential good reasons for you to choose `goes-back`:
- Maybe you root more for Go than for JS 😅
- You may find our header customization support a little easier to understand than `Echo-Server`'s
- We support request/response gzip (de)compression
  - Trying to send compressed request bodies to Echo-Server logs an "incorrect header check" error if I try to use gzip. Other compression formats prompt a `415 Unsupported Media Type` response
  - I didn't test response compression, though
- Unfortunately, `Echo-Server` doesn't see the bright light of new code for a while now (as of jan/2026)
  - This is bad not only for the tracking of issues like the compression one mentioned above, but also for the security vulnerabilities that naturally pile up over time, which can prevent it from being used on environments with strict security rules

Again, this is not by any means an attack on `Ealenn`'s work: `Echo-Server` helped me a lot of times to test API gateway configuration against a "real" upstream. I probably would try and open a PR for the things I need if I wasn't unsure about the status of the project (something that the author surely has their own personal reasons to not acting about). It's the mix of doing a real life side project and the specific need for things like (de)compression that motivated me to create `goes-back`.

If any of this reasons makes sense for you to use `goes-back`, welcome aboard!

## 📦 Example JSON Response
Request:
```bash
curl -X POST "http://localhost:8080/hahaha?uno=due&uno=tre" \
    -H "Content-Type: application/json" \
    -H "Content-Encoding: gzip" \
    --data-binary "H4sIAAAAAAAAA6rmUlBQyk6tVLJSUCpLzClNVeKqBQAAAP//AwA+ENryFAAAAA=="
```

Response:
```json
{
  "Request": {
    "Port": 8080,
    "Method": "POST",
    "URL": "/hahaha?uno=due&uno=tre",
    "Scheme": "http",
    "Host": "localhost:8080",
    "Path": "/hahaha",
    "Body": {
      "SizeRaw": 64,
      "SizeDecompressed": 20,
      "Raw": "H4sIAAAAAAAAA6rmUlBQyk6tVLJSUCpLzClNVeKqBQAAAP//AwA+ENryFAAAAA==",
      "Decompressed": "{\n  \"key\": \"value\"\n}",
      "JSON": {
        "key": "value"
      }
    },
    "Query": {
      "uno": [
        "due",
        "tre"
      ]
    },
    "Headers": {
      "Accept-Encoding": [
        "*/*"
      ],
      "Connection": [
        "close"
      ],
      "Content-Length": [
        "64"
      ],
      "Content-Type": [
        "application/json"
      ],
      "Content-Encoding": [
        "gzip"
      ],
      "User-Agent": [
        "curl/8.15.0"
      ]
    }
  },
  "Client": {
    "Port": 41218,
    "IP": "127.0.0.1",
    "UserAgent": "vscode-restclient"
  },
  "Server": {
    "Name": "GoesBack",
    "Version": "0.2.1", // x-release-please-version
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
