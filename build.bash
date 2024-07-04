#/bin/bash

APP_NAME="module-migration"

docker build -t module-builder -f Dockerfile --build-arg="GITHUB_TOKEN=$GITHUB_TOKEN" --build-arg="APP_NAME=$APP_NAME" .
docker run -d --name module-builder module-builder:latest
docker container cp module-builder:/app/$APP_NAME .
docker rm -f module-builder
