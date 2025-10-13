#!/usr/bin/env bash
set -e  # Exit immediately on error

trap 'echo "🛑 Cleaning up..."; docker-compose -f "$COMPOSE_FILE" down -v --remove-orphans' EXIT

APP_NAME="go_server"
COMPOSE_FILE="docker-compose.yml"

echo "🚀 Building Go app and Docker images..."

# Build Docker images
docker-compose -f "$COMPOSE_FILE" build

echo "✅ Build complete."

# Start the containers
echo "🟢 Starting $APP_NAME and MongoDB..."
docker-compose -f "$COMPOSE_FILE" up -d

# Show running containers
echo "📦 Running containers:"
docker ps --filter "name=${APP_NAME}" --filter "name=mongo"

# Wait a few seconds for MongoDB to initialize
sleep 3

# Tail the logs so user can see output
echo "📜 Streaming logs (Ctrl+C to stop):"
docker-compose -f "$COMPOSE_FILE" logs -f
