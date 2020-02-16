#!/bin/bash

### SETTINGS
gcp_billing_account_id="00C32F-55AB5C-78A691"
domain="docker.vjpatel.me"

gcloud auth configure-docker --quiet

cat <<EOF > terraform/terraform.tfvars
billing_account_id = "${gcp_billing_account_id}"
EOF

cd terraform
terraform init
terraform apply --auto-approve || terraform apply --auto-approve
cd ../

project=$(jq -r '.resources[].instances[].attributes.project' terraform/terraform.tfstate | grep -v null | head -n1)
# gcloud container clusters get-credentials kubelet-image-service --region europe-west1 --project "${project}"


# Deploy Kubelet Image Service
cat <<EOF
Done!

configure kubectl to use the GKE cluster:
$ gcloud container clusters get-credentials kubelet-image-service --region europe-west1 --project "${project}"

Deploy the kubelet-image-service by running:
$ kustomize build k8s/gke | kubectl apply -f-

And the ghost cms using an image from the docker registry proxy
$ kustomize build ghost-k8s/ | kubectl apply -f-

configure kubectl to use the EKS cluster:
$ aws eks update-kubeconfig --name test

Deploy the kubelet-image-service by running:
$ kustomize build k8s/eks | kubectl apply -f-

And the ghost cms using an image from the docker registry proxy
$ kustomize build ghost-k8s/ | kubectl apply -f-
EOF
