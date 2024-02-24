FROM golang:latest
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server cmd/server.go
ENTRYPOINT [ "./server" ]
