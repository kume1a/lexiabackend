FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download

COPY *.go ./
COPY ./internal ./internal
COPY ./ent ./ent

RUN apk add --no-cache curl make && \
    curl -sSf https://atlasgo.sh | sh

RUN make schemagen

RUN CGO_ENABLED=0 GOOS=linux go build -o /lexiabin

FROM alpine:3.21

WORKDIR /app

ENV ENVIRONMENT=production
ENV PATH="/root/.local/bin:$PATH"

ARG DB_CONNECTION_URL
ENV DB_CONNECTION_URL=${DB_CONNECTION_URL}

COPY --from=builder /lexiabin /lexiabin
COPY --from=builder /app/ent/migrate/migrations ./ent/migrate/migrations
COPY Makefile ./

EXPOSE 8002

RUN apk add --no-cache curl make && \
    curl -sSf https://atlasgo.sh | sh

# ENTRYPOINT ["tail", "-f", "/dev/null"]

CMD ["/bin/sh", "-c", "make migration-apply && /lexiabin"]
