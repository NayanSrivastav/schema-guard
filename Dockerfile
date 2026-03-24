FROM golang:alpine

WORKDIR /app

# Copy all source maps securely across
COPY . .

# Install dependencies locally over Alpine ensuring smooth compilation
RUN go mod tidy

# Build the native binary architecture cleanly
RUN go build -o schemaguard-api main.go

# Expose Telemetry and Engine Ports
EXPOSE 8080

# Execute strictly
CMD ["./schemaguard-api"]
