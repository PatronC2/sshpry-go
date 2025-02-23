# First stage: Run tests with TTY simulation
FROM golang:1.21.0 as build

WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o sshpry ./main.go

# Third stage: Copy the binary for host output
FROM alpine:latest AS output
WORKDIR /output
COPY --from=build /app/sshpry .