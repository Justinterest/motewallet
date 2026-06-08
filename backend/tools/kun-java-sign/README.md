# KUN Java 签名工具

与 `backend/internal/pkg/kun/signer.go` 规则一致，便于和 Go 服务、KUN 文档对照调试。

参考文档：

- [鉴权认证机制](https://opendocs.kun.global/docs/api-integration/authentication-and-authorization-mechanism)
- [制备商户密钥](https://opendocs.kun.global/docs/secret-key-configuration/merchant-key-preparation)

## 规则摘要

1. 业务参数：JSON body **顶级字段**（**保持 JSON 原始 key 大小写**，如 `requestNo`、`directorInfo`）
2. 系统参数：Header `Customer-No` → `customerNo`，`Timestamp` → `timestamp`（保持驼峰）
3. 跳过 `null` 与 `sign`
4. 按 key ASCII 排序，拼接 `key=value&...`（value **不** URL encode，与文档 Java 示例一致）
5. `SHA256withRSA` + Base64 → 请求头 `Sign`

## 构建与运行

```bash
cd backend/tools/kun-java-sign

# 编译
mvn -q compile

# 示例：用 PKCS#8 私钥签注册接口
mvn -q exec:java -Dexec.args="../../.kun-keys/app_private_key_pkcs8.pem 10013098 test@example.com MW20260604120000abc"
```

输出 `signString` 与 `Sign`，可与 Go 日志里的 `sign_string` 对比。

## 在业务代码中使用

```java
import com.motewallet.kun.KunSignUtil;
import java.security.PrivateKey;
import java.util.Map;

PrivateKey key = KunSignUtil.loadPrivateKey(Path.of("app_private_key_pkcs8.pem"));

Map<String, Object> body = KunSignUtil.jsonBodyToMap("""
    {"email":"a@b.com","requestNo":"MW001"}
    """);

String timestamp = String.valueOf(System.currentTimeMillis());
KunSignUtil.SignResult r = KunSignUtil.signRequest(body, "10013098", timestamp, key);

// HTTP headers
// Customer-No: 10013098
// Timestamp: {timestamp}
// Sign: {r.signatureBase64()}
```

## 私钥格式

| 格式 | Java 加载 |
|------|-----------|
| `app_private_key_pkcs8.pem` | 推荐，`KunSignUtil.loadPrivateKey(path)` |
| `.env` 一行 Base64（`app_private_key_pkcs8.b64.txt`） | `loadPrivateKey(base64String)`，支持 PKCS#8 / PKCS#1 DER |

## 与 Go 对齐验证

同一组 `email`、`requestNo`、`customerNo`、`timestamp`、同一私钥时，Java 与 Go 的 `signString` 和 `Sign` 应完全一致。

## 模拟 KUN 验签（服务端）

```bash
mvn -q compile
java -cp "target/classes:target/lib/*" com.motewallet.kun.KunVerifyDemo
```

`KunVerifyUtil` 提供：

- `verifyOutboundRequest` — 从 body + `Customer-No` + `Timestamp` 重建待签串后验签（模拟 KUN 收商户请求）
- `verifySignString` — 直接验已有 `sign_string`
- `verifyWebhookOrResponse` — Webhook/响应 `data` + `timestamp` 验签（对齐 Go `VerifyWebhookSignature`）
