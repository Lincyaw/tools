#!/bin/bash

# Build and push Docker images to DockerHub
# Usage: ./scripts/build-and-push.sh [DOCKER_USERNAME]

set -e

# Configuration
DOCKER_USERNAME="${1:-lincyaw}"
VERSION="${VERSION:-latest}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Building and pushing Docker images"
echo "Docker username: ${DOCKER_USERNAME}"
echo "Version: ${VERSION}"
echo "Platforms: ${PLATFORMS}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if logged in to Docker Hub
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running"
    exit 1
fi

# Ensure buildx is available and create/use a builder
echo ""
echo "Setting up Docker buildx..."
docker buildx create --name multiarch-builder --use 2>/dev/null || docker buildx use multiarch-builder 2>/dev/null || docker buildx use default

# Build and push gateway image
echo ""
echo "Building and pushing gateway image for multiple platforms..."
docker buildx build --platform ${PLATFORMS} \
    -t ${DOCKER_USERNAME}/tools-gateway:${VERSION} \
    -f services/gateway/Dockerfile \
    --push \
    services/gateway

# Build and push shortcode service image
echo ""
echo "Building and pushing shortcode service image for multiple platforms..."
docker buildx build --platform ${PLATFORMS} \
    -t ${DOCKER_USERNAME}/tools-shortcode:${VERSION} \
    -f services/shortcode/Dockerfile \
    --push \
    services/shortcode

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✓ All images built and pushed successfully!"
echo ""
echo "Images pushed:"
echo "  - ${DOCKER_USERNAME}/tools-gateway:${VERSION}"
echo "  - ${DOCKER_USERNAME}/tools-shortcode:${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
