#!/bin/bash -e

terraform init

terraform destroy --force
