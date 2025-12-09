# !/bin/sh
set -eu

src_dir=$(dirname "$0")/..

echo "Generating Swagger docs.."

swag init \
    --generalInfo "../meals-service/http/routes/router.go" \
    --output "./docs" \
    --outputTypes go,json,yaml \
    --parseDependency --parseInternal --parseDepth 1