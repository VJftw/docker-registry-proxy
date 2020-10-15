FROM scratch

COPY /kubelet-image-service /bin/
VOLUME /tmp
VOLUME /var/run/kubelet

ENTRYPOINT [ "/bin/kubelet-image-service" ]
