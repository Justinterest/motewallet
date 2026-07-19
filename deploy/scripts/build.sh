#!/usr/bin/env bash
#
# 本地构建 linux/amd64 镜像（可单独构建某一服务）
#
# 用法:
#   ./deploy/scripts/build.sh [backend|frontend|admin|all]
#
# 环境变量:
#   TAG                 镜像标签（默认: YYYYMMDD-<git-short>）
#   NEXT_PUBLIC_API_URL 前端 API baseURL；默认空（同域，由 nginx 转发 /api/）
#   PLATFORM            默认 linux/amd64
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

require_cmd docker

TARGET="${1:-all}"
TAG="$(resolve_tag)"
# 默认空：生产走 nginx 同域反代；本地直连可设 NEXT_PUBLIC_API_URL=http://127.0.0.1:8080
API_URL="${NEXT_PUBLIC_API_URL-}"

mkdir -p "$DIST_DIR"
echo "$TAG" >"$DIST_DIR/TAG"

log "platform=$PLATFORM tag=$TAG api=${API_URL:-<same-origin>}"

while IFS= read -r svc; do
  image="$(service_image "$svc")"
  ctx="$(service_dockerfile_ctx "$svc")"
  log "building ${image}:${TAG}"

  build_args=(--platform "$PLATFORM" --load)
  if [[ "$svc" == "frontend" || "$svc" == "admin" ]]; then
    build_args+=(--build-arg "NEXT_PUBLIC_API_URL=${API_URL}")
  fi

  docker buildx build \
    -t "${image}:${TAG}" \
    -t "${image}:latest" \
    "${build_args[@]}" \
    -f "${ctx}/Dockerfile" \
    "$ctx"

  log "built ${image}:${TAG}"
done < <(normalize_services "$TARGET")

log "done. images tagged with ${TAG}"
