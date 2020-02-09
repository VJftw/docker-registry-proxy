## Docker Registry Proxy

 - Communicates with _Authorisation Verifiers_ to verify a given credential.
 - Routes credentials by username to a map of verifiers. e.g. `username: _gcp` routes to the GCP verifier.
 - Uses HMAC to generate short lived JWTs when authenticated for Docker Registry API v2


## Kubelet Image Service

 - Communicates with _Authorisation Providers_ to retrieve credentials for a given Docker Registry.
 - Routes requests based upon docker registry host. Deployers provide a semi-colon separate list of hostnames with a given endpoint as flags/env vars. e.g. `*.gcr.io;gcr.io=gcp-authenticator.svc.minikube`. Astrisks can be used to denote wildcards like wildcard certs.
