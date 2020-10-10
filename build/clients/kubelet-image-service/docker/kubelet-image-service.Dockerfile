FROM scratch

COPY dist/kubelet-image-service /bin/
COPY dist/plugin* /plugins/
VOLUME /tmp
VOLUME /var/run/kubelet

ENTRYPOINT [ "/bin/kubelet-image-service" ]
