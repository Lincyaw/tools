#!/bin/bash

# Build and push Docker images to DockerHub
# Usage: ./scripts/build-and-push.sh [DOCKER_USERNAME]

set -e

# Configuration
DOCKER_USERNAME="${1:-lincyaw}"
VERSION="${VERSION:-latest}"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Building and pushing Docker images"
echo "Docker username: ${DOCKER_USERNAME}"
echo "Version: ${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"


# Build and push gateway image
echo ""
echo "Building gateway image..."
docker build -t ${DOCKER_USERNAME}/tools-gateway:${VERSION} \
    -f services/gateway/Dockerfile \
    services/gateway

echo "Pushing gateway image..."
docker push ${DOCKER_USERNAME}/tools-gateway:${VERSION}

# Build and push shortcode service image
echo ""
echo "Building shortcode service image..."
docker build -t ${DOCKER_USERNAME}/tools-shortcode:${VERSION} \
    -f services/shortcode/Dockerfile \
    services/shortcode

echo "Pushing shortcode service image..."
docker push ${DOCKER_USERNAME}/tools-shortcode:${VERSION}

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✓ All images built and pushed successfully!"
echo ""
echo "Images pushed:"
echo "  - ${DOCKER_USERNAME}/tools-gateway:${VERSION}"
echo "  - ${DOCKER_USERNAME}/tools-shortcode:${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
