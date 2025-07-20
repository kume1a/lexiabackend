FROM golang:1.24.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache curl make ca-certificates
RUN curl -sSf https://atlasgo.sh | sh

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY ./internal ./internal
COPY ./ent ./ent
COPY Makefile ./

RUN make schemagen

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -trimpath \
    -o /lexiabin

FROM alpine:3.21

WORKDIR /app

ENV ENVIRONMENT=production

ARG DB_CONNECTION_URL
ENV DB_CONNECTION_URL=${DB_CONNECTION_URL}

RUN apk add --no-cache ca-certificates tzdata && \
    apk add --no-cache --virtual .atlas-deps curl && \
    curl -sSf https://atlasgo.sh | sh && \
    apk del .atlas-deps && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

COPY --from=builder /lexiabin /lexiabin
RUN chmod +x /lexiabin

COPY --from=builder /app/ent/migrate/migrations ./ent/migrate/migrations

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

USER appuser

EXPOSE 8002

ENTRYPOINT ["/entrypoint.sh"]
