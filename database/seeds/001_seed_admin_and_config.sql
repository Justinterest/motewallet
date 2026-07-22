-- ============================================================
-- Motewallet 种子数据（开发/测试环境）
-- 幂等：可重复执行（INSERT ... ON DUPLICATE KEY UPDATE）
-- ============================================================

-- Super Admin account (Password: Admin@2024, bcrypt hash)
INSERT INTO `admin_users` (`username`, `email`, `password_hash`, `role`, `status`) VALUES
('superadmin', 'admin@motewallet.com', '$2a$10$3eKTYEkiHEf87WAg1bta9uj0TRqnFEkGPqzyaN3HZen5F8gMA0iz.', 'SUPER_ADMIN', 'ACTIVE')
ON DUPLICATE KEY UPDATE
  `password_hash` = VALUES(`password_hash`),
  `role` = VALUES(`role`),
  `status` = VALUES(`status`);

-- Default fee template
INSERT INTO `fee_templates` (`id`, `name`, `description`, `is_default`) VALUES
(1, 'Default Template', 'System default fee template for new merchants', 1)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `description` = VALUES(`description`),
  `is_default` = VALUES(`is_default`);

-- Exchange fee items (default 0.3% rate)
INSERT INTO `fee_template_exchange_items` (`fee_template_id`, `from_currency`, `to_currency`, `fee_rate`, `min_fee`, `min_fee_currency`) VALUES
(1, 'USDT', 'USD',  0.003000, 1.00000000, 'USD'),
(1, 'USD',  'USDT', 0.003000, 1.00000000, 'USDT'),
(1, 'USDT', 'HKD',  0.003000, 5.00000000, 'HKD'),
(1, 'HKD',  'USDT', 0.003000, 5.00000000, 'USDT'),
(1, 'USDT', 'EUR',  0.003000, 1.00000000, 'EUR'),
(1, 'EUR',  'USDT', 0.003000, 1.00000000, 'USDT'),
(1, 'USDT', 'USDC', 0.001000, 1.00000000, 'USDC'),
(1, 'USDC', 'USDT', 0.001000, 1.00000000, 'USDT'),
(1, 'USDT', 'BTC',  0.003000, 0.00010000, 'BTC'),
(1, 'BTC',  'USDT', 0.003000, 1.00000000, 'USDT')
ON DUPLICATE KEY UPDATE
  `fee_rate` = VALUES(`fee_rate`),
  `min_fee` = VALUES(`min_fee`),
  `min_fee_currency` = VALUES(`min_fee_currency`);

-- Crypto withdrawal fee items
INSERT INTO `fee_template_crypto_withdrawal_items` (`fee_template_id`, `currency`, `chain`, `fee_rate`, `fixed_fee`) VALUES
(1, 'USDT', 'ETH_ERC20',   0.000000, 5.00000000),
(1, 'USDT', 'TRX_TRC20',   0.000000, 1.00000000),
(1, 'USDT', 'SOL_Solana',   0.000000, 1.00000000),
(1, 'USDT', 'BSC_BEP20',   0.000000, 1.00000000),
(1, 'USDC', 'ETH_ERC20',   0.000000, 5.00000000),
(1, 'USDC', 'TRX_TRC20',   0.000000, 1.00000000),
(1, 'BTC',  'BTC',          0.000000, 0.00050000)
ON DUPLICATE KEY UPDATE
  `fee_rate` = VALUES(`fee_rate`),
  `fixed_fee` = VALUES(`fixed_fee`);

-- Fiat withdrawal fee items
INSERT INTO `fee_template_fiat_withdrawal_items` (`fee_template_id`, `currency`, `transfer_type`, `fee_rate`, `fixed_fee`) VALUES
(1, 'USD', 'LOCAL',  0.000000, 25.00000000),
(1, 'USD', 'TT',     0.000000, 40.00000000),
(1, 'HKD', 'LOCAL',  0.000000, 50.00000000),
(1, 'HKD', 'CHATS',  0.000000, 100.00000000),
(1, 'HKD', 'TT',     0.000000, 200.00000000),
(1, 'EUR', 'LOCAL',  0.000000, 20.00000000),
(1, 'EUR', 'TT',     0.000000, 35.00000000)
ON DUPLICATE KEY UPDATE
  `fee_rate` = VALUES(`fee_rate`),
  `fixed_fee` = VALUES(`fixed_fee`);

-- System configs（supported_crypto_chains / default_crypto_chains 也可能由 migration 000005 写入）
INSERT INTO `system_configs` (`config_key`, `config_value`, `description`) VALUES
('kun_api_base_url', 'https://api.kun.global', 'KUN API base URL'),
('kun_region_code', 'KUN_PL', 'KUN region code (Poland)'),
('supported_fiat_currencies', 'USD,HKD,EUR', 'Supported fiat currencies'),
('supported_crypto_currencies', 'USDT,USDC,BTC', 'Supported crypto currencies'),
('supported_crypto_chains', '{"USDT":["ETH_ERC20","TRX_TRC20","SOL_Solana","BSC_BEP20"],"USDC":["ETH_ERC20","TRX_TRC20","SOL_Solana","BSC_BEP20"],"BTC":["BTC"]}', 'Default supported chains per crypto currency (JSON)'),
('default_crypto_chains', '{"USDT":"TRX_TRC20","USDC":"ETH_ERC20","BTC":"BTC"}', 'Default selected chain per crypto currency (JSON)'),
('platform_name', 'Motewallet', 'Platform display name'),
('withdrawal_auto_approve_crypto_usd', '0', 'Auto-approve threshold for crypto withdrawal (USD equivalent, 0 = all manual)'),
('withdrawal_auto_approve_fiat_usd', '0', 'Auto-approve threshold for fiat withdrawal (USD equivalent, 0 = all manual)')
ON DUPLICATE KEY UPDATE
  `config_value` = VALUES(`config_value`),
  `description` = VALUES(`description`);
