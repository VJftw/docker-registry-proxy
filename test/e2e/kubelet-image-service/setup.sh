#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cwd=$(pwd)

# Set up Docker Registry Proxy
cd "${DIR}/../docker-registry-proxy"
./setup.sh "${DIR}/docker-registry-proxy.tfvars"



# Tear down Docker Registry Proxy
cd "${DIR}/../docker-registry-proxy"
./teardown.sh

cd "${cwd}"