#!/usr/bin/env bash
set -euo pipefail

if [[ ! -f secrets/jwt.key || ! -f secrets/jwt.pub ]]; then
  echo "Missing key pair files: secrets/jwt.key or secrets/jwt.pub" >&2
  exit 1
fi

openssl ec -in secrets/jwt.key -noout -text >/dev/null 2>&1
openssl ec -pubin -in secrets/jwt.pub -noout -text >/dev/null 2>&1

msg_file="$(mktemp)"
sig_file="$(mktemp)"
pub_file="$(mktemp)"
trap 'rm -f "$msg_file" "$sig_file" "$pub_file"' EXIT

printf 'jwt-key-check' > "$msg_file"
openssl pkey -in secrets/jwt.key -pubout -out "$pub_file" >/dev/null 2>&1
openssl dgst -sha256 -sign secrets/jwt.key -out "$sig_file" "$msg_file" >/dev/null 2>&1
openssl dgst -sha256 -verify "$pub_file" -signature "$sig_file" "$msg_file" >/dev/null 2>&1

echo "Key pair OK"
