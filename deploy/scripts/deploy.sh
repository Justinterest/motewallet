#!/usr/bin/env bash
#
# 一键发布：本地 build → export → scp 上传 → 服务器 import → compose 重启
# 支持三个服务单独发布。
#
# 用法:
#   ./deploy/scripts/deploy.sh [backend|frontend|admin|all]
#   ./deploy/scripts/deploy.sh backend --skip-build   # 仅上传并重启已有镜像包
#   ./deploy/scripts/deploy.sh all --skip-export      # 已 export 过则跳过
#
# 环境变量:
#   REMOTE_HOST=motewallet-prod
#   REMOTE_DIR=/home/ecs-user/motewallet-withdrawal
#   TAG=...
#   NEXT_PUBLIC_API_URL   默认空（同域 /api/，由 nginx 转发）
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

require_cmd docker ssh scp gzip

TARGET="all"
SKIP_BUILD=false
SKIP_EXPORT=false
SKIP_UPLOAD=false
SKIP_REMOTE=false

for arg in "$@"; do
  case "$arg" in
    backend|frontend|admin|all) TARGET="$arg" ;;
    --skip-build) SKIP_BUILD=true ;;
    --skip-export) SKIP_EXPORT=true ;;
    --skip-upload) SKIP_UPLOAD=true ;;
    --skip-remote) SKIP_REMOTE=true ;;
    -h|--help)
      sed -n '1,25p' "$0"
      exit 0
      ;;
    *)
      echo "未知参数: $arg" >&2
      exit 1
      ;;
  esac
done

TAG="$(resolve_tag)"
mkdir -p "$DIST_DIR"
echo "$TAG" >"$DIST_DIR/TAG"

SELECTED=()
while IFS= read -r svc; do
  SELECTED+=("$svc")
done < <(normalize_services "$TARGET")

if [[ "$SKIP_BUILD" != true ]]; then
  "$SCRIPT_DIR/build.sh" "$TARGET"
fi

if [[ "$SKIP_EXPORT" != true ]]; then
  TAG="$TAG" "$SCRIPT_DIR/export.sh" "$TARGET"
fi

if [[ "$SKIP_UPLOAD" != true ]]; then
  log "ensuring remote dir ${REMOTE_HOST}:${REMOTE_DIR}"
  ssh "$REMOTE_HOST" "mkdir -p '${REMOTE_DIR}/images' '${REMOTE_DIR}/deploy'"

  # sync compose + remote helper (env 不覆盖已有生产配置)
  scp "$DEPLOY_DIR/docker-compose.yml" "${REMOTE_HOST}:${REMOTE_DIR}/docker-compose.yml"
  scp "$SCRIPT_DIR/remote-up.sh" "${REMOTE_HOST}:${REMOTE_DIR}/remote-up.sh"
  ssh "$REMOTE_HOST" "chmod +x '${REMOTE_DIR}/remote-up.sh'"

  if ! ssh "$REMOTE_HOST" "test -f '${REMOTE_DIR}/.env'"; then
    log "remote .env missing — uploading .env.example as template"
    if [[ -f "$DEPLOY_DIR/.env" ]]; then
      scp "$DEPLOY_DIR/.env" "${REMOTE_HOST}:${REMOTE_DIR}/.env"
    else
      scp "$DEPLOY_DIR/.env.example" "${REMOTE_HOST}:${REMOTE_DIR}/.env"
      log "请尽快 ssh ${REMOTE_HOST} 编辑 ${REMOTE_DIR}/.env 填入生产配置"
    fi
  fi

  for svc in "${SELECTED[@]}"; do
    image="$(service_image "$svc")"
    archive="${DIST_DIR}/${image}-${TAG}.tar.gz"
    [[ -f "$archive" ]] || {
      echo "缺少镜像包: $archive（请先 build/export）" >&2
      exit 1
    }
    log "uploading $(basename "$archive")"
    scp "$archive" "${REMOTE_HOST}:${REMOTE_DIR}/images/"
  done
fi

if [[ "$SKIP_REMOTE" != true ]]; then
  svc_list=$(IFS=,; echo "${SELECTED[*]}")
  log "remote import + restart: ${svc_list} tag=${TAG}"
  ssh "$REMOTE_HOST" "TAG='${TAG}' '${REMOTE_DIR}/remote-up.sh' '${svc_list}'"
fi

log "deploy finished: ${SELECTED[*]} @ ${TAG}"
