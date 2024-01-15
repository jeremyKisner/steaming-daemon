FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o bin/server ./cmd/server/main.go
EXPOSE 8080
CMD ["./bin/server"]
