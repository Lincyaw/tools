# Docker Images

This directory contains scripts for building and pushing Docker images to DockerHub.

## Prerequisites

- Docker installed and running
- Docker Hub account
- Logged in to Docker Hub: `docker login`

## Scripts

### build-images.sh

Builds all Docker images locally without pushing to DockerHub.

**Usage:**
```bash
# Use default username (lincyaw)
./scripts/build-images.sh

# Use custom username
./scripts/build-images.sh yourusername

# Build with custom version
VERSION=v1.0.0 ./scripts/build-images.sh
```

### build-and-push.sh

Builds all Docker images and pushes them to DockerHub.

**Usage:**
```bash
# Use default username (lincyaw)
./scripts/build-and-push.sh

# Use custom username
./scripts/build-and-push.sh yourusername

# Build and push with custom version
VERSION=v1.0.0 ./scripts/build-and-push.sh
```

## Makefile Commands

The Makefile provides convenient shortcuts for common Docker operations:

```bash
# Build images locally
make docker-build

# Build images with custom username
make docker-build DOCKER_USERNAME=yourusername

# Build and push images to DockerHub
make docker-push

# Build and push with custom username
make docker-push DOCKER_USERNAME=yourusername

# Start services using DockerHub images
make docker-up

# Start services building locally (development)
make docker-dev-up

# Stop services
make docker-down

# View logs
make docker-logs
```

## Docker Compose Files

- **docker-compose.yml**: Production configuration, uses pre-built images from DockerHub
- **docker-compose.dev.yml**: Development configuration, builds images locally

## Images

The following images are built and pushed:

- `lincyaw/tools-gateway:latest` - Caddy gateway service
- `lincyaw/tools-shortcode:latest` - Shortcode API service

## Workflow

### For Development

```bash
# Build and run services locally
make docker-dev-up

# View logs
make docker-logs

# Stop services
make docker-dev-down
```

### For Production Deployment

```bash
# 1. Build and push images to DockerHub
make docker-push

# 2. On production server, pull and run
make docker-up
```

### Updating Images

```bash
# 1. Make your code changes
# 2. Build and push new images
make docker-push

# 3. On production server, pull latest images and restart
docker-compose pull
docker-compose up -d
```

## Notes

- Default Docker username is `lincyaw`. Override with `DOCKER_USERNAME` environment variable.
- Default version tag is `latest`. Override with `VERSION` environment variable.
- Make sure you're logged in to Docker Hub before pushing: `docker login`
