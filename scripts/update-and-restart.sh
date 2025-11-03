#!/bin/bash

# Pull latest images and restart services
# Usage: ./scripts/update-and-restart.sh

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Updating services with latest images"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo "Pulling latest images from DockerHub..."
docker-compose pull

echo ""
echo "Restarting services..."
docker-compose up -d

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✓ Services updated and restarted successfully!"
echo ""
echo "To view logs, run:"
echo "  docker-compose logs -f"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
