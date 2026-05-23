package response

type WalletBalanceResp struct {
	AccountType      string `json:"account_type"`
	Currency         string `json:"currency"`
	Balance          string `json:"balance"`
	FrozenBalance    string `json:"frozen_balance"`
	AvailableBalance string `json:"available_balance"`
}

type WalletBalancesResp struct {
	Wallets []WalletBalanceResp `json:"wallets"`
}
