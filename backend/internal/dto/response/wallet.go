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

type WalletLedgerEntryResp struct {
	ID                  uint64  `json:"id"`
	AccountType         string  `json:"account_type"`
	Currency            string  `json:"currency"`
	EntryType           string  `json:"entry_type"`
	Amount              string  `json:"amount"`
	BalanceBefore       string  `json:"balance_before"`
	BalanceAfter        string  `json:"balance_after"`
	FrozenBefore        string  `json:"frozen_before"`
	FrozenAfter         string  `json:"frozen_after"`
	TransactionRecordID *uint64 `json:"transaction_record_id,omitempty"`
	PlatformOrderID     *string `json:"platform_order_id,omitempty"`
	BizType             *string `json:"biz_type,omitempty"`
	Remark              *string `json:"remark,omitempty"`
	CreatedAt           string  `json:"created_at"`
}

type WalletLedgerListResp struct {
	Entries  []WalletLedgerEntryResp `json:"entries"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}
