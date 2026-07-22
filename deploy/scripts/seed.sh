#!/usr/bin/env bash
#
# 将 database/seeds 同步到服务器并执行种子 SQL（须在 migrate up 之后）。
#
# 用法:
#   ./deploy/scripts/seed.sh
#   ./deploy/scripts/seed.sh 001_seed_admin_and_config.sql   # 仅执行指定文件
#
# 环境变量:
#   REMOTE_HOST=motewallet-prod
#   REMOTE_DIR=/home/ecs-user/motewallet-withdrawal
#   SEED_LOCAL=1          # 在本机直连 DB（需本机可访问 MySQL）
#   MIGRATE_LOCAL=1       # 与 SEED_LOCAL 等效（与 migrate.sh 保持一致）
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

SEED_LOCAL="${SEED_LOCAL:-${MIGRATE_LOCAL:-0}}"
SEEDS_SRC="${ROOT_DIR}/database/seeds"
REMOTE_SEEDS="${REMOTE_DIR}/seeds"
DOTENV_DB_PY="${SCRIPT_DIR}/dotenv_db.py"
MYSQL_IMAGE="${MYSQL_IMAGE:-mysql:8.0}"

usage() {
  sed -n '1,18p' "$0"
  exit 1
}

TARGET_FILE="${1:-}"

require_cmd python3 docker

load_db_env() {
  local env_file="$1"
  [[ -f "$env_file" ]] || {
    echo "缺少环境文件: $env_file" >&2
    exit 1
  }
  eval "$(python3 "$DOTENV_DB_PY" "$env_file" exports)"
}

sync_remote_env() {
  if [[ ! -f "$DEPLOY_DIR/.env" ]]; then
    if ssh "$REMOTE_HOST" "test -f '${REMOTE_DIR}/.env'"; then
      log "local deploy/.env missing — using existing remote .env"
      return
    fi
    echo "缺少 deploy/.env，且服务器也没有 ${REMOTE_DIR}/.env" >&2
    exit 1
  fi
  log "sync deploy/.env -> ${REMOTE_HOST}:${REMOTE_DIR}/.env"
  ssh "$REMOTE_HOST" "mkdir -p '${REMOTE_DIR}'"
  scp "$DEPLOY_DIR/.env" "${REMOTE_HOST}:${REMOTE_DIR}/.env"
}

sync_seeds() {
  log "sync seeds -> ${REMOTE_HOST}:${REMOTE_SEEDS}"
  ssh "$REMOTE_HOST" "mkdir -p '${REMOTE_SEEDS}' '${REMOTE_DIR}/scripts'"
  rsync -az --delete "${SEEDS_SRC}/" "${REMOTE_HOST}:${REMOTE_SEEDS}/"
  scp "$DOTENV_DB_PY" "${REMOTE_HOST}:${REMOTE_DIR}/scripts/dotenv_db.py"
}

run_mysql_seed_dir() {
  local seeds_dir="$1"
  local env_file="$2"
  local only_file="${3:-}"

  load_db_env "$env_file"

  shopt -s nullglob
  local files=("$seeds_dir"/*.sql)
  shopt -u nullglob

  if [[ ${#files[@]} -eq 0 ]]; then
    log "no seed files in ${seeds_dir}"
    return
  fi

  if [[ -n "$only_file" ]]; then
    local matched=""
    for seed_file in "${files[@]}"; do
      if [[ "$(basename "$seed_file")" == "$only_file" ]]; then
        matched="$seed_file"
        break
      fi
    done
    if [[ -z "$matched" ]]; then
      echo "未找到种子文件: ${only_file}" >&2
      exit 1
    fi
    files=("$matched")
  fi

  for seed_file in "${files[@]}"; do
    log "applying seed $(basename "$seed_file")"
    docker run --rm --network host -i "$MYSQL_IMAGE" \
      mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" \
      <"$seed_file"
  done
}

run_remote() {
  sync_seeds
  sync_remote_env

  log "remote seed ${TARGET_FILE:-all}"
  ssh "$REMOTE_HOST" bash -s -- "$REMOTE_DIR" "$REMOTE_SEEDS" "$MYSQL_IMAGE" "${TARGET_FILE:-}" <<'EOS'
set -euo pipefail
REMOTE_DIR="$1"
REMOTE_SEEDS="$2"
MYSQL_IMAGE="$3"
ONLY_FILE="${4:-}"

ENV_FILE="${REMOTE_DIR}/.env"
DOTENV_DB_PY="${REMOTE_DIR}/scripts/dotenv_db.py"

eval "$(python3 "$DOTENV_DB_PY" "$ENV_FILE" exports)"

shopt -s nullglob
files=("$REMOTE_SEEDS"/*.sql)
shopt -u nullglob

if [[ ${#files[@]} -eq 0 ]]; then
  echo ">>> no seed files in ${REMOTE_SEEDS}"
  exit 0
fi

if [[ -n "$ONLY_FILE" ]]; then
  matched=""
  for seed_file in "${files[@]}"; do
    if [[ "$(basename "$seed_file")" == "$ONLY_FILE" ]]; then
      matched="$seed_file"
      break
    fi
  done
  if [[ -z "$matched" ]]; then
    echo "未找到种子文件: ${ONLY_FILE}" >&2
    exit 1
  fi
  files=("$matched")
fi

for seed_file in "${files[@]}"; do
  echo ">>> applying seed $(basename "$seed_file")"
  docker run --rm --network host -i "$MYSQL_IMAGE" \
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" \
    <"$seed_file"
done
EOS
}

run_local() {
  run_mysql_seed_dir "$SEEDS_SRC" "$DEPLOY_DIR/.env" "$TARGET_FILE"
}

case "${1:-}" in
  -h|--help) usage ;;
  "")
    ;;
  *.sql)
    TARGET_FILE="$1"
    ;;
  *)
    echo "未知参数: $1" >&2
    usage
    ;;
esac

if [[ "$SEED_LOCAL" == "1" ]]; then
  run_local
else
  require_cmd ssh rsync scp
  run_remote
fi

log "seed done"
