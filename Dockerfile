# Start from the official Go image
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o exchange-rate-service ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/exchange-rate-service .
COPY --from=build /app/config.yml.example ./config.yml.example
# Ensure config.yml exists, copy from example if missing
RUN if [ ! -f ./config.yml ]; then cp ./config.yml.example ./config.yml; fi
CMD ["./exchange-rate-service"]
