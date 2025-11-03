#!/bin/bash

# Build Docker images locally
# Usage: ./scripts/build-images.sh [DOCKER_USERNAME]

set -e

# Configuration
DOCKER_USERNAME="${1:-lincyaw}"
VERSION="${VERSION:-latest}"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Building Docker images locally"
echo "Docker username: ${DOCKER_USERNAME}"
echo "Version: ${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Build gateway image
echo ""
echo "Building gateway image..."
docker build -t ${DOCKER_USERNAME}/tools-gateway:${VERSION} \
    -f services/gateway/Dockerfile \
    services/gateway

# Build shortcode service image
echo ""
echo "Building shortcode service image..."
docker build -t ${DOCKER_USERNAME}/tools-shortcode:${VERSION} \
    -f services/shortcode/Dockerfile \
    services/shortcode

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✓ All images built successfully!"
echo ""
echo "Images built:"
echo "  - ${DOCKER_USERNAME}/tools-gateway:${VERSION}"
echo "  - ${DOCKER_USERNAME}/tools-shortcode:${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
