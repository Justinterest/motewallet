ALTER TABLE `merchants`
  ADD COLUMN `supported_crypto_currencies` VARCHAR(128) NULL COMMENT 'Comma-separated crypto codes; NULL = system default' AFTER `fee_template_id`,
  ADD COLUMN `supported_fiat_currencies` VARCHAR(128) NULL COMMENT 'Comma-separated fiat codes; NULL = system default' AFTER `supported_crypto_currencies`;
