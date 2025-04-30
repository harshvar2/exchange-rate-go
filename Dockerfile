# Start from the official Go image
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o exchange-rate-service ./cmd/server

FROM alpine:latest
WORKDIR /root/
WORKDIR /app
COPY --from=build /app /app
COPY --from=build /app/config.yml.example /app/config.yml.example
WORKDIR /app
CMD ["./exchange-rate-service"]
