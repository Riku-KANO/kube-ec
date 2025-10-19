#!/bin/bash
set -e

# Script to build Docker image for web frontend
# Usage: docker_build.sh <path-to-dockerfile> <image-tag>

DOCKERFILE=${1:-"web/Dockerfile"}
IMAGE_TAG=${2:-"kube-ec-web:latest"}

echo "Building Docker image: $IMAGE_TAG"
echo "Using Dockerfile: $DOCKERFILE"

# Determine which container runtime to use
if command -v nerdctl &> /dev/null; then
    RUNTIME="nerdctl"
    echo "Using nerdctl (Rancher Desktop)"
elif command -v docker &> /dev/null; then
    RUNTIME="docker"
    echo "Using docker"
else
    echo "Error: Neither docker nor nerdctl found"
    exit 1
fi

# Get the directory containing the Dockerfile
DOCKER_DIR=$(dirname "$DOCKERFILE")

# Build the image
cd "$DOCKER_DIR"
$RUNTIME build -t "$IMAGE_TAG" .

echo "Successfully built image: $IMAGE_TAG"

# Show the image
$RUNTIME images | grep kube-ec-web || true
