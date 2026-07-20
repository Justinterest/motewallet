# Docker 部署

## 架构

| 服务 | 镜像 | 容器端口 | 宿主机默认绑定 |
|------|------|----------|----------------|
| backend | `motewallet-backend` | 8080 | `127.0.0.1:8080`（不对外） |
| frontend | `motewallet-frontend` | 3000 | `127.0.0.1:3000` |
| admin | `motewallet-admin` | 3001 | `127.0.0.1:3001` |

```
浏览器 → nginx (443)
           ├─ /api/*  → 127.0.0.1:8080  (backend)
           └─ /*      → 127.0.0.1:3000/3001  (frontend/admin)
```

API **不单独暴露域名/端口**，只经 nginx 同域转发。前端构建时 `NEXT_PUBLIC_API_URL` 默认为空（同源请求 `/api/v1/...`）。

流程：**本机 build（linux/amd64）→ export tar.gz → scp 上传 → 服务器 docker load → compose 重启**。

现有 `motewallet.com` / `ops.motewallet.com`（Java）不受影响。

## 首次准备

### 1. 本机

```bash
cp deploy/.env.example deploy/.env
# 编辑 deploy/.env：数据库、JWT、KUN、S3 等

chmod +x deploy/scripts/*.sh
```

### 2. 服务器（`motewallet-prod`）

首次 deploy 会自动创建 `/home/ecs-user/motewallet-withdrawal` 并上传 `.env` 模板。之后务必：

```bash
ssh motewallet-prod
vi ~/motewallet-withdrawal/.env   # 填生产配置
```

数据库需可连通（当前生产机未跑本地 MySQL，一般用 RDS）。

### 3. Nginx 子域

参考 `deploy/nginx/motewallet-withdrawal.conf.example`：

- `wallet.motewallet.com` → frontend + `/api/` → backend
- `portal.motewallet.com` → admin + `/api/` → backend

KUN webhook：`https://wallet.motewallet.com/api/v1/webhook/kun`

上传并启用：

```bash
scp deploy/nginx/motewallet-withdrawal.conf.example \
  motewallet-prod:/tmp/motewallet-withdrawal.conf

ssh motewallet-prod 'sudo mv /tmp/motewallet-withdrawal.conf /etc/nginx/conf.d/motewallet-withdrawal.conf \
  && sudo nginx -t && sudo systemctl reload nginx'
```

## 数据库迁移（线上）

Schema 变更在 `database/migrations/`（golang-migrate）。**不要**把 seed 打到生产。

推荐顺序：**先 migrate，再发 backend**。

```bash
# 查看线上当前版本
./deploy/scripts/migrate.sh version

# 首次：建库 + 跑完全部 migration
./deploy/scripts/migrate.sh create-db
./deploy/scripts/migrate.sh up

# 日常：有新 migration 时
./deploy/scripts/migrate.sh up

# 回滚一步（慎用）
./deploy/scripts/migrate.sh down 1
```

脚本会：
1. rsync `database/migrations/` 到服务器
2. **每次同步** 本地 `deploy/.env` → 服务器（改密码后无需手动 scp）
3. 用 Docker 跑 `migrate/migrate` 连库执行

本机直连 DB（需网络可达）：

```bash
MIGRATE_LOCAL=1 ./deploy/scripts/migrate.sh version
```

## 日常发布（可单独发）

```bash
# 只发后端
./deploy/scripts/deploy.sh backend

# 只发商户端 / 管理端（默认同源，无需设 API URL）
./deploy/scripts/deploy.sh frontend
./deploy/scripts/deploy.sh admin

# 三个一起
./deploy/scripts/deploy.sh all
```

分步执行：

```bash
./deploy/scripts/build.sh backend
./deploy/scripts/export.sh backend
./deploy/scripts/deploy.sh backend --skip-build --skip-export
```

指定版本标签：

```bash
TAG=1.0.0 ./deploy/scripts/deploy.sh backend
```

## 服务器手动操作

```bash
ssh motewallet-prod
cd ~/motewallet-withdrawal

# 加载并重启某一个服务
TAG=20260719-abc1234 ./remote-up.sh backend

# 查看状态 / 日志
docker compose ps
docker compose logs -f --tail=100 backend
```

## 注意

- 本机是 Apple Silicon 时脚本固定 `--platform linux/amd64`，与生产机架构一致。
- backend 端口只绑 `127.0.0.1`，外网无法直连。
- 镜像包在 `deploy/dist/`，已加入 gitignore，勿提交。
- 服务器内存约 3.4G，且已有 Java 容器；上线后留意 `free -h`。
