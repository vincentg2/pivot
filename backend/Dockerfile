# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend ./backend
RUN cd backend && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/pivot-api ./cmd/api \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/collect-catalog ./cmd/collect-catalog \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/collect-sport ./cmd/collect-sport \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/collect-tv ./cmd/collect-tv \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/collect-news ./cmd/collect-news
RUN GOBIN=/out go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/ /usr/local/bin/
COPY backend/migrations ./migrations
COPY backend/docker-entrypoint.sh /usr/local/bin/pivot-entrypoint
RUN chmod 0555 /usr/local/bin/pivot-entrypoint
EXPOSE 8080
ENTRYPOINT ["pivot-entrypoint"]
