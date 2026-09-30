#!/bin/bash
set -e

# ==========================================
# Hardcoded Configuration
# ==========================================
APP_NAME="config-service"
APP_PORT="3012"
NETWORK_NAME="game-local"
FULL_IMAGE="35.220.207.63/kh-game-dev/config-service:20260306133922"
NACOS_HOST="192.168.102.102"
NACOS_PORT="8848"
NACOS_USER="php"
NACOS_PASSWORD="cLp4nhXoitUoPUhVNLPxt2YKDvrUb4FQ"
NACOS_NAMESPACE_ID="kkr-dev-env2"
APP_URL="https://dev-klaklouk-api-gateway.lomasq.com:9025"
GROUP_NAME="DEV_GROUP"
SERVICE_IP="192.168.102.108"
APP_KEY="6pjO7kIbbdOWqTHOr3u9OI8m-EzgUNJ_IaIJ4pr02iA"
TZ="Asia/Phnom_Penh"
DB_USER="postgres"
DB_PASS="Amma@abba@123"
DB_HOST="10.10.13.115"
DB_NAME="admin-service"
DB_PORT="5432"
LOG_DIR="/home/deploy/sit-kla-klouk/logs"

# ==========================================
# Global Error Handler (Jenkins logging)
# ==========================================
error_handler() {
  echo ""
  echo "❌ ERROR: Deployment failed at line $1"
  echo "========================================="

  echo "📊 docker ps:"
  docker ps -a || true

  echo ""
  echo "📜 Container Logs (if exists):"
  docker logs "$APP_NAME" 2>/dev/null || true

  echo "========================================="
  exit 1
}

trap 'error_handler $LINENO' ERR

# ==========================================
# System / Docker Info
# ==========================================
#echo "====== System Info ======"
#date
#hostname
#docker --version
#echo "========================="

echo "🔍 Checking Docker daemon..."
docker info >/dev/null 2>&1 || {
  echo "❌ Docker daemon not running"
  exit 1
}

# ==========================================
# Deployment Info
# ==========================================
echo "========================================="
echo "🚀 Deploying container"
echo "App: $APP_NAME"
echo "Image: $FULL_IMAGE"
echo "Port: $APP_PORT"
echo "Network: $NETWORK_NAME"
echo "========================================="

# ==========================================
# Pull Image
# ==========================================
echo "📥 Pulling image..."
docker pull "$FULL_IMAGE"

echo "🔍 Verifying image pulled..."
docker image inspect "$FULL_IMAGE" >/dev/null 2>&1 || {
  echo "❌ Image pull failed"
  exit 1
}

# ==========================================
# Stop existing container if exists
# ==========================================
if docker ps -a --format '{{.Names}}' | grep -q "^${APP_NAME}$"; then
    echo "Stopping existing container..."
    docker stop -t 10 "$APP_NAME" || true
    docker rm "$APP_NAME" || true
fi

# ==========================================
# Create network if not exists
# ==========================================
docker network inspect "$NETWORK_NAME" >/dev/null 2>&1 || docker network create "$NETWORK_NAME"

# ==========================================
# Run container 
# ==========================================

FINAL_LOG_DIR="$LOG_DIR/$APP_NAME"


mkdir -p "$FINAL_LOG_DIR"

# Give Promtail read access to logs
#chmod -R 755 "$LOG_DIR"
#chmod -R 777 "$FINAL_LOG_DIR"   # safe for dev environment
echo "🚀 Starting container..."


docker run -d \
  --name "$APP_NAME" \
  --network "$NETWORK_NAME" \
  -p "$APP_PORT:3000" \
  -v "$FINAL_LOG_DIR:/app/logs" \
  --memory="512m" \
  --restart unless-stopped \
  "$FULL_IMAGE"

# ==========================================
# Verify Container Deployment (Jenkins logs)
# ==========================================
# echo "📊 Showing platforms and platform_api_keys tables data..."

# docker run --rm \
#   -e PGPASSWORD="$DB_PASS" \
#   postgres:15 \
#   psql -h "$DB_HOST" \
#        -p "$DB_PORT" \
#        -U "$DB_USER" \
#        -d "$DB_NAME" \
#        -c "SELECT * FROM platforms; SELECT * FROM platform_api_keys;"

# echo "✅ Data displayed successfully"

echo "========================================="
echo "🔍 Checking container status..."
echo "========================================="

# Wait for container startup
sleep 12

echo ""
echo "📊 docker ps (running containers):"
docker ps --filter "name=$APP_NAME" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" || true

echo ""
echo "📊 docker ps -a (all containers):"
docker ps -a --filter "name=$APP_NAME" --format "table {{.Names}}\t{{.Status}}\t{{.Image}}" || true

# Get container state
STATUS=$(docker inspect --format='{{.State.Status}}' "$APP_NAME" 2>/dev/null || echo "not_found")

echo ""
echo "Container Status: $STATUS"

# If container not running → show logs + fail Jenkins
if [ "$STATUS" != "running" ]; then
    echo "❌ Container is NOT running!"
    echo "====== Container Logs ======"
    docker logs "$APP_NAME" 2>/dev/null || true
    echo "============================"
    exit 1
fi

echo "✅ Container is running successfully"

# Show restart count
echo ""
echo "🔁 Restart Count:"
docker inspect --format='{{.RestartCount}}' "$APP_NAME" 2>/dev/null || true

# Show container details
echo ""
echo "📦 Container Details:"
docker inspect --format='
Status: {{.State.Status}}
Started: {{.State.StartedAt}}
Health: {{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}
IP: {{.NetworkSettings.IPAddress}}
' "$APP_NAME" 2>/dev/null || true

# Show last logs even if success
echo ""
echo "📜 Last 100 container logs:"
#docker logs --tail 100 "$APP_NAME" 2>/dev/null || true
#docker logs "$APP_NAME" 2>&1 | tail -200
docker logs "$APP_NAME" 2>&1

echo "========================================="
echo "✅ Deployment verification completed"
echo "========================================="

echo ""
# echo "🌐 Testing service locally..."

# curl -X POST http://localhost:8002/api/banners \
# -H "Content-Type: application/json" \
# -d '{
#   "display_duration": 10,
#   "display_platform": "h5",
#   "start_date": 1710000000,
#   "end_date": 1710600000,
#   "image_name": "banner.png"
# }'

# ==========================================
# List all files and folders in LOG_DIR
# ==========================================
echo ""
#LOG_FOLDER="../../logs/"
echo "📂 Listing all files and folders recursively in $LOG_DIR:"

if [ -d "$LOG_DIR" ]; then
    ls -alhR "$LOG_DIR"
else
    echo "⚠️ Directory $LOG_DIR does not exist!"
fi
