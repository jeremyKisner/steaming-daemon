FROM golang:latest
WORKDIR /app
COPY . .
RUN mkdir audio
RUN go build -o bin/server ./cmd/server/main.go
EXPOSE 8082
CMD ["./bin/server"]
