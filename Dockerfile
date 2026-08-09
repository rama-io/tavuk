# syntax=docker/dockerfile:1

FROM golang:1.26.5-alpine3.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/tavuk ./cmd/tavuk
RUN go test ./...

FROM alpine:3.24.1
RUN apk add --no-cache ca-certificates
RUN adduser -D -u 10001 app

ENV DATA_DIR=/data
RUN mkdir -p /data && chown app:app /data

WORKDIR /app
COPY --from=build /out/tavuk /app/tavuk
USER app
CMD ["/app/tavuk"]
