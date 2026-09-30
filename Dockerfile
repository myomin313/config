# ======================
# Build stage
# ======================
ARG HARBOR_REGISTRY=35.220.207.63
FROM ${HARBOR_REGISTRY}/base_images_php/golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go

RUN go build -o server .

# ======================
# Run stage
# ======================
ARG HARBOR_REGISTRY
FROM ${HARBOR_REGISTRY}/base_images_php/nginx:1.27.2-alpine

WORKDIR /app

# copy env into image
COPY .env /app/.env

# create user
# RUN getent group www-data || addgroup -g 1000 www-data \
 # && getent passwd www-data || adduser -D -u 1000 -G www-data www-data

RUN mkdir -p /app/logs
#RUN mkdir -p /app/logs/sit-kla-klouk/logs/golang

COPY --from=builder /app/server .
COPY --from=builder /app/docs ./docs

#RUN chown -R www-data:www-data /app

EXPOSE 3000
# -----------------------------
# Healthcheck for Go service
# -----------------------------
# Checks if the Go server responds on port 3000
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -q --spider http://localhost:3000/health || exit 1

# USER www-data

CMD ["./server"]