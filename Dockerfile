FROM golang:1.24.0 AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o sshpry ./main.go
