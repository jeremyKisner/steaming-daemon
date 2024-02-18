# streaming daemon api

# Local Development

## Docker
1. Launch application
```
docker compose up --build -d
```
1. Insert Record
```
go run ./cmd/audio/insert/main.go -name "to be titled" -artist "me" -album "album title"
```

1. Tear down
```
docker compose down
```
