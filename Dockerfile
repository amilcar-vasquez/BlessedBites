# -----------------------
# Stage 1: Build Go Binary
# -----------------------
FROM golang:1.24-alpine AS builder

# Install Git (needed by go get sometimes)
RUN apk add --no-cache git

WORKDIR /app

# Copy the entire source early for modules using internal/
COPY . .

# Ensure modules are tidy and cached
RUN go mod tidy && go mod download

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o blessedbites ./cmd/web

# -----------------------
# Stage 2: Slim Runtime Image
# -----------------------
FROM alpine:latest

RUN adduser -D amilcar
WORKDIR /app

COPY --from=builder /app/blessedbites .
COPY ui/ ./ui/
COPY migrations/ ./migrations/
COPY tls/ ./tls/
COPY documentation/ ./documentation/

# 🛠️ Fix permission for tls so non-root user can read it
RUN chown -R amilcar:amilcar /app && chmod -R 755 /app/tls

USER amilcar

EXPOSE 4000

HEALTHCHECK --interval=30s --timeout=5s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:4000/healthz || exit 1

CMD ["./blessedbites"]
