#!/usr/bin/env bash
#
# 【服务器】加载镜像包并重启对应 compose 服务
# 由 deploy.sh 通过 ssh 调用，也可手动执行。
#
# 用法（在服务器 REMOTE_DIR）:
#   ./remote-up.sh backend
#   ./remote-up.sh frontend,admin
#   ./remote-up.sh all
#   TAG=20260719-abc1234 ./remote-up.sh backend
#
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGES_DIR="${ROOT_DIR}/images"
COMPOSE_FILE="${ROOT_DIR}/docker-compose.yml"
ENV_FILE="${ROOT_DIR}/.env"

cd "$ROOT_DIR"

if [[ ! -f "$COMPOSE_FILE" ]]; then
  echo "缺少 docker-compose.yml" >&2
  exit 1
fi
if [[ ! -f "$ENV_FILE" ]]; then
  echo "缺少 .env，请先从 .env.example 复制并填写" >&2
  exit 1
fi

TAG="${TAG:-}"
if [[ -z "$TAG" && -f "${IMAGES_DIR}/../deploy/dist/TAG" ]]; then
  TAG="$(cat "${IMAGES_DIR}/../deploy/dist/TAG")"
fi

INPUT="${1:-all}"
IFS=',' read -r -a RAW <<<"$INPUT"

SERVICES=()
for item in "${RAW[@]}"; do
  item="$(echo "$item" | tr -d '[:space:]')"
  case "$item" in
    all)
      SERVICES=(backend frontend admin)
      break
      ;;
    backend|frontend|admin)
      SERVICES+=("$item")
      ;;
    "")
      ;;
    *)
      echo "未知服务: $item" >&2
      exit 1
      ;;
  esac
done

if [[ ${#SERVICES[@]} -eq 0 ]]; then
  echo "未指定服务" >&2
  exit 1
fi

compose_tag_var() {
  case "$1" in
    backend) echo BACKEND_TAG ;;
    frontend) echo FRONTEND_TAG ;;
    admin) echo ADMIN_TAG ;;
  esac
}

set_env_tag() {
  local key="$1"
  local val="$2"
  if grep -q "^${key}=" "$ENV_FILE"; then
    # portable in-place replace
    local tmp
    tmp="$(mktemp)"
    awk -v k="$key" -v v="$val" 'BEGIN{FS=OFS="="} $1==k{$2=v} {print}' "$ENV_FILE" >"$tmp"
    mv "$tmp" "$ENV_FILE"
  else
    echo "${key}=${val}" >>"$ENV_FILE"
  fi
}

for svc in "${SERVICES[@]}"; do
  image="motewallet-${svc}"
  archive=""

  if [[ -n "$TAG" && -f "${IMAGES_DIR}/${image}-${TAG}.tar.gz" ]]; then
    archive="${IMAGES_DIR}/${image}-${TAG}.tar.gz"
  else
    # 取该服务最新包
    archive="$(ls -1t "${IMAGES_DIR}/${image}-"*.tar.gz 2>/dev/null | head -1 || true)"
    if [[ -n "$archive" ]]; then
      base="$(basename "$archive" .tar.gz)"
      TAG="${base#${image}-}"
    fi
  fi

  if [[ -z "$archive" || ! -f "$archive" ]]; then
    echo "找不到镜像包: ${IMAGES_DIR}/${image}-<tag>.tar.gz" >&2
    exit 1
  fi

  echo ">>> loading $(basename "$archive")"
  gunzip -c "$archive" | docker load

  # 确保 latest 也指向该 tag（compose 可用 latest 兜底）
  if docker image inspect "${image}:${TAG}" >/dev/null 2>&1; then
    docker tag "${image}:${TAG}" "${image}:latest"
  fi

  key="$(compose_tag_var "$svc")"
  set_env_tag "$key" "$TAG"
  echo ">>> ${key}=${TAG}"
done

echo ">>> docker compose up -d ${SERVICES[*]}"
# shellcheck disable=SC2086
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --no-deps "${SERVICES[@]}"

echo ">>> status"
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" ps
