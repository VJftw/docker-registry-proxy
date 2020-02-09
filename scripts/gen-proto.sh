#!/bin/sh -e

srcDefs="${PWD}/api/proto/v1"
destGo="${PWD}/pkg/genproto/v1"

mkdir -p "${destGo}"
uid=$(id -u)
gid=$(id -g)
docker run --rm \
    -v "${srcDefs}:/defs" \
    -v "${destGo}:/tmp/protogen" \
    --user "${uid}:${git}" \
    namely/protoc-all \
    -d /defs -l go -o /tmp/protogen
