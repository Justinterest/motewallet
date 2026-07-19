#!/usr/bin/env python3
"""Safely read DB_* from a .env file without bash source (handles * $ ` etc)."""
from __future__ import annotations

import argparse
import shlex
import sys
import urllib.parse
from pathlib import Path


def load_env(path: Path) -> dict[str, str]:
    env: dict[str, str] = {}
    for raw in path.read_text(encoding="utf-8").splitlines():
        line = raw.rstrip("\r")
        if not line.strip() or line.lstrip().startswith("#"):
            continue
        if "=" not in line:
            continue
        key, val = line.split("=", 1)
        key = key.strip()
        if len(val) >= 2 and val[0] == val[-1] and val[0] in "'\"":
            val = val[1:-1]
        env[key] = val
    return env


def require_db(env: dict[str, str]) -> dict[str, str]:
    keys = ("DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME")
    missing = [k for k in keys if not env.get(k)]
    if missing:
        raise SystemExit(f"missing in .env: {', '.join(missing)}")
    return {k: env[k] for k in keys}


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("env_file")
    parser.add_argument(
        "mode",
        choices=("exports", "dsn", "create-sql"),
        help="exports: shell export lines; dsn: migrate mysql URL; create-sql: CREATE DATABASE",
    )
    args = parser.parse_args()
    db = require_db(load_env(Path(args.env_file)))

    if args.mode == "exports":
        for k, v in db.items():
            print(f"export {k}={shlex.quote(v)}")
        return

    if args.mode == "create-sql":
        name = db["DB_NAME"].replace("`", "``")
        print(
            f"CREATE DATABASE IF NOT EXISTS `{name}` "
            "CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        )
        return

    user = urllib.parse.quote(db["DB_USER"], safe="")
    password = urllib.parse.quote(db["DB_PASSWORD"], safe="")
    host = db["DB_HOST"]
    port = db["DB_PORT"]
    name = urllib.parse.quote(db["DB_NAME"], safe="")
    print(
        f"mysql://{user}:{password}@tcp({host}:{port})/{name}?multiStatements=true"
    )


if __name__ == "__main__":
    main()
