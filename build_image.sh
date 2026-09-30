#!/usr/bin/env bash
set -e

# ==========================================
# Global Error Handler (Build Server)
# ==========================================
error_handler() {
  echo ""
  echo "❌ ERROR: Image build/push failed at line $1"

  #echo "====== Docker Info ======"
  #docker --version || true
  #docker info || true

  #echo "====== Local Images ======"
  #docker images || true

  exit 1
}

trap 'error_handler $LINENO' ERR

# ==========================================
# System / Docker Info
# ==========================================
echo "====== System Info ======"
date
hostname
docker --version
echo "========================="

echo "🔍 Checking Docker daemon..."
docker info >/dev/null 2>&1 || {
  echo "❌ Docker daemon not running"
  exit 1
}

# ==========================================
# Load configuration
# ==========================================
[ -f .env ] || { echo ".env file not found"; exit 1; }
source .env

# Validate required variables
: "${IMAGE_NAME:?IMAGE_NAME not set in .env}"
: "${TAG:?TAG not set in .env}"
: "${HARBOR_REGISTRY:?HARBOR_REGISTRY not set in .env}"
: "${HARBOR_PROJECT:?HARBOR_PROJECT not set in .env}"
: "${APP_NAME:?APP_NAME not set in .env}"

# ==========================================
# Get current branch
# ==========================================
BRANCH=${BRANCH}
ENVIRONMENT=${BRANCH}


# ==========================================
# IMAGE_NAME already has full path from .env
# ==========================================
# .env has: IMAGE_NAME="${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${APP_NAME}"
# So FULL_IMAGE is just IMAGE_NAME:TAG
FULL_IMAGE="${IMAGE_NAME}:${TAG}"

echo "========================================="
echo "🚀 Building Docker Image"
echo "Branch: $BRANCH"
echo "Environment: $ENVIRONMENT"
echo "Harbor Project: $HARBOR_PROJECT"
echo "App Name: $APP_NAME"
echo "Image: $FULL_IMAGE"
echo "Tag: $TAG"
echo "========================================="


# ==========================================
# Build Image
# ==========================================
docker build \
  --build-arg HARBOR_REGISTRY=${HARBOR_REGISTRY} \
  --build-arg HARBOR_PROJECT=${HARBOR_PROJECT} \
  --build-arg ENVIRONMENT=${ENVIRONMENT} \
  --build-arg APP_NAME=${APP_NAME} \
  -t "$FULL_IMAGE" .

echo "🔍 Verifying image built..."
docker image inspect "$FULL_IMAGE" >/dev/null 2>&1 || {
  echo "❌ Image build failed"
  exit 1
}

echo "📦 Built image details:"
docker images | grep "$APP_NAME" || true

# ==========================================
# Push Image
# ==========================================
echo "📤 Pushing image with tag: $TAG"
docker push "$FULL_IMAGE"

echo "🔍 Verifying push by pulling image..."
docker pull "$FULL_IMAGE" >/dev/null 2>&1 || {
  echo "❌ Image push verification failed"
  exit 1
}

echo "========================================="
echo "✅ Image pushed successfully!"
echo "========================================="

echo "🧹 Cleaning local Docker images..."
docker rmi -f "$FULL_IMAGE" || true
#docker image prune -f || true


# ==========================================
# Update run_container.sh with latest values
# ==========================================
if [ -f "run_container.sh" ]; then
  echo "📝 Updating run_container.sh configuration..."
  
  sed -i "s|^APP_NAME=.*|APP_NAME=\"$APP_NAME\"|" run_container.sh
  sed -i "s|^APP_PORT=.*|APP_PORT=\"$APP_PORT\"|" run_container.sh
  sed -i "s|^NETWORK_NAME=.*|NETWORK_NAME=\"$NETWORK_NAME\"|" run_container.sh
  sed -i "s|^FULL_IMAGE=.*|FULL_IMAGE=\"${FULL_IMAGE}\"|" run_container.sh
  sed -i "s|^NACOS_USER=.*|NACOS_USER=\"$NACOS_USER\"|" run_container.sh
  sed -i "s|^NACOS_PASSWORD=.*|NACOS_PASSWORD=\"$NACOS_PASSWORD\"|" run_container.sh
  sed -i "s|^NACOS_HOST=.*|NACOS_HOST=\"$NACOS_HOST\"|" run_container.sh
  sed -i "s|^NACOS_PORT=.*|NACOS_PORT=\"$NACOS_PORT\"|" run_container.sh
  sed -i "s|^NACOS_NAMESPACE_ID=.*|NACOS_NAMESPACE_ID=\"$NACOS_NAMESPACE_ID\"|" run_container.sh
  sed -i "s|^GROUP_NAME=.*|GROUP_NAME=\"$NACOS_GROUP\"|" run_container.sh
  sed -i "s|^SERVICE_IP=.*|SERVICE_IP=\"$SERVICE_IP\"|" run_container.sh
  sed -i "s|^APP_KEY=.*|APP_KEY=\"$APP_KEY\"|" run_container.sh
  sed -i "s|^TZ=.*|TZ=\"$TZ\"|" run_container.sh
  sed -i "s|^DB_USER=.*|DB_USER=\"$DB_USER\"|" run_container.sh
  sed -i "s|^DB_PASS=.*|DB_PASS=\"$DB_PASS\"|" run_container.sh
  sed -i "s|^DB_HOST=.*|DB_HOST=\"$DB_HOST\"|" run_container.sh
  sed -i "s|^DB_NAME=.*|DB_NAME=\"$DB_NAME\"|" run_container.sh
  sed -i "s|^DB_PORT=.*|DB_PORT=\"$DB_PORT\"|" run_container.sh
  sed -i "s|^LOG_DIR=.*|LOG_DIR=\"$LOG_DIR\"|" run_container.sh

  echo "✅ run_container.sh updated"
fi

echo "========================================="
echo "🎉 All done! Image built and GitOps updated"
echo "========================================="