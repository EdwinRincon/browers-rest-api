# syntax=docker/dockerfile:1.4

ARG TARGETARCH

# =========================
# STAGE 1: Builder
# =========================
FROM --platform=$BUILDPLATFORM golang:1.23.1-alpine AS builder

RUN apk add --no-cache git && \
    apk del --purge git && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=0
ARG TARGETARCH
RUN GOOS=linux GOARCH="$TARGETARCH" go build -ldflags="-s -w" -o /app/app . && \
    test -f /app/app || exit 1

# =========================
# STAGE 2: Production
# =========================
FROM --platform=$BUILDPLATFORM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup && \
    rm -rf /var/cache/apk/*

COPY --from=builder /app/app /usr/local/bin/app
WORKDIR /home/appuser
ENV GIN_MODE=release TZ=Europe/Madrid
EXPOSE 5050
USER appuser

ENTRYPOINT ["/usr/local/bin/app"]
