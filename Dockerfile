# Start from the official Go image
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o exchange-rate-service ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/exchange-rate-service .
CMD ["./exchange-rate-service"]
