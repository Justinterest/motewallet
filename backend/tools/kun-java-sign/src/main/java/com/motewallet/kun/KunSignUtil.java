package com.motewallet.kun;

import com.google.gson.Gson;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.reflect.TypeToken;
import org.bouncycastle.asn1.DERNull;
import org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers;
import org.bouncycastle.asn1.pkcs.PrivateKeyInfo;
import org.bouncycastle.asn1.pkcs.RSAPrivateKey;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.openssl.PEMKeyPair;
import org.bouncycastle.openssl.PEMParser;
import org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter;

import java.io.IOException;
import java.io.StringReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Signature;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.TreeMap;

/**
 * KUN OpenAPI request signing — aligned with backend/internal/pkg/kun/signer.go
 * and https://opendocs.kun.global/docs/api-integration/authentication-and-authorization-mechanism
 */
public final class KunSignUtil {

    private static final Gson GSON = new Gson();

    private KunSignUtil() {}

    public record SignResult(String signString, String signatureBase64) {}

    /** Sign business params + Customer-No + Timestamp. */
    public static SignResult signRequest(
            Map<String, Object> bizParams,
            String customerNo,
            String timestamp,
            PrivateKey privateKey) throws Exception {
        String signString = buildCanonicalKVString(bizParams, Map.of(
                "customerNo", customerNo,
                "timestamp", timestamp
        ));
        String signature = sign(signString, privateKey);
        return new SignResult(signString, signature);
    }

    /** Parse JSON body object to top-level map (same idea as Go StructToMap). */
    public static Map<String, Object> jsonBodyToMap(String jsonBody) {
        return GSON.fromJson(jsonBody, new TypeToken<Map<String, Object>>() {}.getType());
    }

    public static Map<String, Object> objectToMap(Object dto) {
        JsonElement tree = GSON.toJsonTree(dto);
        return GSON.fromJson(tree, new TypeToken<Map<String, Object>>() {}.getType());
    }

    /**
     * Build canonical string: sorted key=value joined by &.
     * Business keys keep JSON body casing (e.g. requestNo, directorInfo).
     * System keys keep customerNo / timestamp casing.
     */
    public static String buildCanonicalKVString(
            Map<String, Object> biz,
            Map<String, String> extraPlain) {
        Map<String, String> params = new LinkedHashMap<>();

        if (biz != null) {
            for (Map.Entry<String, Object> entry : biz.entrySet()) {
                if (entry.getValue() == null) {
                    continue;
                }
                String key = entry.getKey().trim();
                if (key.isEmpty() || "sign".equalsIgnoreCase(key)) {
                    continue;
                }
                params.put(key, valueToString(entry.getValue()));
            }
        }

        if (extraPlain != null) {
            for (Map.Entry<String, String> entry : extraPlain.entrySet()) {
                String key = entry.getKey().trim();
                if (key.isEmpty() || "sign".equalsIgnoreCase(key)) {
                    continue;
                }
                params.put(key, entry.getValue());
            }
        }

        TreeMap<String, String> sorted = new TreeMap<>(params);
        StringBuilder sb = new StringBuilder();
        for (Map.Entry<String, String> entry : sorted.entrySet()) {
            if (sb.length() > 0) {
                sb.append('&');
            }
            sb.append(entry.getKey()).append('=').append(entry.getValue());
        }
        return sb.toString();
    }

    public static String valueToString(Object value) {
        if (value == null) {
            return "";
        }
        if (value instanceof String s) {
            return s;
        }
        if (value instanceof Boolean b) {
            return b ? "true" : "false";
        }
        if (value instanceof Number n) {
            double d = n.doubleValue();
            if (d == Math.rint(d) && !Double.isInfinite(d)) {
                return String.valueOf((long) d);
            }
            return n.toString();
        }
        return GSON.toJson(value);
    }

    /** SHA256withRSA + Base64 (same as Java Signature SHA256withRSA). */
    public static String sign(String signString, PrivateKey privateKey) throws Exception {
        Signature signature = Signature.getInstance("SHA256withRSA");
        signature.initSign(privateKey);
        signature.update(signString.getBytes(StandardCharsets.UTF_8));
        return Base64.getEncoder().encodeToString(signature.sign());
    }

    public static boolean verify(String signString, String signatureBase64, PublicKey publicKey) throws Exception {
        Signature signature = Signature.getInstance("SHA256withRSA");
        signature.initVerify(publicKey);
        signature.update(signString.getBytes(StandardCharsets.UTF_8));
        return signature.verify(Base64.getDecoder().decode(signatureBase64));
    }

    /** Load private key from PEM file, PEM string, or raw Base64 DER (PKCS#8 / PKCS#1). */
    public static PrivateKey loadPrivateKey(String keyMaterial) throws Exception {
        String trimmed = keyMaterial.trim().replace("\\n", "\n");
        if (trimmed.contains("-----BEGIN")) {
            return loadPrivateKeyFromPem(trimmed);
        }
        byte[] der = Base64.getDecoder().decode(trimmed);
        return loadPrivateKeyFromDer(der);
    }

    public static PrivateKey loadPrivateKey(Path pemPath) throws Exception {
        return loadPrivateKey(Files.readString(pemPath));
    }

    public static PublicKey loadPublicKey(String keyMaterial) throws Exception {
        String trimmed = keyMaterial.trim().replace("\\n", "\n");
        if (trimmed.contains("-----BEGIN")) {
            try (PEMParser parser = new PEMParser(new StringReader(trimmed))) {
                Object obj = parser.readObject();
                JcaPEMKeyConverter converter = new JcaPEMKeyConverter();
                if (obj instanceof org.bouncycastle.cert.X509CertificateHolder cert) {
                    return converter.getPublicKey(cert.getSubjectPublicKeyInfo());
                }
                if (obj instanceof org.bouncycastle.asn1.x509.SubjectPublicKeyInfo info) {
                    return converter.getPublicKey(info);
                }
                throw new IllegalArgumentException("unsupported public PEM object");
            }
        }
        byte[] der = Base64.getDecoder().decode(trimmed);
        return KeyFactory.getInstance("RSA").generatePublic(new X509EncodedKeySpec(der));
    }

    public static PublicKey loadPublicKey(Path pemPath) throws Exception {
        return loadPublicKey(Files.readString(pemPath));
    }

    private static PrivateKey loadPrivateKeyFromPem(String pem) throws IOException {
        try (PEMParser parser = new PEMParser(new StringReader(pem))) {
            Object obj = parser.readObject();
            JcaPEMKeyConverter converter = new JcaPEMKeyConverter();
            if (obj instanceof PEMKeyPair pair) {
                return converter.getPrivateKey(pair.getPrivateKeyInfo());
            }
            if (obj instanceof PrivateKeyInfo info) {
                return converter.getPrivateKey(info);
            }
            if (obj instanceof PrivateKey pk) {
                return pk;
            }
            throw new IllegalArgumentException("unsupported private PEM object: " + obj.getClass());
        }
    }

    private static PrivateKey loadPrivateKeyFromDer(byte[] der) throws Exception {
        KeyFactory kf = KeyFactory.getInstance("RSA");
        try {
            return kf.generatePrivate(new PKCS8EncodedKeySpec(der));
        } catch (Exception ignored) {
            // PKCS#1 RSAPrivateKey DER (e.g. .b64.txt from openssl pkcs8 -outform DER)
            RSAPrivateKey rsa = RSAPrivateKey.getInstance(der);
            PrivateKeyInfo info = new PrivateKeyInfo(
                    new AlgorithmIdentifier(PKCSObjectIdentifiers.rsaEncryption, DERNull.INSTANCE),
                    rsa);
            return new JcaPEMKeyConverter().getPrivateKey(info);
        }
    }
}
