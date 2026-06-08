package com.motewallet.kun;

import java.nio.file.Files;
import java.nio.file.Path;
import java.security.PrivateKey;
import java.security.PublicKey;

/**
 * Compare Java sign output with a captured Go KUN API log for sub-merchant register.
 *
 * Run from backend/tools/kun-java-sign:
 *   mvn -q compile
 *   java -cp "target/classes:target/lib/*" com.motewallet.kun.KunSignSubMerchantCompare \
 *     ../../.kun-keys/app_private_key_pkcs8.pem
 */
public final class KunSignSubMerchantCompare {

    private static final String CUSTOMER_NO = "10038701";
    private static final String TIMESTAMP = "1780904848142";
    private static final String GO_SIGN_STRING = """
            customerNo=10038701&directorInfo=[{"gender":"Male","idCard":"111","nameEN":"111","country":"HK","nameCHS":"111","surname":"11","authType":"00000001003","birthday":"2026-06-05","surnameCHS":"111","certificateTerm":["2026-06-27","2026-07-04"],"residenceAddress":"111111","residenceCountry":"HK","idHolding":[{"path":"https://img.kun.global/kun-api/e7593fd9ba0948c69a09967ba5dd3623.png"}],"verificationType":"idHolding"}]&enterpriseInfo={"incorporationCertificate":[{"path":"https://img.kun.global/kun-api/bfc20d02130f4c2f8b41848e52ffbea2.png"}],"incorporationCertificateNo":"121231","establishTime":"2026-05-01","enterpriseEN":"1212","enterpriseNameCHS":"无","registerRegion":"HK","registerAddress":"111","businessRegistration":[{"path":"https://img.kun.global/kun-api/ad1b76894a8042e9a4ca1b00b2276b1b.png"}],"businessRegistrationNo":"12121","phone":"1111","isChangeEnterpriseNameInFiveYears":"No","enterpriseType":"Private limited company","mainBusinessAddress":"111","industry":"Goods Trade","subIndustry":"111","initialFundingSource":"Business income","wealthSource":"Business income","continuousFundingSource":"Business income","salesVolumeLastyear":"HKD 0-2,500,000","employeeNum":"<10","openAccountPurpose":"Business operation","associationRules":[{"path":"https://img.kun.global/kun-api/3ba6898fe31c44e4b7832bd985c037bb.png"}],"authenticMaterials":[{"path":"https://img.kun.global/kun-api/0167e8b5a7cc454f84da6ecfc0b8ce89.png"}],"managerCountry":"HK","managerAuthType":"00000001003","managerVerificationType":"idHolding","managerIdHolding":[{"path":"https://img.kun.global/kun-api/7a62932b882744668aef60e2e107faf2.png"},{"path":"https://img.kun.global/kun-api/fe819213464648468796a8ba5b1a6485.png"},{"path":"https://img.kun.global/kun-api/64b0b258fbfa4195b364078791050254.png"}],"managerCertificateTerm":["2026-06-12","2026-06-28"],"managerSurnameCHS":"121","managerNameCHS":"111","managerSurname":"111","managerNameEN":"111","managerBirthday":"2026-06-13","managerGender":"Female","managerIdCard":"1111","managerResidenceCountry":"HK","managerResidenceAddress":"12121","middleTierShareholders":"No","nnc1":[{"path":"https://img.kun.global/kun-api/b71af5f61263429fb80cb3e397696bc4.png"}]}&requestNo=MW202606081547141b8a5709&shareholdersInfo=[{"gender":"Male","idCard":"11","nameEN":"111","country":"HK","nameCHS":"11","surname":"11","authType":"00000001003","birthday":"2026-06-19","surnameCHS":"111","certificateTerm":["2026-06-06","2026-06-19"],"residenceAddress":"11111","residenceCountry":"HK","shareholdingRatio":"11","idHolding":[{"path":"https://img.kun.global/kun-api/5bd12dda340b48629ca1eb03bd6441f8.png"}],"verificationType":"idHolding"}]&timestamp=1780904848142""";
    private static final String GO_SIGN_BASE64 =
            "e7SO2i9/p/4M6f/DGkPjz1kl3sO8mDLp1Xa6sNncjTS235zGB0/driw2PqHbx1e6maJMyuKjHawuHFLbOdO0ZuGAvmNtULO3mLOL4RO1CAQsVpjpbRru9TC2dtqvWiFL3CMtHW3576cgLCHEsp13U81DtMYS3uId5/kBV69lgwuNMfY3l7mE+1FxQ4S1EQitQJrMFfq/zkR2otR4X4qBDAeR6ITBlap53tFptzIcoIddcTuWLJ32h66wnuDYQvSJMvR5pBvCzTLvh7YWgYbfRQhDnvlt4/RjHFo0hKwsO4YD0WczIAe1rg337pH1z9xYdDZhBwT10YRKscz+S3yc7Q==";

    public static void main(String[] args) throws Exception {
        Path keyPath = args.length > 0
                ? Path.of(args[0])
                : Path.of("../../.kun-keys/app_private_key_pkcs8.pem");

        PrivateKey privateKey = KunSignUtil.loadPrivateKey(keyPath);
        Path pubPath = keyPath.getParent().resolve("app_public_key.pem");
        PublicKey pub = pubPath.toFile().exists() ? KunSignUtil.loadPublicKey(pubPath) : null;

        System.out.println("=== Java vs Go subMerchant/register log ===");
        System.out.println("key: " + keyPath.toAbsolutePath());
        System.out.println();

        Path goBody = Path.of("fixtures/sub-merchant-register-body-go-marshal.json");
        if (!goBody.toFile().exists()) {
            System.err.println("missing fixture: run Go test TestSubMerchantRegister_goMarshalBodyOrder first");
            System.exit(1);
        }

        boolean matched = compareCase("go-http-body-json", goBody, privateKey, pub);
        if (!matched) {
            System.exit(1);
        }
    }

    private static boolean compareCase(String label, Path bodyPath, PrivateKey privateKey, PublicKey pub) throws Exception {
        String jsonBody = Files.readString(bodyPath);
        var body = KunSignUtil.jsonBodyToMap(jsonBody);
        KunSignUtil.SignResult result = KunSignUtil.signRequest(body, CUSTOMER_NO, TIMESTAMP, privateKey);

        boolean signStringMatch = GO_SIGN_STRING.equals(result.signString());
        boolean signMatch = GO_SIGN_BASE64.equals(result.signatureBase64());

        System.out.println("--- " + label + " ---");
        System.out.println("body: " + bodyPath.toAbsolutePath());
        System.out.println("[signString MATCH] " + signStringMatch);
        if (!signStringMatch) {
            printFirstDiff(GO_SIGN_STRING, result.signString());
        }
        System.out.println("[Sign Base64 MATCH] " + signMatch);
        if (pub != null) {
            System.out.println("[verify Go signString + Go Sign] " + KunSignUtil.verify(GO_SIGN_STRING, GO_SIGN_BASE64, pub));
            System.out.println("[verify Java signString + Go Sign] " + KunSignUtil.verify(result.signString(), GO_SIGN_BASE64, pub));
        }
        System.out.println();
        return signStringMatch && signMatch;
    }

    private static void printFirstDiff(String go, String java) {
        int max = Math.min(go.length(), java.length());
        for (int i = 0; i < max; i++) {
            if (go.charAt(i) != java.charAt(i)) {
                int from = Math.max(0, i - 40);
                int toGo = Math.min(go.length(), i + 40);
                int toJava = Math.min(java.length(), i + 40);
                System.out.println("  first diff at index " + i);
                System.out.println("  Go:   ..." + go.substring(from, toGo) + "...");
                System.out.println("  Java: ..." + java.substring(from, toJava) + "...");
                return;
            }
        }
        if (go.length() != java.length()) {
            System.out.println("  strings equal up to shorter length; lengths differ");
        }
    }
}
