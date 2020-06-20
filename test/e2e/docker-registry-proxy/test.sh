#!/bin/bash -e

echo "-- Performing Test --"
# 1. get test image uri from terraform outputs
image_uri=$(terraform output test_image_uri)
echo "-> Attempting to pull Docker Image ${image_uri}"
# 2. pull test image via uri
docker pull "${image_uri}"
# 3
echo "-- Success!"
