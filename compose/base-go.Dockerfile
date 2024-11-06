FROM golang:1.23-alpine

COPY ../../pools-listener/go.mod ./
COPY ../../pools-listener/go.sum ./

RUN go mod download

COPY ../../abis .
COPY ../../pools-listener .

RUN go build -o myapp cmd/pool_listener/main.go