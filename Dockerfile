# Stage 1: Build frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm i && npm ci

COPY frontend/ ./
RUN cp public/config.js.template public/config.js
RUN npm run build

# Stage 2: Build backend
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git

COPY server/ .

RUN go build -buildvcs=false -o server .

# Stage 3: Runtime
FROM alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata gettext

COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /frontend/dist ./frontend/dist

COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

RUN mkdir -p /app/logs

EXPOSE 3002

ENTRYPOINT ["./docker-entrypoint.sh"]
