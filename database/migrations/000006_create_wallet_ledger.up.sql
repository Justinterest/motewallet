-- ------------------------------------------------------------
-- wallet_ledger（钱包资金变动账本 — 只追加）
-- 每次余额/冻结变动写一行，用于对账与重建钱包状态
-- ------------------------------------------------------------
CREATE TABLE `wallet_ledger` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `wallet_id` BIGINT UNSIGNED NOT NULL,
  `account_type` VARCHAR(10) NOT NULL COMMENT 'FUNDING / TRADING',
  `currency` VARCHAR(10) NOT NULL,
  `entry_type` VARCHAR(20) NOT NULL COMMENT 'CREDIT / FREEZE / UNFREEZE / DEDUCT_FROZEN',
  `amount` DECIMAL(28,8) NOT NULL COMMENT '变动金额（正数）',
  `balance_before` DECIMAL(28,8) NOT NULL,
  `balance_after` DECIMAL(28,8) NOT NULL,
  `frozen_before` DECIMAL(28,8) NOT NULL,
  `frozen_after` DECIMAL(28,8) NOT NULL,
  `transaction_record_id` BIGINT UNSIGNED NULL COMMENT '关联业务流水',
  `biz_type` VARCHAR(20) NULL COMMENT 'DEPOSIT / WITHDRAWAL / EXCHANGE / TRANSFER',
  `remark` VARCHAR(500) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_wallet_ledger_merchant_currency_time` (`merchant_id`, `currency`, `created_at`),
  KEY `idx_wallet_ledger_wallet_time` (`wallet_id`, `created_at`),
  KEY `idx_wallet_ledger_txn` (`transaction_record_id`),
  KEY `idx_wallet_ledger_biz_type` (`biz_type`),
  CONSTRAINT `fk_wallet_ledger_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`),
  CONSTRAINT `fk_wallet_ledger_wallet`
    FOREIGN KEY (`wallet_id`) REFERENCES `merchant_wallets` (`id`),
  CONSTRAINT `fk_wallet_ledger_txn`
    FOREIGN KEY (`transaction_record_id`) REFERENCES `transaction_records` (`id`),
  CONSTRAINT `chk_wallet_ledger_amount` CHECK (`amount` > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
