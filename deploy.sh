#!/bin/bash
set -e

# =========================
# Load environment variables
# =========================
source .env

# Generate dynamic tag
DYN_TAG=$(date +%Y%m%d%H%M)
FULL_IMAGE="${IMAGE_NAME}:${DYN_TAG}"

echo "🚀 Building $FULL_IMAGE"

# Build using docker compose (latest tag in Compose)
docker compose build --no-cache

# Tag image dynamically
docker tag ${IMAGE_NAME}:${TAG} $FULL_IMAGE

# Push to Harbor
docker push $FULL_IMAGE
echo "✅ Image pushed to Harbor"

# =========================
# Deploy Container
# =========================
if [ -z "$APP_PORT" ]; then
  echo "❌ APP_PORT missing in .env"
  exit 1
fi

echo "🚀 Deploying $FULL_IMAGE on port ${APP_PORT}"

# Stop and remove existing container if it exists
if docker ps -a --format '{{.Names}}' | grep -q "^${APP_NAME}$"; then
    echo "Stopping existing container..."
    docker stop ${APP_NAME}
    docker rm ${APP_NAME}
fi

# Create log directory if it doesn't exist
LOG_DIR="./logs"
mkdir -p ${LOG_DIR}
ERROR_LOG="${LOG_DIR}/${APP_NAME}_error_$(date +%Y%m%d).log"
ACCESS_LOG="${LOG_DIR}/${APP_NAME}_access_$(date +%Y%m%d).log"

# Run new container with logging driver
docker run -d \
  --name "${APP_NAME}" \
  --network "${NETWORK_NAME}" \
  -p "${APP_PORT}:3000" \
  --restart unless-stopped \
  --log-driver json-file \
  --log-opt max-size=10m \
  --log-opt max-file=3 \
  --log-opt tag="${APP_NAME}" \
  -e GIN_MODE=release \
  -e LOG_LEVEL=debug \
  "$FULL_IMAGE"

# Wait for container to start
sleep 5

# =========================
# HEALTH CHECK & LOGGING
# =========================
echo ""
echo "🔍 Performing health checks..."

# Function to check container status
check_container_status() {
    local status=$(docker inspect ${APP_NAME} --format='{{.State.Status}}')
    local health=$(docker inspect ${APP_NAME} --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}No health check{{end}}' 2>/dev/null)
    
    echo "Container Status: ${status}"
    echo "Health Status: ${health}"
    
    if [ "$status" != "running" ]; then
        echo "❌ Container is not running! Status: ${status}"
        capture_error_logs
        return 1
    fi
    return 0
}

# Function to capture error logs
capture_error_logs() {
    echo ""
    echo "⚠️  Capturing error logs..."
    {
        echo "========== ERROR LOGS =========="
        echo "Timestamp: $(date)"
        echo "Container: ${APP_NAME}"
        echo "Image: ${FULL_IMAGE}"
        echo ""
        echo "--- Container Logs (Last 100 lines) ---"
        docker logs --tail 100 ${APP_NAME} 2>&1
        echo ""
        echo "--- Container Inspect ---"
        docker inspect ${APP_NAME}
        echo ""
        echo "--- Container Stats ---"
        docker stats --no-stream ${APP_NAME}
        echo "================================"
    } >> ${ERROR_LOG}
    
    echo "📝 Error logs saved to: ${ERROR_LOG}"
}

# Function to check HTTP endpoint
check_http_endpoint() {
    local max_attempts=30
    local wait_time=2
    
    echo "Checking HTTP endpoint on port ${APP_PORT}..."
    
    for i in $(seq 1 $max_attempts); do
        if curl -s -f -o /dev/null "http://localhost:${APP_PORT}/health" || \
           curl -s -f -o /dev/null "http://localhost:${APP_PORT}/ready" || \
           curl -s -f -o /dev/null "http://localhost:${APP_PORT}/"; then
            echo "✅ HTTP endpoint is responding"
            return 0
        fi
        
        # Check if container is still running
        if ! docker ps -q -f name=${APP_NAME} > /dev/null; then
            echo "❌ Container stopped during startup"
            capture_error_logs
            return 1
        fi
        
        echo "Waiting for HTTP endpoint... (${i}/${max_attempts})"
        sleep $wait_time
    done
    
    echo "❌ HTTP endpoint not responding"
    capture_error_logs
    return 1
}

# Function to monitor container resources
monitor_resources() {
    echo ""
    echo "📊 Container Resource Usage:"
    docker stats --no-stream ${APP_NAME} --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
    
    # Log resource usage
    {
        echo "========== RESOURCE USAGE =========="
        echo "Timestamp: $(date)"
        docker stats --no-stream ${APP_NAME} --no-trunc
        echo "==================================="
    } >> "${LOG_DIR}/${APP_NAME}_resources_$(date +%Y%m%d).log"
}

# Function to stream logs in background (optional)
stream_logs() {
    echo "📜 Streaming logs to ${ACCESS_LOG}"
    docker logs -f ${APP_NAME} >> ${ACCESS_LOG} 2>&1 &
    echo $! > "${LOG_DIR}/${APP_NAME}_log_pid.txt"
}

# Run health checks
if ! check_container_status; then
    echo "❌ Deployment failed - container not healthy"
    exit 1
fi

if ! check_http_endpoint; then
    echo "❌ Deployment failed - HTTP endpoint not responding"
    exit 1
fi

# Monitor resources
monitor_resources

# Optional: Start log streaming
# stream_logs

echo "✅ Deployment successful on port ${APP_PORT}"

# =========================
# DISPLAY INFORMATION
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
docker logs --tail 50 ${APP_NAME} 2>&1

echo ""
echo "📍 Log files location:"
echo "   Error logs: ${ERROR_LOG}"
echo "   Access logs: ${ACCESS_LOG}"
echo "   Resource logs: ${LOG_DIR}/${APP_NAME}_resources_*.log"

echo "=============================="