#!/usr/bin/env bash
#
# 将 database/migrations 同步到服务器并用 golang-migrate 应用到线上库。
# 默认在 motewallet-prod 上执行（RDS 通常只对内网开放）。
#
# 用法:
#   ./deploy/scripts/migrate.sh version
#   ./deploy/scripts/migrate.sh up
#   ./deploy/scripts/migrate.sh down 1
#   ./deploy/scripts/migrate.sh force 7          # 修复脏版本（慎用）
#   ./deploy/scripts/migrate.sh create-db        # 首次建库
#
# 环境变量:
#   REMOTE_HOST=motewallet-prod
#   REMOTE_DIR=/home/ecs-user/motewallet-withdrawal
#   MIGRATE_LOCAL=1   # 在本机直连 DB（需本机可访问 MySQL）
#   MIGRATE_IMAGE=migrate/migrate:v4.18.3
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

MIGRATE_IMAGE="${MIGRATE_IMAGE:-migrate/migrate:v4.18.3}"
MIGRATE_LOCAL="${MIGRATE_LOCAL:-0}"
MIGRATIONS_SRC="${ROOT_DIR}/database/migrations"
REMOTE_MIGRATIONS="${REMOTE_DIR}/migrations"
DOTENV_DB_PY="${SCRIPT_DIR}/dotenv_db.py"

usage() {
  sed -n '1,22p' "$0"
  exit 1
}

ACTION="${1:-}"
[[ -n "$ACTION" ]] || usage
shift || true

require_cmd python3

load_db_env() {
  local env_file="$1"
  [[ -f "$env_file" ]] || {
    echo "缺少环境文件: $env_file" >&2
    exit 1
  }
  # 不用 source .env（密码里的 * $ 等会炸）；只安全导出 DB_*
  eval "$(python3 "$DOTENV_DB_PY" "$env_file" exports)"
}

mysql_dsn_from_env_file() {
  python3 "$DOTENV_DB_PY" "$1" dsn
}

sync_migrations() {
  log "sync migrations -> ${REMOTE_HOST}:${REMOTE_MIGRATIONS}"
  ssh "$REMOTE_HOST" "mkdir -p '${REMOTE_MIGRATIONS}' '${REMOTE_DIR}/scripts'"
  rsync -az --delete "${MIGRATIONS_SRC}/" "${REMOTE_HOST}:${REMOTE_MIGRATIONS}/"
  scp "$DOTENV_DB_PY" "${REMOTE_HOST}:${REMOTE_DIR}/scripts/dotenv_db.py"
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

run_remote() {
  local remote_action="$1"
  shift

  sync_migrations
  sync_remote_env

  log "remote migrate ${remote_action} $*"
  ssh "$REMOTE_HOST" bash -s -- "$REMOTE_DIR" "$REMOTE_MIGRATIONS" "$MIGRATE_IMAGE" "$remote_action" "$@" <<'EOS'
set -euo pipefail
REMOTE_DIR="$1"
REMOTE_MIGRATIONS="$2"
MIGRATE_IMAGE="$3"
ACTION="$4"
shift 4

ENV_FILE="${REMOTE_DIR}/.env"
DOTENV_DB_PY="${REMOTE_DIR}/scripts/dotenv_db.py"

eval "$(python3 "$DOTENV_DB_PY" "$ENV_FILE" exports)"

if [[ "$ACTION" == "create-db" ]]; then
  SQL="$(python3 "$DOTENV_DB_PY" "$ENV_FILE" create-sql)"
  if command -v mysql >/dev/null 2>&1; then
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "$SQL"
  else
    docker run --rm --network host mysql:8.0 \
      mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "$SQL"
  fi
  echo "database ${DB_NAME} ready"
  exit 0
fi

DSN="$(python3 "$DOTENV_DB_PY" "$ENV_FILE" dsn)"

docker run --rm --network host \
  -v "${REMOTE_MIGRATIONS}:/migrations:ro" \
  "$MIGRATE_IMAGE" \
  -path=/migrations \
  -database "$DSN" \
  "$ACTION" "$@"
EOS
}

run_local() {
  local local_action="$1"
  shift

  require_cmd docker
  load_db_env "$DEPLOY_DIR/.env"

  if [[ "$local_action" == "create-db" ]]; then
    log "create database ${DB_NAME} (local)"
    SQL="$(python3 "$DOTENV_DB_PY" "$DEPLOY_DIR/.env" create-sql)"
    docker run --rm --network host mysql:8.0 \
      mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "$SQL"
    log "database ${DB_NAME} ready"
    return
  fi

  local dsn
  dsn="$(mysql_dsn_from_env_file "$DEPLOY_DIR/.env")"
  log "local migrate ${local_action} $*"
  docker run --rm --network host \
    -v "${MIGRATIONS_SRC}:/migrations:ro" \
    "$MIGRATE_IMAGE" \
    -path=/migrations \
    -database "$dsn" \
    "$local_action" "$@"
}

case "$ACTION" in
  -h|--help) usage ;;
  up|down|version|force|goto|drop|create-db)
    if [[ "$MIGRATE_LOCAL" == "1" ]]; then
      run_local "$ACTION" "$@"
    else
      require_cmd ssh rsync scp
      run_remote "$ACTION" "$@"
    fi
    ;;
  *)
    echo "未知操作: $ACTION" >&2
    usage
    ;;
esac

log "migrate ${ACTION} done"
