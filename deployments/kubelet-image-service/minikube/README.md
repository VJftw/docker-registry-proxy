## Minikube Demo instructions

```
# Setup
scripts/build.sh
scripts/minikube.sh

# view logs in minikube
sudo journalctl -u kubelet -f
sudo journalctl -u k8s-image-service -f


# Run image
kubectl run nginx --image=docker.io/library/nginx --replicas=1
kubectl get pod -w

kubectl delete deployment nginx
```
