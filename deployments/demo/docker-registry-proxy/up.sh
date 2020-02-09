#!/bin/bash -e

### SETTINGS
acme_email="meetthevj@gmail.com"
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

# Place SA into k8s dir
jq -r '.resources[].instances[].attributes.private_key' \
  terraform/terraform.tfstate \
  | grep -v null \
  | base64 -d - \
  > ./k8s/gcp-sa.json

project=$(jq -r '.resources[].instances[].attributes.project' terraform/terraform.tfstate | grep -v null | head -n1)

gcloud container clusters get-credentials docker-registry-proxy --region europe-west1 --project "${project}"

kubectl label nodes kubernetes.io/os=linux --all

kubectl create namespace ingress-nginx
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.28.0/deploy/static/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.28.0/deploy/static/provider/cloud-generic.yaml

kubectl create namespace cert-manager
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.13.0/cert-manager.yaml

kubectl --namespace cert-manager wait --for condition=available deployment cert-manager-webhook

cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt
spec:
  acme:
    email: ${acme_email}
    # server: https://acme-staging-v02.api.letsencrypt.org/directory
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-issuer-account-key
    solvers:
    - http01:
        ingress:
          class: nginx
EOF

cat <<EOF > k8s/ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: docker-registry-proxy
  namespace: docker-registry-proxy
  annotations:
    # add an annotation indicating the issuer to use.
    cert-manager.io/cluster-issuer: letsencrypt
spec:
  rules:
    - host: ${domain}
      http:
        paths:
          - backend:
              serviceName: docker-registry-proxy
              servicePort: 80
            path: /
  tls: # < placing a host in the TLS config will indicate a certificate should be created
    - hosts:
        - ${domain}
      secretName: cert-${domain//./-} # < cert-manager will store the created certificate in this secret.
EOF

cat <<EOF > k8s/gcp-provider.yaml
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: docker-registry-proxy
spec:
  template:
    spec:
      containers:
        - name: docker-registry-proxy
          env:
            - name: DRP_UPSTREAM_REPOSITORY
              value: https://gcr.io/${project}
            - name: DRP_PLUGINS
              value: "auth-provider_static:upstream_static auth-verifier_gcp-instanceidentitydocument:gcpverifier"
            - name: DRP_UPSTREAM_AUTHENTICATION
              value: upstream_static
            - name: DRP_UPSTREAM_STATIC_USERNAME
              value: _json_key
            - name: DRP_UPSTREAM_STATIC_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: gcp-sa
                  key: gcp-sa.json
            # - name: DRP_AUTHENTICATION_VERIFIER
            #   value: _gcpidd:gcpverifier
            # - name: DRP_GCPVERIFIER_PROJECT_IDS
            #   value: <project id>

EOF

# Deploy Docker Registry Proxy
cat <<EOF
Done! You'll likely need to update you dns record to point to the LoadBalancer IP.
You can get this with:
$ kubectl --namespace ingress-nginx get svc
Deploy the docker registry proxy by running:
$ kustomize build k8s/ | kubectl apply -f-
EOF
