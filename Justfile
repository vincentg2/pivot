set dotenv-load := true

default:
    @just --list

install:
    pnpm install
    cd backend && go mod download

doctor:
    @command -v go >/dev/null && go version
    @command -v node >/dev/null && node --version
    @command -v pnpm >/dev/null && pnpm --version
    @command -v docker >/dev/null && docker --version
    @command -v docker >/dev/null && docker compose version

db-up:
    docker compose up -d postgres

db-down:
    docker compose down

migrate:
    cd backend && go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 -dir migrations postgres "${DATABASE_URL}" up

dev:
    # Run both processes; Ctrl-C stops the process group.
    trap 'kill 0' INT TERM EXIT; (cd backend && go run github.com/air-verse/air@v1.62.0) & (cd frontend && pnpm dev) & wait

seed-demo:
    cd backend && go run ./cmd/seed-demo

collect-sport:
    cd backend && go run ./cmd/collect-sport

create-admin email nickname:
    cd backend && go run ./cmd/create-admin "{{email}}" "{{nickname}}"

fmt:
    cd backend && gofmt -w $$(find . -name '*.go' -not -path './tmp/*')
    pnpm --filter @pivot/frontend format

test:
    cd backend && go test ./...
    pnpm --filter @pivot/frontend test --run

lint:
    cd backend && go vet ./...
    pnpm --filter @pivot/frontend lint

build:
    cd backend && go build ./...
    pnpm --filter @pivot/frontend build

check: lint test build
