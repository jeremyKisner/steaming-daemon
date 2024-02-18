# streaming daemon api

# Local Development

## Docker
1. Stand up the infrastructure.
```
docker compose up -d
```
1. Build application
```
*docker build -t streaming-daemon .*
```
1. Create Tables
```
go run cmd/scripts/createtable.go
```
1. Run the Server
```
docker run -p 8080:8080 streaming-daemon
```
1. Tear down
```
docker compose down
```
=

## Run Locally
```
go run cmd/server/main.go
```