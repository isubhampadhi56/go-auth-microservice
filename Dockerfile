# Build stage
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools (optional: git for go mod download)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build the application
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
RUN go build -ldflags="-s -w" -o main ./cmd/main.go

# Runtime stage
FROM alpine:3.21

# Install CA certificates and tzdata for time zone support
RUN apk --no-cache add ca-certificates tzdata

# Optionally set the TZ env var so apps/commands see local timezone
ENV TZ=Asia/Kolkata

# (Optional) Configure localtime file
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .
# COPY --from=builder /app/.env.example .env

# Copy any other necessary files (if needed)
# COPY --from=builder /app/templates ./templates
# COPY --from=builder /app/static ./static

# Change ownership of the application files to appuser
RUN chown -R appuser:appgroup /root/

# Switch to non-root user
USER appuser

# Expose port (default: 8080, can be overridden with API_PORT environment variable)
# To publish the port when running: docker run -p 8080:8080 or docker run -p 3000:8080
EXPOSE 8080
ENV APP_ENV=PRODUCTION
# Command to run the application
CMD ["./main"]
