FROM golang:1.22-alpine AS builder

# Install git and certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the actual code
COPY . .

# Build a statically linked binary (no C dependencies)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ticketrush cmd/api/main.go


FROM alpine:latest

WORKDIR /app

# Import necessary files from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/ticketrush .

COPY migrations/ ./migrations/

EXPOSE 8080

# Run the binary
CMD ["./ticketrush"]