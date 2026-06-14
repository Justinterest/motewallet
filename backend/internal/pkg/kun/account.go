package kun

import "fmt"

const (
	AccountHK = "KUN_HK" // 香港站 = 资金账户
	AccountPL = "KUN_PL" // 波兰站 = 交易账户
)

func PlatformAccountToKUN(platformType string) (string, error) {
	switch platformType {
	case "FUNDING":
		return AccountHK, nil
	case "TRADING":
		return AccountPL, nil
	default:
		return "", fmt.Errorf("unsupported platform account type: %s", platformType)
	}
}
