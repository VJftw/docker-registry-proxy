#!/bin/sh -xe

mv /home/docker/k8s-image-service /var/lib/k8s-image-service
mv /home/docker/k8s-image-service.conf /usr/lib/systemd/system/k8s-image-service.service

/home/docker/kubelet-patch.sh

systemctl daemon-reload
systemctl start k8s-image-service
systemctl restart kubelet
