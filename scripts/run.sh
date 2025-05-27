#!/bin/bash

command=${1:-"start"}

PROJECT_ROOT="$(dirname "$(dirname "$0")")"

export AMBULANCE_API_ENVIRONMENT="Development"
export AMBULANCE_API_PORT="8080"
export AMBULANCE_API_MONGODB_USERNAME="root"
export AMBULANCE_API_MONGODB_PASSWORD="neUhaDnes"

mongo() {
    docker compose --file "${PROJECT_ROOT}/deployments/docker-compose/compose.yaml" "$@"
}

case $command in
    "start")
        mongo up --detach
        trap 'mongo down' EXIT
        cd "${PROJECT_ROOT}" && go run ./cmd/ambulance-api-service/main.go
        ;;
    "mongo")
        mongo up
        ;;
    "test")
        go test -v ./...
        ;;
    "docker")
        docker build -t xprikryl/ambulance-wl-webapi:local-build -f ${PROJECT_ROOT}/build/docker/Dockerfile .
        ;;
    "openapi")
        docker run --rm -ti -v "${PROJECT_ROOT}:/local" openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
        ;;
    *)
        echo "Unknown command: $command"
        exit 1
        ;;
esac 