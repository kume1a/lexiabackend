# Optimized Dockerfile - should reduce size from 280MB to ~50-80MB
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Install only essential build dependencies
RUN apk add --no-cache curl make ca-certificates

# Install atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./
COPY ./internal ./internal
COPY ./ent ./ent
COPY Makefile ./

# Generate schema
RUN make schemagen

# Build optimized binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -trimpath \
    -o /lexiabin

# Production stage - minimal alpine
FROM alpine:3.21

WORKDIR /app

ENV ENVIRONMENT=production

ARG DB_CONNECTION_URL
ENV DB_CONNECTION_URL=${DB_CONNECTION_URL}

# Install minimal runtime dependencies and atlas
RUN apk add --no-cache ca-certificates tzdata && \
    apk add --no-cache --virtual .atlas-deps curl && \
    curl -sSf https://atlasgo.sh | sh && \
    apk del .atlas-deps && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy optimized binary
COPY --from=builder /lexiabin /lexiabin
RUN chmod +x /lexiabin

# Copy migration files
COPY --from=builder /app/ent/migrate/migrations ./ent/migrate/migrations

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Use non-root user for security
USER appuser

EXPOSE 8002

ENTRYPOINT ["/entrypoint.sh"]
