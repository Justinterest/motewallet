package com.motewallet.kun;

import java.security.PublicKey;
import java.util.Map;

/**
 * KUN signature verification — aligned with backend/internal/pkg/kun/verifier.go
 * and https://opendocs.kun.global/docs/api-integration/authentication-and-authorization-mechanism
 */
public final class KunVerifyUtil {

    private KunVerifyUtil() {}

    public record VerifyResult(boolean ok, String signString, String message) {}

    /**
     * Verify outbound API request (merchant → KUN): rebuild canonical string from body + headers, then SHA256withRSA.
     * This simulates what KUN server does when it receives your POST.
     */
    public static VerifyResult verifyOutboundRequest(
            Map<String, Object> requestBody,
            String customerNo,
            String timestamp,
            String signatureBase64,
            PublicKey kunPlatformPublicKey) throws Exception {
        String signString = KunSignUtil.buildCanonicalKVString(requestBody, Map.of(
                "customerNo", customerNo,
                "timestamp", timestamp
        ));
        boolean ok = KunSignUtil.verify(signString, signatureBase64, kunPlatformPublicKey);
        String msg = ok ? "signature valid" : "signature invalid (public key mismatch or wrong canonical string)";
        return new VerifyResult(ok, signString, msg);
    }

    /**
     * Verify when canonical string is already known (e.g. from Go log sign_string).
     */
    public static VerifyResult verifySignString(
            String signString,
            String signatureBase64,
            PublicKey publicKey) throws Exception {
        boolean ok = KunSignUtil.verify(signString, signatureBase64, publicKey);
        String msg = ok ? "signature valid" : "signature invalid";
        return new VerifyResult(ok, signString, msg);
    }

    /**
     * Verify webhook / response (KUN → merchant): data fields + timestamp.
     * Same as Go VerifyWebhookSignature.
     */
    public static VerifyResult verifyWebhookOrResponse(
            Map<String, Object> data,
            String timestamp,
            String signatureBase64,
            PublicKey kunPublicKey) throws Exception {
        String signString = KunSignUtil.buildCanonicalKVString(data, Map.of("timestamp", timestamp));
        boolean ok = KunSignUtil.verify(signString, signatureBase64, kunPublicKey);
        String msg = ok ? "webhook/response signature valid" : "webhook/response signature invalid";
        return new VerifyResult(ok, signString, msg);
    }
}
