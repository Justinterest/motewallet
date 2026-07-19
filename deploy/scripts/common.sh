#!/usr/bin/env bash
# Shared helpers for deploy scripts.
set -euo pipefail

DEPLOY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ROOT_DIR="$(cd "$DEPLOY_DIR/.." && pwd)"
DIST_DIR="${DIST_DIR:-$DEPLOY_DIR/dist}"
REMOTE_HOST="${REMOTE_HOST:-motewallet-prod}"
REMOTE_DIR="${REMOTE_DIR:-/home/ecs-user/motewallet-withdrawal}"
PLATFORM="${PLATFORM:-linux/amd64}"

SERVICES=(backend frontend admin)

IMAGE_PREFIX=motewallet

service_image() {
  local svc="$1"
  echo "${IMAGE_PREFIX}-${svc}"
}

service_dockerfile_ctx() {
  local svc="$1"
  echo "${ROOT_DIR}/${svc}"
}

resolve_tag() {
  if [[ -n "${TAG:-}" ]]; then
    echo "$TAG"
    return
  fi
  local short
  short="$(git -C "$ROOT_DIR" rev-parse --short HEAD 2>/dev/null || true)"
  if [[ -n "$short" ]]; then
    echo "$(date +%Y%m%d)-${short}"
  else
    date +%Y%m%d%H%M
  fi
}

normalize_services() {
  local input="${1:-all}"
  case "$input" in
    all)
      printf '%s\n' "${SERVICES[@]}"
      ;;
    backend|frontend|admin)
      echo "$input"
      ;;
    *)
      echo "未知服务: $input（可选: backend | frontend | admin | all）" >&2
      exit 1
      ;;
  esac
}

compose_tag_var() {
  local svc="$1"
  case "$svc" in
    backend) echo BACKEND_TAG ;;
    frontend) echo FRONTEND_TAG ;;
    admin) echo ADMIN_TAG ;;
  esac
}

log() {
  printf '>>> %s\n' "$*"
}

require_cmd() {
  local c
  for c in "$@"; do
    command -v "$c" >/dev/null 2>&1 || {
      echo "缺少命令: $c" >&2
      exit 1
    }
  done
}
