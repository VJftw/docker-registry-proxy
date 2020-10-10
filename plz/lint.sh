#!/bin/bash
set -Eeuo pipefail

plz query alltargets //plz/lint/... | plz run sequential -
