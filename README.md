# streaming daemon api

# Local Development

## Docker
1. Launch Application
```
docker compose up --build -d
```
2. Insert Record
```
go run ./cmd/audio/insert/main.go -name "to be titled" -artist "me" -album "album title" -filepath examples/beep.wav
```
3. Check Record. Go to [http://localhost:8080/audio/1](http://localhost:8080/audio/1)
4. Tear down
```
docker compose down
```
