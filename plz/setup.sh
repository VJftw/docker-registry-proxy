#!/bin/bash
set -Eeuo pipefail

echo "-> Configuring githooks"
git config core.hooksPath .githooks
