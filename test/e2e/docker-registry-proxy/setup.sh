#!/bin/bash -e

terraform init

terraform apply --auto-approve

image_uri=$(terraform output test_image_uri)
hostname=$(echo "${image_uri}" | cut -f1 -d/)
echo "-> checking for DNS propagation of: ${hostname}"

while [ ! "$(dig +short ${hostname}.)" ]; do
    ts=$(date --rfc-3339=seconds)
    echo "${ts}: waiting for DNS propagation of ${hostname}"
    sleep 10
done
