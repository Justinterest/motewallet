package com.motewallet.kun;

import java.nio.file.Path;
import java.security.PublicKey;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Simulate KUN platform signature verification for a captured request log.
 *
 * Run:
 *   mvn -q compile
 *   java -cp "target/classes:target/lib/*" com.motewallet.kun.KunVerifyDemo
 */
public final class KunVerifyDemo {

    // Captured from Go [kun api call] log
    private static final String SIGN_STRING =
            "customerNo=10013098&email=jinfenggt2@gmail.com&requestno=MW202606042325044aba2468&timestamp=1780586704766";
    private static final String SIGN =
            "XFBTGsQponSjBhLGc5iEM5hOVwaPWOdB63OZvOKMlJHJyMWZ7W1++DByfr3tZ6p8C+pUxlMuPNPu6hYn+GdKuyULx7W82ewhCPohhsRzSkG3/mQzeC8PzdL9lIvYHu8YdUtsTi3Wtltp/Mo0ui/E8BWQLV2pBTV+6OOKzQR5zOcn/XsaLYeiXlkgp6ePX7ewqEZy/Su0UeVPMg71Gw8zeCK+2+AbRsAG2ig6toJ/4yW6i9ziYMsMBp3jPmI5bU7m9tOOZ5hc0SBo0ceXJ+5ulzkzBdyfAG5/RRJFFJ9Y4WM0SA2Z242uPxWeCv+YW/eP75G6u0yww/7nJe2y0arXqQ==";

    /** Previous public key (before key rotation) — simulates KUN if old key still on file. */
    private static final String OLD_PUBLIC_KEY_B64 =
            "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjWzZzs9rMMYVVZtmq6s9713ZLVE3O7dWryuLT/fyBtL6EgSRalX+dNgLn2OlZlf6RdotHrcoYSBPQFprGJy5P3xzH+ZtirCEHOKINl4b26nXA1wrRS9UQLNRTh/Y4G0jjfYogDZdyrvW3x3FtTQ+r7eLlw9Dabz36UBZOvYvdpCz6dDS8CUb9YY6J9/p/w4ZFcmv/xNkzvHogl4kiKxSc+OuuXqUdkIMCO0irbAvbUQ7/002FzlbqB2Xa/Xrj5iVcWIzIYzmTYGMMY9F6HK8OnUhhufCFjvdK78WIgU6E2SsAGopPIM8EWEL3q5mwxoJ1Tw3eR9FWUtEYzi/fYP/zwIDAQAB";

    public static void main(String[] args) throws Exception {
        Path keysDir = args.length > 0 ? Path.of(args[0]) : Path.of("../../.kun-keys");

        PublicKey currentPub = KunSignUtil.loadPublicKey(keysDir.resolve("app_public_key.pem"));
        PublicKey oldPub = KunSignUtil.loadPublicKey(OLD_PUBLIC_KEY_B64);

        Map<String, Object> body = new LinkedHashMap<>();
        body.put("email", "jinfenggt2@gmail.com");
        body.put("requestNo", "MW202606042325044aba2468");
        String customerNo = "10013098";
        String timestamp = "1780586704766";

        System.out.println("=== Simulate KUN verify (SHA256withRSA) ===\n");

        // 1) Direct: use sign_string from log
        System.out.println("--- [1] Verify prebuilt sign_string + Sign ---");
        printVerify("Current public key (matches .kun-keys private)",
                KunVerifyUtil.verifySignString(SIGN_STRING, SIGN, currentPub));
        printVerify("OLD public key (simulates wrong key on KUN console)",
                KunVerifyUtil.verifySignString(SIGN_STRING, SIGN, oldPub));

        // 2) Rebuild canonical from body + headers (KUN server path)
        System.out.println("\n--- [2] Rebuild sign_string from request_body + Customer-No + Timestamp ---");
        KunVerifyUtil.VerifyResult rebuilt = KunVerifyUtil.verifyOutboundRequest(
                body, customerNo, timestamp, SIGN, currentPub);
        System.out.println("signString rebuilt: " + rebuilt.signString());
        System.out.println("equals log sign_string: " + SIGN_STRING.equals(rebuilt.signString()));
        printVerify("Current public key", rebuilt);

        KunVerifyUtil.VerifyResult rebuiltOld = KunVerifyUtil.verifyOutboundRequest(
                body, customerNo, timestamp, SIGN, oldPub);
        printVerify("OLD public key", rebuiltOld);

        // 3) Summary
        System.out.println("\n=== Conclusion ===");
        boolean kunWouldAccept = rebuilt.ok();
        System.out.println("If KUN has CURRENT public key (app_public_key.pem): "
                + (kunWouldAccept ? "PASS (100033 should NOT be signature mismatch)" : "FAIL"));
        System.out.println("If KUN still has OLD public key: "
                + (rebuiltOld.ok() ? "PASS" : "FAIL → explains code 100033"));
    }

    private static void printVerify(String label, KunVerifyUtil.VerifyResult r) {
        System.out.println(label + ": " + (r.ok() ? "PASS ✓" : "FAIL ✗") + " — " + r.message());
    }
}
