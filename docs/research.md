# How Google Kubernetes Engine (GKE) transparently authenticates with Google Container Registry (GCR)

## Preface
If we push an image to a private GCR and then start a GKE cluster, we are able to use images from that private GCR transparently in GKE:
```
# Push an image to GCR
meetthevj@cloudshell:~ (gke-test-257811)$ docker pull kindest/node:v1.15.3
meetthevj@cloudshell:~ (gke-test-257811)$ docker tag kindest/node:v1.15.3 eu.gcr.io/gke-test-257811/kindest/node:v1.15.3
meetthevj@cloudshell:~ (gke-test-257811)$ docker push eu.gcr.io/gke-test-257811/kindest/node:v1.15.3

# Run a Kubernetes container with that image
meetthevj@cloudshell:~ (gke-test-257811)$ kubectl run test --image=eu.gcr.io/gke-test-257811/kindest/node:v1.15.3 --replicas=1
meetthevj@cloudshell:~ (gke-test-257811)$ kubectl get pod
Events:
  Type     Reason     Age                From                                          Message
  ----     ------     ----               ----                                          -------
  Normal   Scheduled  109s               default-scheduler                             Successfully assigned default/test-7878b75597-n2mmm to gke-test-default-pool-e36d59bc-86n7
  Normal   Pulling    109s               kubelet, gke-test-default-pool-e36d59bc-86n7  pulling image "eu.gcr.io/gke-test-257811/kindest/node:v1.15.3"
  Normal   Pulled     76s                kubelet, gke-test-default-pool-e36d59bc-86n7  Successfully pulled image "eu.gcr.io/gke-test-257811/kindest/node:v1.15.3"
  Normal   Created    29s (x4 over 73s)  kubelet, gke-test-default-pool-e36d59bc-86n7  Created container
  Normal   Started    29s (x4 over 73s)  kubelet, gke-test-default-pool-e36d59bc-86n7  Started container
```


## Discovery

### GCR does not authenticate nodes at the provider/infrastructure level
If this was the case, we should be able to run `docker pull eu.gcr.io/gke-test-257811/kindest/node:v1.15.3` on the node successfully as authentication is handled transparently from the node's perspective. The below proves this:
```
meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ docker pull eu.gcr.io/gke-test-257811/kindest/node:v1.15.3
Error response from daemon: unauthorized: You don't have the needed permissions to perform this operation, and you may have invalid credentials. To authenticate your request, follow the steps in: https://cloud.google.com/container-registry/docs/advanced-authentication
```

### Kubelet performs authentication
https://github.com/kubernetes/kubernetes/tree/master/pkg/credentialprovider/gcp

Kubelet performs authentication via this package (there are sibling packages for `aws` and `azure` so assume that they work in a similar way). This package simply obtains the service account token from the node metadata to use as the password for the docker registry `eu.gcr.io`.

#### Proof by using `docker login`
```
meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ curl -H "Metadata-Flavor: Google" 
\ http://metadata.google.internal./computeMetadata/v1/instance/service-accounts/default/token
{"access_token":"xxxxxx","expires_in":2895,"token_type":"Bearer"}

meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ docker login eu.gcr.io
Username: _token
Password: xxxxxx
WARNING! Your password will be stored unencrypted in /home/meetthevj/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded

meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ docker pull eu.gcr.io/gke-test-257811/kindest/node:v1.15.3
v1.15.3: Pulling from gke-test-257811/kindest/node
09fe536fe3e4: Pull complete
4326ed393c12: Pull complete
73a8f76105d8: Pull complete
61a147631503: Pull complete
ebf284b29a9b: Pull complete
45edd976aead: Pull complete
25a49a5ef18f: Pull complete
ac9964da2ccf: Pull complete
a9e631f193b5: Pull complete
3d2abb89014a: Pull complete
c4dcd9c54845: Pull complete
1a7ee4f67711: Pull complete
Digest: sha256:27e388752544890482a86b90d8ac50fcfa63a2e8656a96ec5337b902ec8e5157
Status: Downloaded newer image for eu.gcr.io/gke-test-257811/kindest/node:v1.15.3
```

### Going Forward

We definitely don't want to have to add to kubelet code even though kubelet has has vendor specific logic to handle authentication. We should attempt to find a method to specify other transparent sources of credentials. E.g. There is an `--image-service-endpoint` on kubelet which we may be able to implement (https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/, https://github.com/kubernetes/kubernetes/pull/30212) to perform custom authentication against an external private registry. (as demonstrated in this repository)

### Multi-cloud single-registry authentication

#### Instance Identity Documents

##### GCP GCE
https://cloud.google.com/compute/docs/instances/verifying-instance-identity

On GCE, we can obtain an instance identity JWT via:
```
meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ curl -H "Metadata-Flavor: Google" \
> 'http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=registry.example.org&format=full&licenses=TRUE'
eyJhbGciOiJSUzI1NiIsImtpZCI6IjljZWY1MzQwNjQyYjE1N2ZhOGE0ZjBkODc0ZmU3OTAwMzYyZDgyZGIiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJyZWdpc3RyeS5leGFtcGxlLm9yZyIsImF6cCI6IjExMDQ4MTU3MDQ4MTgwNzE2ODA0MyIsImVtYWlsIjoiMjI0MTg5MDMzODg2LWNvbXB1dGVAZGV2ZWxvcGVyLmdzZXJ2aWNlYWNjb3VudC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZXhwIjoxNTcyNzA2ODI3LCJnb29nbGUiOnsiY29tcHV0ZV9lbmdpbmUiOnsiaW5zdGFuY2VfY3JlYXRpb25fdGltZXN0YW1wIjoxNTcyNjk0MTYzLCJpbnN0YW5jZV9pZCI6IjY1NjAyMzgyOTM5ODc5Mzk5NjUiLCJpbnN0YW5jZV9uYW1lIjoiZ2tlLXRlc3QtZGVmYXVsdC1wb29sLWUzNmQ1OWJjLTg2bjciLCJsaWNlbnNlX2lkIjpbIjEwMDEwMTAiLCIxMDAxMDAzIiwiMTY2NzM5NzEyMjMzNjU4NzY2IiwiNjg4MDA0MTk4NDA5NjU0MDEzMiJdLCJwcm9qZWN0X2lkIjoiZ2tlLXRlc3QtMjU3ODExIiwicHJvamVjdF9udW1iZXIiOjIyNDE4OTAzMzg4Niwiem9uZSI6ImV1cm9wZS13ZXN0Mi1iIn19LCJpYXQiOjE1NzI3MDMyMjcsImlzcyI6Imh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbSIsInN1YiI6IjExMDQ4MTU3MDQ4MTgwNzE2ODA0MyJ9.Ghc-6R30ioMZZuFKYhOo1_lnSwnsJB1KwqQyXqXxRgiIQJJE83KUD3Sfhs8rAtde0S0e71XcHF_EcQMjwO5xOt3hrVZd56ENt-L0kjQZXACNic2pq7nh-qw8nU1tjWOPrxxb0t5AX1lu8CTMmMVsjTQq7BkG0ruRumPoZUgq8vLFO7AgRRfC3wBtGPyRzcTrnfB5GORT9et3VQNa773MeMcbRz8ggdWxGncMYWvXd57ZOMzS_LAHbgECr3kvBT1_ki6jcmINeNbxz7WxCdmoBR-gCP24o-p9mYZPjgpDr8a1bG6_HBOROUdwoIXbhRvGDpoSqhaXPdL7mZ2rFtsmcA
```
This token has the following payload:
```
{
  "aud": "registry.example.org",
  "azp": "110481570481807168043",
  "email": "224189033886-compute@developer.gserviceaccount.com",
  "email_verified": true,
  "exp": 1572706827,
  "google": {
    "compute_engine": {
      "instance_creation_timestamp": 1572694163,
      "instance_id": "6560238293987939965",
      "instance_name": "gke-test-default-pool-e36d59bc-86n7",
      "license_id": [
        "1001010",
        "1001003",
        "166739712233658766",
        "6880041984096540132"
      ],
      "project_id": "gke-test-257811",
      "project_number": 224189033886,
      "zone": "europe-west2-b"
    }
  },
  "iat": 1572703227,
  "iss": "https://accounts.google.com",
  "sub": "110481570481807168043"
}
```
To verify this document's authenticity, our registry verifies it against GCP's public keys (example in above documentation) and then whitelists based on `project_id`/`project_number`, `zone`, `instance_name` etc.

##### AWS EC2
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html

On AWS EC2, we can obtain an instance identity document via:
```
curl http://169.254.169.254/latest/dynamic/instance-identity/document
{
    "devpayProductCodes": null,
    "marketplaceProductCodes": ["1abc2defghijklm3nopqrs4tu"],
    "availabilityZone": "us-west-2b",
    "privateIp": "10.158.112.84",
    "version": "2017-09-30",
    "instanceId": "i-1234567890abcdef0",
    "billingProducts": null,
    "instanceType": "t2.micro",
    "accountId": "123456789012",
    "imageId": "ami-5fb8c835",
    "pendingTime": "2016-11-19T16: 32: 11Z",
    "architecture": "x86_64",
    "kernelId": null,
    "ramdiskId": null,
    "region": "us-west-2"
}
```
And the document's PKCS7 signature via:
```
curl http://169.254.169.254/latest/dynamic/instance-identity/pkcs7
```
To verify the document's authenticity, our registry verifies it against AWS's regional public keys (example in above documentation) and then whitelists based on `accountId`, `region`, `availabilityZone` etc.


# Appendix

## Files of interest
### /etc/default/docker
```
meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ cat /etc/default/docker
DOCKER_OPTS="-p /var/run/docker.pid --iptables=false --ip-masq=false --log-level=warn --bip=169.254.123.1/24 --registry-mirror=https://mirror.gcr.io --log-driver=json-file --log-opt=max-size=10m --log-opt=max-file=5 --insecure-registry 10.0.0.0/8"
```
### /etc/default/kubelet
```
meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ cat /etc/default/kubelet
KUBELET_OPTS="--v=2 --cloud-provider=gce --experimental-check-node-capabilities-before-mount=true --allow-privileged=true --experimental-mounter-path=/home/kubernetes/containerized_mounter/mounter --cert-dir=/var/lib/kubelet/pki/ --cni-bin-dir=/home/kubernetes/bin --kubeconfig=/var/lib/kubelet/kubeconfig --experimental-kernel-memcg-notification=true --max-pods=110 --network-plugin=kubenet --node-labels=beta.kubernetes.io/fluentd-ds-ready=true,cloud.google.com/gke-nodepool=default-pool,cloud.google.com/gke-os-distribution=cos --volume-plugin-dir=/home/kubernetes/flexvolume --bootstrap-kubeconfig=/var/lib/kubelet/bootstrap-kubeconfig --node-status-max-images=25 --registry-qps=10 --registry-burst=20 --config /home/kubernetes/kubelet-config.yaml --pod-sysctls='net.core.somaxconn=1024,net.ipv4.conf.all.accept_redirects=0,net.ipv4.conf.all.forwarding=1,net.ipv4.conf.all.route_localnet=1,net.ipv4.conf.default.forwarding=1,net.ipv4.ip_forward=1,net.ipv4.tcp_fin_timeout=60,net.ipv4.tcp_keepalive_intvl=75,net.ipv4.tcp_keepalive_probes=9,net.ipv4.tcp_keepalive_time=7200,net.ipv4.tcp_max_syn_backlog=128,net.ipv4.tcp_max_tw_buckets=16384,net.ipv4.tcp_syn_retries=6,net.ipv4.tcp_tw_reuse=0,net.netfilter.nf_conntrack_generic_timeout=600,net.netfilter.nf_conntrack_tcp_timeout_close_wait=3600,net.netfilter.nf_conntrack_tcp_timeout_established=86400'"
KUBE_COVERAGE_FILE="/var/log/kubelet.cov"
```
### /usr/lib/systemd/system/docker.service
```
meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ cat /usr/lib/systemd/system/docker.service
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
PartOf=containerd.service
After=network-online.target docker.socket firewalld.service containerd.service
Wants=network-online.target containerd.service
Requires=docker.socket

[Service]
Type=notify
EnvironmentFile=-/etc/default/docker
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
ExecStartPre=/bin/sh -c 'if [[ -f /var/lib/docker/daemon.json ]]; then cp -f /var/lib/docker/daemon.json /etc/docker/daemon.json; fi'
ExecStart=/usr/bin/dockerd --registry-mirror=https://mirror.gcr.io --host=fd:// --containerd=/var/run/containerd/containerd.sock $DOCKER_OPTS
ExecReload=/bin/kill -s HUP $MAINPID
ExecStopPost=/bin/echo "docker daemon exited"
OOMScoreAdjust=-999
LimitNOFILE=1048576
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNPROC=infinity
LimitCORE=infinity
# Uncomment TasksMax if your systemd version supports it.
# Only systemd 226 and above support this version.
TasksMax=infinity
TimeoutStartSec=0
# set delegate yes so that systemd does not reset the cgroups of docker containers
Delegate=yes
# kill only the docker process, not all processes in the cgroup
KillMode=process
# restart the docker process if it exits
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```
### /etc/systemd/system/kubelet.service
```
meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ cat /etc/systemd/system/kubelet.service
[Unit]
Description=Kubernetes kubelet
Requires=network-online.target
After=network-online.target

[Service]
Restart=always
RestartSec=10
EnvironmentFile=/etc/default/kubelet
ExecStart=/home/kubernetes/bin/kubelet $KUBELET_OPTS

[Install]
WantedBy=multi-user.target
```

### `journalctl -u kubelet` with -v=4
```
meetthevj@gke-test-default-pool-e36d59bc-86n7 ~ $ journalctl -u kubelet
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.469365   19859 kuberuntime_manager.go:750] Creating container &Container{Name:test,Image:eu.gcr.io/gke-test-257811/kindest/node:v1.15.3,Command:[],Args:[],WorkingDir:,Ports:[],Env:[],Resources:ResourceRequirements{Limits:ResourceList{},Requests:ResourceList{cpu: {{100 -3} {<nil>} 100m DecimalSI},},},VolumeMounts:[{default-token-r6ssk true /var/run/secrets/kubernetes.io/serviceaccount  <nil>}],LivenessProbe:nil,ReadinessProbe:nil,Lifecycle:nil,TerminationMessagePath:/dev/termination-log,ImagePullPolicy:IfNotPresent,SecurityContext:nil,Stdin:false,StdinOnce:false,TTY:false,EnvFrom:[],TerminationMessagePolicy:File,VolumeDevices:[],} in pod test-7878b75597-295wc_default(aca32e7d-fd74-11e9-96b3-42010a9a0084)
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.470540   19859 provider.go:116] Refreshing cache for provider: *credentialprovider.defaultDockerConfigProvider
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.470807   19859 config.go:131] looking for config.json at /var/lib/kubelet/config.json
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.471013   19859 config.go:131] looking for config.json at /config.json
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.471235   19859 config.go:131] looking for config.json at /.docker/config.json
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.471431   19859 config.go:131] looking for config.json at /.docker/config.json
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.471620   19859 config.go:101] looking for .dockercfg at /var/lib/kubelet/.dockercfg
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.471808   19859 config.go:101] looking for .dockercfg at /.dockercfg
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.472004   19859 config.go:101] looking for .dockercfg at /.dockercfg
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.472227   19859 config.go:101] looking for .dockercfg at /.dockercfg
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.472417   19859 provider.go:86] Unable to parse Docker config file: couldn't find valid .dockercfg after checking in [/var/lib/kubelet   /]
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.473105   19859 server.go:459] Event(v1.ObjectReference{Kind:"Pod", Namespace:"default", Name:"test-7878b75597-295wc", UID:"aca32e7d-fd74-11e9-96b3-42010a9a0084", APIVersion:"v1", ResourceVersion:"22964", FieldPath:"spec.containers{test}"}): type: 'Normal' reason: 'Pulling' pulling image "eu.gcr.io/gke-test-257811/kindest/node:v1.15.3"
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.475201   19859 provider.go:116] Refreshing cache for provider: *gcp_credentials.dockerConfigKeyProvider
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.476157   19859 config.go:191] body of failing http response: &{0x6dbb60 0xc000c89600 0x6e6510}
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: E1102 13:28:37.476445   19859 metadata.go:142] while reading 'google-dockercfg' metadata: http status code: 404 while fetching url http://metadata.google.internal./computeMetadata/v1/instance/attributes/google-dockercfg
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.476717   19859 provider.go:116] Refreshing cache for provider: *gcp_credentials.dockerConfigUrlKeyProvider
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: I1102 13:28:37.477809   19859 config.go:191] body of failing http response: &{0x6dbb60 0xc000c89700 0x6e6510}
Nov 02 13:28:37 gke-test-default-pool-e36d59bc-86n7 kubelet[19859]: E1102 13:28:37.478111   19859 metadata.go:159] while reading 'google-dockercfg-url' metadata: http status code: 404 while fetching url http://metadata.google.internal./computeMetadata/v1/instance/attributes/google-dockercfg-url
