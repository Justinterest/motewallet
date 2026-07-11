DELETE FROM `system_configs`
WHERE `config_key` IN ('supported_crypto_chains', 'default_crypto_chains');

ALTER TABLE `merchants`
  DROP COLUMN `default_crypto_chains`,
  DROP COLUMN `supported_crypto_chains`;
