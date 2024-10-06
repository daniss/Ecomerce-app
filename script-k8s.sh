#!/bin/bash

# Set the image name and tag
IMAGE_NAME="product-service user-service order-service"
TAG="latest"
DOCKER_USERNAME="your_docker_username"
K8S_DEPLOYMENTS="product-service user-service order-service"

if [ -n "$1" ]; then
    IMAGE_NAME="$1"
    K8S_DEPLOYMENTS="$1"
fi

# Build the Docker image
echo "Building Docker image..."
for service in ${IMAGE_NAME}; do
    sudo docker build -t ${service}:${TAG} ${service}/.
    if [ $? -ne 0 ]; then
        echo "Docker build for ${service} failed. Exiting."
        exit 1
    fi
done

# Tag the Docker image
echo "Tagging Docker image..."
for service in ${IMAGE_NAME}; do
    sudo docker tag ${service}:${TAG} ${DOCKER_USERNAME}/${service}:${TAG}
done

# Push the Docker image to Docker Hub
echo "Pushing Docker image to Docker Hub..."
for service in ${IMAGE_NAME}; do
    sudo docker push ${DOCKER_USERNAME}/${service}:${TAG}
    if [ $? -ne 0 ]; then
        echo "Docker push for ${service} failed. Exiting."
        exit 1
    fi
    echo "Docker image for ${service} pushed successfully!"
done

# Restart the Kubernetes deployment
echo "Restarting Kubernetes deployment..."
for deployment in ${K8S_DEPLOYMENTS}; do
    sudo kubectl rollout restart deployment/${deployment}
done

echo "Kubernetes deployment restarted successfully!"
