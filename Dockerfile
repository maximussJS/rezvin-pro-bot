FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main src/cmd/*.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/main /app/main

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 443

VOLUME ["/app/certs"]


CMD ["/app/main"]
