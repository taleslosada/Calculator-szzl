# ---- 1. Build the React frontend ----
FROM node:22-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---- 2. Build the Go backend (static binary) ----
FROM golang:1.24-alpine AS backend
WORKDIR /src
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- 3. Minimal runtime image ----
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=backend /out/server ./server
COPY --from=frontend /app/dist ./public
ENV PORT=8080 \
    STATIC_DIR=/app/public
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
