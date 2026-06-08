package com.motewallet.kun;

import com.google.gson.Gson;

import java.nio.file.Path;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Demo: sign the same payload as POST /rest/v2.0/customer/register
 *
 * Run from backend/tools/kun-java-sign:
 *   mvn -q compile exec:java \
 *     -Dexec.args="../../.kun-keys/app_private_key_pkcs8.pem 10013098"
 *
 * Or with explicit email / requestNo:
 *   mvn -q compile exec:java \
 *     -Dexec.args="../../.kun-keys/app_private_key_pkcs8.pem 10013098 jinfeng@example.com MW20260604120000abc"
 */
public final class KunSignDemo {

    public static void main(String[] args) throws Exception {
        if (args.length < 2) {
            System.err.println("Usage: KunSignDemo <private-key-pem-path> <customerNo> [email] [requestNo]");
            System.exit(1);
        }

        Path keyPath = Path.of(args[0]);
        String customerNo = args[1];
        String email = args.length > 2 ? args[2] : "test@example.com";
        String requestNo = args.length > 3 ? args[3] : "MW" + System.currentTimeMillis();

        String timestamp = String.valueOf(System.currentTimeMillis());

        Map<String, Object> body = new LinkedHashMap<>();
        body.put("email", email);
        body.put("requestNo", requestNo);

        PrivateKey privateKey = KunSignUtil.loadPrivateKey(keyPath);
        KunSignUtil.SignResult result = KunSignUtil.signRequest(body, customerNo, timestamp, privateKey);

        System.out.println("=== KUN Java Sign Demo ===");
        System.out.println("signString: " + result.signString());
        System.out.println("Sign (Base64): " + result.signatureBase64());
        System.out.println();
        System.out.println("Request headers:");
        System.out.println("  Customer-No: " + customerNo);
        System.out.println("  Timestamp: " + timestamp);
        System.out.println("  Sign: " + result.signatureBase64());
        System.out.println("Request body JSON: " + new Gson().toJson(body));

        Path pubPath = keyPath.getParent().resolve("app_public_key.pem");
        if (pubPath.toFile().exists()) {
            PublicKey publicKey = KunSignUtil.loadPublicKey(pubPath);
            boolean ok = KunSignUtil.verify(result.signString(), result.signatureBase64(), publicKey);
            System.out.println();
            System.out.println("Local verify with app_public_key.pem: " + ok);
        }
    }
}
