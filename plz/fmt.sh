#!/bin/bash
set -Eeuo pipefail

./pleasew query alltargets //plz/format/... | ./pleasew run sequential -
