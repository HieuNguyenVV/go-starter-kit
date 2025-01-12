FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY ./config.yaml ./config.yaml
COPY --from=builder /app/main /app/main
ENTRYPOINT ["./main"]
