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
gcloud container clusters get-credentials kubelet-image-service --region europe-west1 --project "${project}"

cat <<EOF > k8s/gcp-provider.yaml
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kubelet-image-service
  labels:
    k8s-app: kubelet-image-service
spec:
  selector:
    matchLabels:
      name: kubelet-image-service
  template:
    metadata:
      labels:
        name: kubelet-image-service
spec:
  template:
    spec:
      containers:
        - name: kubelet-image-service
          env:
            - name: DOCKER_API_VERSION
              value: "1.39"
            - name: KIS_PLUGINS
              value: "auth-provider_gcp-instanceidentitydocument:gcpidd"
            - name: KIS_GCPIDD_USERNAME
              value: "_gcpidd"
            - name: KIS_GCPIDD_AUDIENCE
              value: "${domain}"
            - name: KIS_AUTHENTICATION_PROVIDER
              value: ${domain}=gcpidd
EOF

# Deploy Kubelet Image Service
cat <<EOF
Done!
Deploy the kubelet-image-service by running:
$ kustomize build k8s/ | kubectl apply -f-
And the ghost cms using an image from the docker registry proxy
$ kustomize build ghost-k8s/ | kubectl apply -f-
EOF
