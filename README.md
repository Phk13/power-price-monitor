# Power Price Monitor

A Go application for monitoring Spanish electricity prices using the ESIOS REE API. Find optimal hours for energy consumption based on price thresholds and get real-time pricing data.

## Features

- 🔍 **Optimal Hours Detection**: Find the cheapest hours to consume electricity
- 💰 **Price Thresholds**: Filter hours based on your preferred price limits
- 🕐 **Real-time Data**: Get current electricity prices for Spain
- 📊 **Smart Caching**: Efficient data caching to minimize API calls
- 🌍 **Timezone Support**: Proper timezone handling for accurate time calculations
- 🐳 **Docker Ready**: Optimized Docker setup for easy deployment
- 🏥 **Health Checks**: Built-in health monitoring

## API Endpoints

### GET /optimal-hours

Find optimal hours for electricity consumption based on price and time constraints.

**Query Parameters:**
- `max_hours` (required): Maximum number of hours to return (1-24)
- `threshold` (required): Maximum price threshold in €/kWh

**Example:**
```bash
curl "http://localhost:8080/optimal-hours?max_hours=8&threshold=0.15"
```

**Response:**
```json
{
  "optimal_hours": ["02:00", "03:00", "04:00", "05:00", "14:00", "15:00", "16:00", "17:00"],
  "total_hours_selected": 8,
  "current_price": 0.123,
  "threshold_used": 0.15,
  "max_hours_used": 8,
  "next_start": "2024-01-15T14:00:00+01:00"
}
```

### GET /health

Health check endpoint for monitoring service status.

**Response:**
```json
{
  "status": "healthy"
}
```

## Configuration

The application uses environment variables for configuration. You can set these in a `.env` file or as system environment variables.

### Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit the `.env` file and set your configuration values:
   ```bash
   TOKEN=your_actual_api_token
   TIMEZONE=Europe/Madrid
   PORT=8080
   ```

### Environment Variables

- `TOKEN` (required): API token for accessing ESIOS REE API data
- `TIMEZONE` (optional): Timezone for calculations (default: Europe/Madrid)
- `PORT` (optional): Server port (default: 8080)

## Development

### Prerequisites

- Go 1.21 or later
- Valid ESIOS REE API token

### Running Locally

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the application:
   ```bash
   make run
   # or
   go run cmd/main.go
   ```

3. Test the API:
   ```bash
   curl "http://localhost:8080/optimal-hours?max_hours=5&threshold=0.12"
   ```

### Available Make Commands

```bash
make help                 # Show available commands
make run                  # Run the application
make lint                 # Run linters
make docker-build         # Build Docker image
make docker-run           # Run Docker container
make docker-compose-up    # Start with docker-compose
make docker-compose-down  # Stop docker-compose services
make docker-compose-logs  # Show docker-compose logs
make docker-clean         # Clean Docker resources
```

## Docker Deployment

### Using Docker Compose (Recommended)

1. Configure environment:
   ```bash
   cp .env.example .env
   # Edit .env with your values
   ```

2. Start the service:
   ```bash
   make docker-compose-up
   # or
   docker-compose up -d
   ```

3. Check logs:
   ```bash
   make docker-compose-logs
   ```

### Using Docker directly

1. Build the image:
   ```bash
   make docker-build
   ```

2. Run the container:
   ```bash
   make docker-run
   # or
   docker run --rm -p 8080:8080 --env-file .env power-price-monitor
   ```

## Production Deployment

The Docker image is optimized for production with:

- ✅ Multi-stage build for minimal image size (~15MB)
- ✅ Non-root user for security
- ✅ Health checks included
- ✅ Proper timezone and CA certificates
- ✅ Resource limits support

### Example production docker-compose.yml

```yaml
version: '3.8'
services:
  power-price-monitor:
    image: power-price-monitor:latest
    ports:
      - "8080:8080"
    environment:
      - TOKEN=${TOKEN}
      - TIMEZONE=Europe/Madrid
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 128M
          cpus: '0.5'
```

## Architecture

The application follows a clean architecture pattern:

- **Handler Layer**: HTTP request handling and validation
- **Service Layer**: Business logic and data processing
- **Model Layer**: Data structures and types
- **Config Layer**: Configuration management

Key features:
- Context-aware request handling for proper cancellation
- Smart caching to minimize API calls
- Timezone-aware calculations
- Comprehensive error handling

## License

This project is licensed under the MIT License.
