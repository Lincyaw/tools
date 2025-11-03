# Tools

A collection of microservices including URL shortcode service and more.

## Services

- **Gateway**: Caddy-based reverse proxy and static file server
- **Shortcode**: URL shortening service with analytics
- **PostgreSQL**: Database for persistent storage
- **Redis**: Cache and session storage

## Quick Start

### Using Pre-built Images from DockerHub

```bash
# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Local Development

```bash
# Build and run services locally
make docker-dev-up

# View logs
make docker-logs

# Stop services
make docker-dev-down
```

## Docker Image Management

### Building Images

```bash
# Build images locally
make docker-build

# Build with custom username
make docker-build DOCKER_USERNAME=yourusername
```

### Pushing to DockerHub

```bash
# Build and push images
make docker-push

# Build and push with custom username
make docker-push DOCKER_USERNAME=yourusername

# Build and push with version tag
VERSION=v1.0.0 make docker-push
```

For more details, see [scripts/README.md](scripts/README.md).

## Development

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- golangci-lint (optional, will be installed automatically)

### Running Tests

```bash
# Run all tests
make test

# Run client tests only
make client-test

# Run server tests only
make server-test
```

### Code Quality

```bash
# Format code
make fmt

# Run linters
make lint

# Run all quality checks (format check, vet, lint, test)
make all
```

### Available Make Commands

```bash
make install-tools    # Install development tools
make lint            # Run linters on all modules
make fmt             # Format all Go code
make fmt-check       # Check if code is formatted
make vet             # Run go vet
make test            # Run all tests
make all             # Run all quality checks

# Docker commands
make docker-up       # Start services (DockerHub images)
make docker-down     # Stop services
make docker-logs     # View logs
make docker-build    # Build images locally
make docker-push     # Build and push to DockerHub
make docker-dev-up   # Start services (local build)
make docker-dev-down # Stop development services
```

## Project Structure

```
.
├── client/                 # CLI client for services
│   ├── cmd/               # Command implementations
│   └── pkg/               # Client packages
├── services/
│   ├── gateway/           # Caddy gateway
│   │   ├── Caddyfile     # Caddy configuration
│   │   └── html/         # Static files
│   └── shortcode/        # URL shortcode service
│       ├── cmd/          # Service entry point
│       └── internal/     # Internal packages
├── scripts/              # Build and deployment scripts
├── docker-compose.yml    # Production compose file (uses DockerHub)
└── docker-compose.dev.yml # Development compose file (local build)
```

## Configuration

Service configurations are managed through environment variables in `docker-compose.yml`:

- `APP_ENV`: Application environment (production/development)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`: PostgreSQL connection
- `REDIS_HOST`, `REDIS_PORT`: Redis connection
- `BASE_URL`: Base URL for shortcode service

## License

MIT
