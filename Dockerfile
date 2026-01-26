FROM dhi.io/golang:1.24-debian12 AS builder

WORKDIR $GOPATH/src/github.com/jolaz-dev/goes-back
ADD . .

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o goes-back .

FROM dhi.io/debian-base:bookworm

WORKDIR /app
COPY --from=builder /go/src/github.com/jolaz-dev/goes-back/goes-back .

EXPOSE 8080

CMD ["./goes-back"]
