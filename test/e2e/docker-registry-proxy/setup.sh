#!/bin/bash -e

var_file="${1}"

if [ -z "${var_file}" ]; then 
    var_file="terraform.tfvars"
fi
echo "-> using ${var_file}"

terraform init

terraform apply --auto-approve -var-file="${var_file}"

image_uri=$(terraform output test_image_uri)
hostname=$(echo "${image_uri}" | cut -f1 -d/)
echo "-> checking for DNS propagation of: ${hostname}"

while [ ! "$(dig +short ${hostname}.)" ]; do
    ts=$(date --rfc-3339=seconds)
    echo "${ts}: waiting for DNS propagation of ${hostname}"
    sleep 10
done
