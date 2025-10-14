#!/usr/bin/env bash
set -e 

trap 'docker-compose -f "$COMPOSE_FILE" down -v --remove-orphans' EXIT

APP_NAME="go_server"
COMPOSE_FILE="docker-compose.yml"

echo "Building HUNT server"

docker-compose -f "$COMPOSE_FILE" build

docker-compose -f "$COMPOSE_FILE" up -d

docker ps --filter "name=${APP_NAME}" --filter "name=mongo"

sleep 3

docker-compose -f "$COMPOSE_FILE" logs -f
