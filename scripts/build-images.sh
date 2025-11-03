#!/bin/bash

# Build Docker images locally
# Usage: ./scripts/build-images.sh [DOCKER_USERNAME]

set -e

# Configuration
DOCKER_USERNAME="${1:-lincyaw}"
VERSION="${VERSION:-latest}"
# For local builds, only build for current platform unless specified
PLATFORMS="${PLATFORMS:-}"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Building Docker images locally"
echo "Docker username: ${DOCKER_USERNAME}"
echo "Version: ${VERSION}"
if [ -n "${PLATFORMS}" ]; then
    echo "Platforms: ${PLATFORMS}"
else
    echo "Platform: current architecture (for local use)"
fi
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "${PLATFORMS}" ]; then
    # Multi-platform build (won't be loadable to local Docker)
    echo ""
    echo "Setting up Docker buildx..."
    docker buildx create --name multiarch-builder --use 2>/dev/null || docker buildx use multiarch-builder 2>/dev/null || docker buildx use default

    # Build gateway image
    echo ""
    echo "Building gateway image for multiple platforms..."
    docker buildx build --platform ${PLATFORMS} \
        -t ${DOCKER_USERNAME}/tools-gateway:${VERSION} \
        -f services/gateway/Dockerfile \
        services/gateway

    # Build shortcode service image
    echo ""
    echo "Building shortcode service image for multiple platforms..."
    docker buildx build --platform ${PLATFORMS} \
        -t ${DOCKER_USERNAME}/tools-shortcode:${VERSION} \
        -f services/shortcode/Dockerfile \
        services/shortcode
else
    # Single platform build for local use
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
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✓ All images built successfully!"
echo ""
echo "Images built:"
echo "  - ${DOCKER_USERNAME}/tools-gateway:${VERSION}"
echo "  - ${DOCKER_USERNAME}/tools-shortcode:${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
