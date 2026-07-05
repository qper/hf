#!/usr/bin/env bash
set -euo pipefail

mkdir -p secrets

openssl ecparam -name prime256v1 -genkey -noout -out secrets/jwt.key
openssl ec -in secrets/jwt.key -pubout -out secrets/jwt.pub
chmod 600 secrets/jwt.key
chmod 644 secrets/jwt.pub

echo "Generated secrets/jwt.key and secrets/jwt.pub"
