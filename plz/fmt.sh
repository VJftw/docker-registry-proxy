#!/bin/bash
set -Eeuo pipefail

plz query alltargets //plz/format/... | plz run sequential -
