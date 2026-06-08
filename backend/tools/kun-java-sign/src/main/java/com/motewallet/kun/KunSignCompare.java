package com.motewallet.kun;

import java.nio.file.Path;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.LinkedHashMap;
import java.util.Map;

/** Compare Java sign output with a captured Go request log. */
public final class KunSignCompare {

    private static final String EXPECTED_SIGN_STRING =
            "customerNo=10013098&email=jinfenggt2@gmail.com&requestNo=MW202606042325044aba2468&timestamp=1780586704766";
    private static final String EXPECTED_GO_SIGN =
            "BFL0JJnEmR8dmH4hJPvIth754x08576+XkVgxzJm3LDV+/D3lymWSQK9sRM2Y+sm7UcjOoQNWnBmQ1o1khYdA3yhQ3A7d2s/vfrDXL5at/3wNeOnYezuk268W3Ciy8/03X+xP40+cd4H3PNh4hJHt4VSCx26WXXY3bJpHCtu2fvM0po4f3rIuvHYpR8Zuc2Q1bJoTODW1beFiZTAcJ2sq4hULzM9EOKMry8sHennNYhS08tELgi3kpVl8Mv/n0Oi7sqrn3smqlCbCG2cbkiUbxBw/j5GJsXGSM5Y4uZMpa0Hh5a9+USKnKSfti0WF9vQSyN7PB5NcNUG+CexBnih3w==";

    public static void main(String[] args) throws Exception {
        Path keyPath = args.length > 0
                ? Path.of(args[0])
                : Path.of("../../.kun-keys/app_private_key_pkcs8.pem");

        String customerNo = "10013098";
        String timestamp = "1780586704766";

        Map<String, Object> body = new LinkedHashMap<>();
        body.put("email", "jinfenggt2@gmail.com");
        body.put("requestNo", "MW202606042325044aba2468");

        PrivateKey privateKey = KunSignUtil.loadPrivateKey(keyPath);
        KunSignUtil.SignResult result = KunSignUtil.signRequest(body, customerNo, timestamp, privateKey);

        boolean signStringMatch = EXPECTED_SIGN_STRING.equals(result.signString());
        boolean signMatch = EXPECTED_GO_SIGN.equals(result.signatureBase64());

        System.out.println("=== Java vs Go log comparison ===");
        System.out.println("private key: " + keyPath.toAbsolutePath());
        System.out.println();
        System.out.println("[signString]");
        System.out.println("  Go log:  " + EXPECTED_SIGN_STRING);
        System.out.println("  Java:    " + result.signString());
        System.out.println("  MATCH:   " + signStringMatch);
        System.out.println();
        System.out.println("[Sign Base64]");
        System.out.println("  Go log:  " + EXPECTED_GO_SIGN);
        System.out.println("  Java:    " + result.signatureBase64());
        System.out.println("  MATCH:   " + signMatch);
        System.out.println();

        Path pubPath = keyPath.getParent().resolve("app_public_key.pem");
        if (pubPath.toFile().exists()) {
            PublicKey pub = KunSignUtil.loadPublicKey(pubPath);
            boolean verifyGo = KunSignUtil.verify(EXPECTED_SIGN_STRING, EXPECTED_GO_SIGN, pub);
            boolean verifyJava = KunSignUtil.verify(result.signString(), result.signatureBase64(), pub);
            System.out.println("[local RSA verify with app_public_key.pem]");
            System.out.println("  Go Sign + signString:   " + verifyGo);
            System.out.println("  Java Sign + signString: " + verifyJava);
        }

        if (!signStringMatch || !signMatch) {
            System.exit(1);
        }
    }
}
