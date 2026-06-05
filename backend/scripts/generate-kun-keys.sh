#!/usr/bin/env bash
# Generate RSA2048 merchant key pair per KUN docs:
# https://opendocs.kun.global/docs/secret-key-configuration/merchant-key-preparation
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="${ROOT}/.kun-keys"

mkdir -p "$OUT_DIR"
cd "$OUT_DIR"

echo "Generating RSA2048 keys in ${OUT_DIR} ..."

openssl genrsa -out app_private_key.pem 2048
openssl pkcs8 -topk8 -inform PEM -in app_private_key.pem -outform PEM -nocrypt -out app_private_key_pkcs8.pem
openssl rsa -in app_private_key.pem -pubout -out app_public_key.pem

# Base64 DER (for .env / KUN console paste — same format as dashboard)
openssl rsa -in app_private_key_pkcs8.pem -pubout -outform DER | base64 | tr -d '\n' > app_public_key.b64.txt
openssl pkcs8 -in app_private_key_pkcs8.pem -outform DER -nocrypt | base64 | tr -d '\n' > app_private_key_pkcs8.b64.txt

echo ""
echo "Done. Files:"
ls -1 app_*.pem app_*.b64.txt
echo ""
echo "Next steps:"
echo "  1. Upload app_public_key.pem (or app_public_key.b64.txt) to KUN merchant console"
echo "  2. Set in backend/.env (Go — use PKCS#8-derived DER base64):"
echo "       KUN_PRIVATE_KEY=<contents of app_private_key_pkcs8.b64.txt>"
echo "       KUN_PUBLIC_KEY=<contents of app_public_key.b64.txt>"
echo "  3. Java projects: use app_private_key_pkcs8.pem (PKCS#8 PEM), not app_private_key.pem"
echo "  3. Never commit .kun-keys/ (already in .gitignore)"
