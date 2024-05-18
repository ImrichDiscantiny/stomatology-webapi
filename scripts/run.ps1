param (
    $command
)

if (-not $command)  {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:STOMATOLOGY_API_ENVIRONMENT="Development"
$env:STOMATOLOGY_API_PORT="8080"
$env:AMBULANCE_API_MONGODB_USERNAME="root"
$env:AMBULANCE_API_MONGODB_PASSWORD="neUhaDnes"

function mongo {
    docker compose --file ${ProjectRoot}/deployments/docker-compose/compose.yaml $args
}

switch ($command) {
    "openapi" {
        docker run --rm -ti -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    "start" {
        try {
            mongo up --detach
            go run ${ProjectRoot}/cmd/stomatology-api-service
        } finally {
            mongo down
        }
    }
    "docker" {
        docker build -t imrdsc/stomatology-al-webapi:local-build -f ${ProjectRoot}/build/docker/Dockerfile .
    }
    "mongo" {
        mongo up
        }
    default {
        throw "Unknown command: $command"
    }
}