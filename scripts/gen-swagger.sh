#!/bin/sh
set -eu

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "Generating Swagger docs..."

swag init \
  --dir "$PROJECT_ROOT" \
  --generalInfo "internal/api/router/routes.go" \
  --output "internal/api/docs" \
  --outputTypes go,json,yaml \
  --parseDependency \
  --parseInternal \
  --parseDepth 1