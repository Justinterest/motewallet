package kun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	kundto "motewallet/internal/pkg/kun/dto"
)

const subMerchantGoSignString = `customerNo=10038701&directorInfo=[{"gender":"Male","idCard":"111","nameEN":"111","country":"HK","nameCHS":"111","surname":"11","authType":"00000001003","birthday":"2026-06-05","surnameCHS":"111","certificateTerm":["2026-06-27","2026-07-04"],"residenceAddress":"111111","residenceCountry":"HK","idHolding":[{"path":"https://img.kun.global/kun-api/e7593fd9ba0948c69a09967ba5dd3623.png"}],"verificationType":"idHolding"}]&enterpriseInfo={"incorporationCertificate":[{"path":"https://img.kun.global/kun-api/bfc20d02130f4c2f8b41848e52ffbea2.png"}],"incorporationCertificateNo":"121231","establishTime":"2026-05-01","enterpriseEN":"1212","enterpriseNameCHS":"无","registerRegion":"HK","registerAddress":"111","businessRegistration":[{"path":"https://img.kun.global/kun-api/ad1b76894a8042e9a4ca1b00b2276b1b.png"}],"businessRegistrationNo":"12121","phone":"1111","isChangeEnterpriseNameInFiveYears":"No","enterpriseType":"Private limited company","mainBusinessAddress":"111","industry":"Goods Trade","subIndustry":"111","initialFundingSource":"Business income","wealthSource":"Business income","continuousFundingSource":"Business income","salesVolumeLastyear":"HKD 0-2,500,000","employeeNum":"<10","openAccountPurpose":"Business operation","associationRules":[{"path":"https://img.kun.global/kun-api/3ba6898fe31c44e4b7832bd985c037bb.png"}],"authenticMaterials":[{"path":"https://img.kun.global/kun-api/0167e8b5a7cc454f84da6ecfc0b8ce89.png"}],"managerCountry":"HK","managerAuthType":"00000001003","managerVerificationType":"idHolding","managerIdHolding":[{"path":"https://img.kun.global/kun-api/7a62932b882744668aef60e2e107faf2.png"},{"path":"https://img.kun.global/kun-api/fe819213464648468796a8ba5b1a6485.png"},{"path":"https://img.kun.global/kun-api/64b0b258fbfa4195b364078791050254.png"}],"managerCertificateTerm":["2026-06-12","2026-06-28"],"managerSurnameCHS":"121","managerNameCHS":"111","managerSurname":"111","managerNameEN":"111","managerBirthday":"2026-06-13","managerGender":"Female","managerIdCard":"1111","managerResidenceCountry":"HK","managerResidenceAddress":"12121","middleTierShareholders":"No","nnc1":[{"path":"https://img.kun.global/kun-api/b71af5f61263429fb80cb3e397696bc4.png"}]}&requestNo=MW202606081547141b8a5709&shareholdersInfo=[{"gender":"Male","idCard":"11","nameEN":"111","country":"HK","nameCHS":"11","surname":"11","authType":"00000001003","birthday":"2026-06-19","surnameCHS":"111","certificateTerm":["2026-06-06","2026-06-19"],"residenceAddress":"11111","residenceCountry":"HK","shareholdingRatio":"11","idHolding":[{"path":"https://img.kun.global/kun-api/5bd12dda340b48629ca1eb03bd6441f8.png"}],"verificationType":"idHolding"}]&timestamp=1780904848142`

func loadSubMerchantFixture(t *testing.T) kundto.SubMerchantRegisterReq {
	t.Helper()
	bodyPath := filepath.Join("..", "..", "..", "tools", "kun-java-sign", "fixtures", "sub-merchant-register-body.json")
	raw, err := os.ReadFile(bodyPath)
	if err != nil {
		t.Skipf("fixture not found: %v", err)
	}
	var req kundto.SubMerchantRegisterReq
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatal(err)
	}
	return req
}

func TestSubMerchantRegisterSignString_fromStructToMap(t *testing.T) {
	req := loadSubMerchantFixture(t)
	params, err := StructToMap(req)
	if err != nil {
		t.Fatal(err)
	}

	got := buildCanonicalKVString(params, map[string]string{
		"customerNo": "10038701",
		"timestamp":  "1780904848142",
	})
	if got != subMerchantGoSignString {
		t.Fatalf("sign string mismatch\nwant len=%d\ngot len=%d\nfirst diff: %s",
			len(subMerchantGoSignString), len(got), diffSnippet(subMerchantGoSignString, got))
	}
}

func TestSubMerchantRegister_goMarshalBodyOrder(t *testing.T) {
	req := loadSubMerchantFixture(t)
	body, err := MarshalRequestBody(req)
	if err != nil {
		t.Fatal(err)
	}
	outPath := filepath.Join("..", "..", "..", "tools", "kun-java-sign", "fixtures", "sub-merchant-register-body-go-marshal.json")
	if err := os.WriteFile(outPath, body, 0o644); err != nil {
		t.Fatal(err)
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, `"requestNo"`) {
		t.Fatalf("expected requestNo in body: %s", bodyStr[:min(200, len(bodyStr))])
	}
	genderIdx := strings.Index(bodyStr, `"gender"`)
	authIdx := strings.Index(bodyStr, `"authType"`)
	if genderIdx < 0 || authIdx < 0 || genderIdx > authIdx {
		t.Fatalf("expected gender before authType in Go marshaled body\ngender@%d auth@%d", genderIdx, authIdx)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestSubMerchantRegisterSignString_jsonUnmarshalMapDiffers(t *testing.T) {
	bodyPath := filepath.Join("..", "..", "..", "tools", "kun-java-sign", "fixtures", "sub-merchant-register-body.json")
	raw, err := os.ReadFile(bodyPath)
	if err != nil {
		t.Skipf("fixture not found: %v", err)
	}
	var req map[string]interface{}
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatal(err)
	}
	got := buildCanonicalKVString(req, map[string]string{
		"customerNo": "10038701",
		"timestamp":  "1780904848142",
	})
	if got == subMerchantGoSignString {
		t.Fatal("expected json.Unmarshal map path to differ from StructToMap production path")
	}
}

func diffSnippet(want, got string) string {
	max := len(want)
	if len(got) < max {
		max = len(got)
	}
	for i := 0; i < max; i++ {
		if want[i] != got[i] {
			from := i - 40
			if from < 0 {
				from = 0
			}
			toWant := i + 40
			if toWant > len(want) {
				toWant = len(want)
			}
			toGot := i + 40
			if toGot > len(got) {
				toGot = len(got)
			}
			return "want ..." + want[from:toWant] + "... vs got ..." + got[from:toGot] + "..."
		}
	}
	return "equal prefix, length differs"
}
