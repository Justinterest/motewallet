ALTER TABLE `fee_templates`
  ADD COLUMN `exchange_fee_deduction_method` VARCHAR(20) NOT NULL DEFAULT 'WALLET' COMMENT 'WALLET / RECEIVED_AMOUNT' AFTER `is_default`,
  ADD COLUMN `crypto_withdrawal_fee_deduction_method` VARCHAR(20) NOT NULL DEFAULT 'WALLET' COMMENT 'WALLET / RECEIVED_AMOUNT' AFTER `exchange_fee_deduction_method`,
  ADD COLUMN `fiat_withdrawal_fee_deduction_method` VARCHAR(20) NOT NULL DEFAULT 'WALLET' COMMENT 'WALLET / RECEIVED_AMOUNT' AFTER `crypto_withdrawal_fee_deduction_method`;

ALTER TABLE `transaction_records`
  ADD COLUMN `fee_deduction_method` VARCHAR(20) NOT NULL DEFAULT 'WALLET' COMMENT '订单创建时的手续费扣取方式快照' AFTER `platform_fee_currency`;
