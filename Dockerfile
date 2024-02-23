FROM golang:latest
WORKDIR /app
COPY . .
RUN mkdir audio
RUN go build -o bin/server ./cmd/server/main.go
EXPOSE 8080
CMD ["./bin/server"]
