
## Running Locally

### Proxy on its own

```bash
# Docker Registry Proxy
$ go run cmd/docker-registry-proxy/main.go \
  --network_address="tcp://:8888" \
  --upstream_repository="https://index.docker.io"
```

### Static credentials

```bash

$ DRP_PLUGINS="auth-provider_static:upstream_auth" \
  go run cmd/docker-registry-proxy/main.go \
  --network_address="tcp://:8888" \
  --upstream_repository="https://index.docker.io" \
  --upstream_authentication="upstream_auth" \
  --upstream_auth_username="foo" \
  --upstream_auth_password="bar"
```

### GCP Instance Identity Document

```bash
$ DRP_PLUGINS="auth-verifier_gcp-instanceidentitydocument:gcpverifier auth-provider_static:upstream_auth" \
  go run cmd/docker-registry-proxy/main.go \
  --network_address="tcp://:8888" \
  --upstream_repository="https://index.docker.io" \
  --upstream_authentication="upstream_auth" \
  --upstream_auth_username="foo" \
  --upstream_auth_password="bar" \
  --authentication_verifier="_gcp:gcpverifier"

# Run on GCE instance
$ dist/auth-provider_gcp \
  --network_address="tcp://:8890" \
  --username="_gcp"  \
  --type="instanceidentitydocument" \
  --instanceidentitydocument_audience="foo"

# Verifier for Docker Registry Proxy
$ go run cmd/auth-verifier/gcp/main.go \
  --network_address="tcp://:8890" \
  --instanceidentitydocument_audiences="foo"

# Docker Registry Proxy
$ go run cmd/docker-registry-proxy/main.go \
  --network_address="tcp://:8888" \
  --grpc_insecure \
  --upstream_repository="https://index.docker.io" \
  --authentication_verifiers="_gcp=tcp://:8890"
```
