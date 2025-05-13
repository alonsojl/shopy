FROM golang:1.21-alpine

WORKDIR /home/lalupe

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk update && apk add --no-cache build-base git zip curl
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
