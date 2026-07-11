ALTER TABLE `merchants`
  ADD COLUMN `supported_crypto_chains` TEXT NULL COMMENT 'JSON map currency->chains[]; NULL = system default' AFTER `supported_fiat_currencies`,
  ADD COLUMN `default_crypto_chains` TEXT NULL COMMENT 'JSON map currency->default chain; NULL = system default' AFTER `supported_crypto_chains`;

INSERT INTO `system_configs` (`config_key`, `config_value`, `description`) VALUES
('supported_crypto_chains', '{"USDT":["ETH_ERC20","TRX_TRC20","SOL_Solana","BSC_BEP20"],"USDC":["ETH_ERC20","TRX_TRC20","SOL_Solana","BSC_BEP20"],"BTC":["BTC"]}', 'Default supported chains per crypto currency (JSON)'),
('default_crypto_chains', '{"USDT":"TRX_TRC20","USDC":"ETH_ERC20","BTC":"BTC"}', 'Default selected chain per crypto currency (JSON)')
ON DUPLICATE KEY UPDATE
  `config_value` = VALUES(`config_value`),
  `description` = VALUES(`description`);
