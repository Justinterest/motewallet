ALTER TABLE `transaction_records`
  DROP COLUMN `fee_deduction_method`;

ALTER TABLE `fee_templates`
  DROP COLUMN `fiat_withdrawal_fee_deduction_method`,
  DROP COLUMN `crypto_withdrawal_fee_deduction_method`,
  DROP COLUMN `exchange_fee_deduction_method`;
