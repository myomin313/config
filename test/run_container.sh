#!/bin/bash
set -e

# =========================
# Configuration with defaults
# =========================
APP_NAME="${APP_NAME:-config-service}"
APP_PORT="${APP_PORT:-3012}"
IMAGE_NAME="${IMAGE_NAME:-harbor.kk001.net/kh-game-dev/${APP_NAME}}"
NETWORK_NAME="${NETWORK_NAME:-game-local}"

FULL_IMAGE="${IMAGE_NAME}:latest"

echo "🚀 Deploying container:"
echo "  Image: $FULL_IMAGE"
echo "  Name: $APP_NAME"
echo "  Port: $APP_PORT"
echo "  Network: $NETWORK_NAME"
echo "=============================="

# =========================
# Pull the image from Harbor
# =========================
echo "📥 Pulling image $FULL_IMAGE..."
docker pull "$FULL_IMAGE"

# =========================
# Stop and remove existing container if it exists
# =========================
if docker ps -a --format '{{.Names}}' | grep -q "^${APP_NAME}$"; then
    echo "🛑 Stopping existing container..."
    docker stop ${APP_NAME}
    docker rm ${APP_NAME}
fi

# =========================
# Ensure Docker network exists
# =========================
docker network inspect "${NETWORK_NAME}" >/dev/null 2>&1 || \
docker network create "${NETWORK_NAME}"

# =========================
# Run the container
# =========================
docker run -d \
  --name "${APP_NAME}" \
  --network "${NETWORK_NAME}" \
  -p "${APP_PORT}:3000" \
  --restart unless-stopped \
  "$FULL_IMAGE"

echo "✅ Deployment successful on port ${APP_PORT}"

# =========================
# Show container info
# =========================
echo ""
echo "=============================="
echo "📦 Running containers:"
docker ps --filter "name=${APP_NAME}"

echo ""
echo "🔌 Port mapping:"
docker port ${APP_NAME}

echo ""
echo "📊 Container status:"
docker inspect ${APP_NAME} --format='Status={{.State.Status}}, Running={{.State.Running}}, Started={{.State.StartedAt}}'

echo ""
echo "📜 Container logs (last 50 lines):"
docker logs --tail 50 ${APP_NAME}
echo "=============================="
