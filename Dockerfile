# syntax=docker/dockerfile:1.4

ARG TARGETARCH=amd64
ARG BUILDPLATFORM=linux/amd64

# =========================
# STAGE 1: Builder
# =========================
FROM --platform=$BUILDPLATFORM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache git && go mod download && apk del git


# Copy explicitly the source files
# to ensure that the build context is clean and only contains necessary files
COPY cmd/      ./cmd/
COPY api/      ./api/
COPY config/   ./config/
COPY docs/     ./docs/
COPY helper/   ./helper/
COPY pkg/      ./pkg/
COPY server/   ./server/


ENV CGO_ENABLED=0
ARG TARGETARCH
RUN GOOS=linux GOARCH="$TARGETARCH" go build -ldflags="-s -w" -o /app/app ./cmd/browersfc && \
    test -f /app/app || exit 1

# =========================
# STAGE 2: Production
# =========================
FROM --platform=$BUILDPLATFORM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/app /usr/local/bin/app

WORKDIR /home/appuser

ENV GIN_MODE=release TZ=Europe/Madrid

EXPOSE 5050

USER appuser

ENTRYPOINT ["/usr/local/bin/app"]

