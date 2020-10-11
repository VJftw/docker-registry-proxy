#!/bin/sh -xe

file="/etc/systemd/system/kubelet.service.d/10-kubeadm.conf"
endpoint="/var/run/k8s-image-service.sock"
grep "image-service-endpoint" "${file}" || sed -i "s|^.*kubelet.*|& -v=2 --image-service-endpoint=unix://${endpoint}|" "${file}"
