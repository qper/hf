# Habits Flow

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](#)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](#)

This repository contains the backend, frontend, Helm charts, documentation, scripts, and CI workflow scaffolding.

## JWT key generation

Generate a local JWT key pair for development:

```bash
bash scripts/generate-keys.sh
```

The script creates:

- `secrets/jwt.key`
- `secrets/jwt.pub`

You can verify the generated key pair with:

```bash
bash scripts/verify-keys.sh
```

To pass the keys to Kubernetes, create a Secret from the generated files:

```bash
kubectl create secret generic jwt-keys \
  --from-file=jwt.key=secrets/jwt.key \
  --from-file=jwt.pub=secrets/jwt.pub
```

Then mount the Secret into the application pod and point the service to the mounted paths.
