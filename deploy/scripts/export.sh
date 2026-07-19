#!/usr/bin/env bash
#
# 导出镜像为压缩包（docker save | gzip）
#
# 用法:
#   ./deploy/scripts/export.sh [backend|frontend|admin|all]
#
# 输出:
#   deploy/dist/motewallet-<svc>-<tag>.tar.gz
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

require_cmd docker gzip

TARGET="${1:-all}"
if [[ -n "${TAG:-}" ]]; then
  :
elif [[ -f "$DIST_DIR/TAG" ]]; then
  TAG="$(cat "$DIST_DIR/TAG")"
else
  TAG="$(resolve_tag)"
fi

mkdir -p "$DIST_DIR"

while IFS= read -r svc; do
  image="$(service_image "$svc")"
  out="${DIST_DIR}/${image}-${TAG}.tar.gz"
  log "exporting ${image}:${TAG} -> ${out}"
  docker save "${image}:${TAG}" | gzip -1 >"$out"
  ls -lh "$out"
done < <(normalize_services "$TARGET")

log "export done"
