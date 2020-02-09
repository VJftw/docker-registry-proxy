#!/bin/bash -e

docker build \
  -f build/docker/docker-registry-proxy.Dockerfile \
  -t vjftw/docker-registry-proxy:latest \
  .

docker build \
  -f build/docker/kubelet-image-service.Dockerfile \
  -t vjftw/kubelet-image-service:latest \
  .

docker push vjftw/docker-registry-proxy:latest
docker push vjftw/kubelet-image-service:latest
