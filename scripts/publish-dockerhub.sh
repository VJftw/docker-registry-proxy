#!/bin/bash -e

if [ -n "${DOCKERHUB_USERNAME}" ]; then
  docker login -u "${DOCKERHUB_USERNAME}" -p "${DOCKERHUB_PASSWORD}"
fi

version=$(git describe --always)

dockerfiles=$(find ./build -name "*.Dockerfile")

base_repo="vjftw"

for dockerfile in ${dockerfiles}; do
  echo "${dockerfile}"
  repository="${dockerfile//\.\/build\/docker\//}"
  repository="${repository//\.Dockerfile/}"
  echo "-> building Docker image ${base_repo}/${repository}:${version}"
  docker build \
    -f "${dockerfile}" \
    -t "${base_repo}/${repository}:${version}" \
    -t "${base_repo}/${repository}:latest" \
    .
done

for dockerfile in ${dockerfiles}; do
  echo "${dockerfile}"
  repository="${dockerfile//\.\/build\/docker\//}"
  repository="${repository//\.Dockerfile/}"
  echo "-> pushing Docker image ${base_repo}/${repository}:${version}"
  docker push "${base_repo}/${repository}:${version}"
  docker push "${base_repo}/${repository}:latest"
done
