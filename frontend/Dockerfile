# syntax=docker/dockerfile:1
FROM node:22-alpine AS builder
WORKDIR /src
RUN corepack enable
COPY pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/package.json
RUN pnpm install --frozen-lockfile
COPY frontend ./frontend
ARG VITE_REMOTE_LOGOS_ENABLED=false
ENV VITE_REMOTE_LOGOS_ENABLED=$VITE_REMOTE_LOGOS_ENABLED
RUN pnpm --filter @pivot/frontend build

FROM nginx:1.29-alpine
COPY frontend/nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --from=builder /src/frontend/dist /usr/share/nginx/html
ENV API_UPSTREAM=http://api:8080
EXPOSE 80
