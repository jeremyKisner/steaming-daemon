# streaming daemon api

# Local Development

## Run Docker
1. Launch Application
```
docker compose up --build -d
```
2. Insert Record
```
go run ./cmd/audioinsert/main.go -name "to be titled" -artist "me" -album "album title" -filepath examples/song.wav
```
3. Check Record. Go to [http://localhost:8080/audio/1](http://localhost:8080/audio/1)
4. Tear down
```
docker compose down
```
