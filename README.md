# Exchange Rate Service

A backend service in Go (using go-kit) for real-time and historical currency exchange rates.

## Features
- Fetch latest and historical exchange rates
- Currency conversion
- In-memory caching
- Clean architecture
- Dockerized

## Environment Variables

- `EXCHANGE_API_KEY`: Your exchangerate.host API key (required)

## How to Run

1. Set your API key environment variable:
   - On Windows (cmd):
     ```cmd
     set EXCHANGE_API_KEY=your_api_key_here
     ```
   - On Linux/macOS:
     ```sh
     export EXCHANGE_API_KEY=your_api_key_here
     ```

2. Tidy dependencies:
   ```sh
   go mod tidy
   ```

3. Build the server:
   ```sh
   go build -o exchange-rate-service ./cmd/server
   ```

4. Run the server:
   ```sh
   ./exchange-rate-service
   ```
   By default, it runs on port 8080. To use a different port:
   ```sh
   set PORT=9090
   ./exchange-rate-service
   ```

5. Test the endpoints (examples):
   - Convert currency:
     ```sh
     curl "http://localhost:8080/convert?from=USD&to=INR&amount=100"
     ```
   - Get exchange rate:
     ```sh
     curl "http://localhost:8080/rate?from=USD&to=INR"
     ```
   - Get historical rates:
     ```sh
     curl "http://localhost:8080/history?from=USD&to=INR&start=2025-01-01&end=2025-01-10"
     ```

## Docker

1. Build the Docker image:
   ```sh
   docker build -t exchange-rate-service .
   ```
2. Run the container:
   ```sh
   docker run -e EXCHANGE_API_KEY=your_api_key_here -p 8080:8080 exchange-rate-service
   ```


