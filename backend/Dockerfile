# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata wget \
    && adduser -D -H -u 10001 appuser

WORKDIR /app

COPY --from=builder /out/server /app/server

ENV TZ=Asia/Shanghai \
    APP_PORT=8080 \
    APP_ENV=production

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${APP_PORT}/health" || exit 1

CMD ["/app/server"]
