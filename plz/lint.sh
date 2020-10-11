#!/bin/bash
set -Eeuo pipefail

./pleasew query alltargets //plz/lint/... | ./pleasew run sequential -
