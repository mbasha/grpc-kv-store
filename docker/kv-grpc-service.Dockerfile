FROM golang:1.18 AS builder

WORKDIR /app

COPY ./kv-grpc-service/go.mod ./go.mod
COPY ./kv-grpc-service/go.sum ./go.sum
RUN go mod download

COPY ./kv-grpc-service/ ./ 

RUN go build -o kv-grpc-service ./cmd/main.go

FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=builder /app/kv-grpc-service .

CMD ["./kv-grpc-service"]