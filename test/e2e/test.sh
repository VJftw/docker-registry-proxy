#!/bin/bash -e

cwd="${PWD}"
cd "test/e2e/docker-registry-proxy/"
./setup.sh

set +e
./test.sh
status=$?
set -e

./teardown.sh

cd "${cwd}"
if [ $status -eq 0 ]; then
    echo "Success!"
else
    echo "Failed: ${status}" 
    exit 1
fi
