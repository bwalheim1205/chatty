
# --------------------------
# BUILD STAGE
# --------------------------
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build source
COPY cmd/ cmd/
COPY internal/ internal/
RUN go build -o chatty ./cmd/cli

# --------------------------
# RUNTIME STAGE
# --------------------------
FROM alpine:3.18

# Add CA certs
RUN apk add --no-cache ca-certificates

# Copy binary from builder
WORKDIR /app
COPY --from=builder /app/chatty .


ENTRYPOINT [ "./chatty" ]
