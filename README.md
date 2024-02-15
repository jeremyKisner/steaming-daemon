# streaming daemon api

# Local Development

## Docker Compose
Stand up the infrastructure.
```
docker compose up
```
```
docker compose down
```

## Run Locally
```
go run cmd/server/main.go
```

## Build Docker Image
```
docker build -t streaming-daemon .
```

## Run Docker Image
```
docker run -p 8080:8080 streaming-daemon
```
