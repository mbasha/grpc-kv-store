FROM golang:1.20 AS builder

WORKDIR /app

COPY ./rest-api-service/go.mod ./rest-api-service/go.sum ./
RUN go mod download

COPY ./rest-api-service/ ./
RUN go build -o rest-api-service ./cmd/main.go

FROM gcr.io/distroless/base

COPY --from=builder /app/rest-api-service /rest-api-service

ENTRYPOINT ["/rest-api-service"]