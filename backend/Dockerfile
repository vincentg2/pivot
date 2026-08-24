# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder
WORKDIR /src
RUN corepack enable
COPY pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/package.json
RUN pnpm install --frozen-lockfile
COPY frontend ./frontend
ARG VITE_REMOTE_LOGOS_ENABLED=true
ENV VITE_REMOTE_LOGOS_ENABLED=$VITE_REMOTE_LOGOS_ENABLED
RUN pnpm --filter @pivot/frontend build

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend ./backend
RUN cd backend && CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -trimpath -ldflags="-s -w" -o /out/pivot-api ./cmd/api \
    && CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -trimpath -ldflags="-s -w" -o /out/collect-catalog ./cmd/collect-catalog \
    && CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -trimpath -ldflags="-s -w" -o /out/collect-sport ./cmd/collect-sport \
    && CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -trimpath -ldflags="-s -w" -o /out/collect-tv ./cmd/collect-tv \
    && CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -trimpath -ldflags="-s -w" -o /out/collect-news ./cmd/collect-news
RUN CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go install github.com/pressly/goose/v3/cmd/goose@v3.24.1 \
    && find "$(go env GOPATH)/bin" -type f -name goose -exec cp {} /out/goose \; \
    && test -x /out/goose

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/ /usr/local/bin/
COPY --from=frontend-builder /src/frontend/dist ./public
COPY backend/migrations ./migrations
COPY backend/docker-entrypoint.sh /usr/local/bin/pivot-entrypoint
RUN chmod 0555 /usr/local/bin/pivot-entrypoint
EXPOSE 8080
ENV WEB_DIST_DIR=/app/public
ENTRYPOINT ["pivot-entrypoint"]
