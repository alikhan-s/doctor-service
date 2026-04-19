# ==========================================
# STAGE 1: Builder
# ==========================================
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# --- DEPENDENCY CACHING STEP ---
COPY go.mod go.sum ./
RUN go mod download

# --- COPY ALL SERVICE ---
COPY . .

# --- COMPILATION ---
# CGO_ENABLED=0           - Disable CGO to get a clean static binary (required for running on Alpine/Scratch).
# GOOS=linux GOARCH=amd64 - Compile for Linux.
# -ldflags="-s -w"        - We remove debug information to reduce the size of the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o doctor-service ./cmd/doctor-s/main.go

# ==========================================
# STAGE 2: Runner
# ==========================================
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

# SECURITY: Creating a user and a group without root privileges
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the compiled binary from the "builder" stage
COPY --from=builder /app/doctor-service .

# Changing file user to our non-root
RUN chown appuser:appgroup /app/doctor-service

# Switch to non-root user
USER appuser

# Port
EXPOSE 8081

# Running comand
CMD ["./doctor-service"]