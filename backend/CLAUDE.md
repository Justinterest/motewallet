# Backend — Golang

## 技术栈
- Go 1.22+
- Gin（Web 框架）
- GORM v2（ORM）
- golang-jwt/jwt（认证）
- go-playground/validator（校验）
- zap / slog（日志）
- viper（配置管理）
- wire（依赖注入）

## 目录约定
- `cmd/server/` — 入口 main.go
- `internal/handler/` — HTTP Handler（Controller）
- `internal/service/` — 业务逻辑层
- `internal/repository/` — 数据访问层
- `internal/model/` — GORM 数据模型
- `internal/dto/` — 请求/响应 DTO
- `internal/middleware/` — 中间件
- `internal/config/` — 配置
- `internal/pkg/` — 内部公共包

## 架构分层
```
Handler → Service → Repository → Model
```
- Handler 只做参数绑定和校验
- Service 处理业务规则和事务
- Repository 封装数据库操作

## 开发命令
```bash
go run cmd/server/main.go    # 启动服务
go test ./...                # 运行测试
go vet ./...                 # 静态检查
make swagger                 # 生成 Swagger 文档
```

## 注意事项
- 金额使用 `shopspring/decimal`，与数据库 `DECIMAL(28,8)` 对齐，API JSON 传字符串（如 `"100.50000000"`）
- 统一响应格式：`{"code": 0, "message": "success", "data": {}}`
- API 路由前缀：`/api/v1/`
- 密码使用 bcrypt 哈希
- 配置通过环境变量注入
- 数据库 Schema 由 Database Agent 管理，不自行修改 migration
- API 契约见 `docs/api/openapi.yaml`
