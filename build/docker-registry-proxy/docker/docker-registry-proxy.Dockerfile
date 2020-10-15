FROM alpine:latest AS builder

RUN apk --update add ca-certificates

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY /docker-registry-proxy /bin/
VOLUME /tmp

ENTRYPOINT [ "/bin/docker-registry-proxy" ]
