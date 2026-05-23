-- ============================================================
-- Motewallet 初始化 Schema
-- 共 18 张表，按 FK 依赖顺序创建
-- ============================================================

-- ------------------------------------------------------------
-- 1. fee_templates（手续费模板）
-- ------------------------------------------------------------
CREATE TABLE `fee_templates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,
  `description` VARCHAR(255) NULL,
  `is_default` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fee_templates_is_default` (`is_default`),
  KEY `idx_fee_templates_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 2. merchants（商户）
-- ------------------------------------------------------------
CREATE TABLE `merchants` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(128) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `kun_sub_customer_no` VARCHAR(64) NULL,
  `fee_template_id` BIGINT UNSIGNED NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'PENDING_AGREEMENT',
  `kyc_auth_id` VARCHAR(128) NULL,
  `kyc_status` VARCHAR(20) NOT NULL DEFAULT 'NONE',
  `kyc_fail_reason` TEXT NULL,
  `kyc_submitted_at` DATETIME(3) NULL,
  `kyc_completed_at` DATETIME(3) NULL,
  `agreement_signed_at` DATETIME(3) NULL,
  `frozen_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_merchants_email` (`email`),
  UNIQUE KEY `uk_merchants_kun_sub_customer_no` (`kun_sub_customer_no`),
  KEY `idx_merchants_status` (`status`),
  KEY `idx_merchants_kyc_status` (`kyc_status`),
  KEY `idx_merchants_fee_template_id` (`fee_template_id`),
  KEY `idx_merchants_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_merchants_fee_template`
    FOREIGN KEY (`fee_template_id`) REFERENCES `fee_templates` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 3. fee_template_exchange_items（兑换手续费配置）
-- ------------------------------------------------------------
CREATE TABLE `fee_template_exchange_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `fee_template_id` BIGINT UNSIGNED NOT NULL,
  `from_currency` VARCHAR(10) NOT NULL,
  `to_currency` VARCHAR(10) NOT NULL,
  `fee_rate` DECIMAL(10,6) NOT NULL DEFAULT 0,
  `min_fee` DECIMAL(28,8) NOT NULL DEFAULT 0,
  `min_fee_currency` VARCHAR(10) NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_exchange_items_template_pair` (`fee_template_id`, `from_currency`, `to_currency`),
  CONSTRAINT `fk_exchange_items_template`
    FOREIGN KEY (`fee_template_id`) REFERENCES `fee_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 4. fee_template_crypto_withdrawal_items（数币提现手续费配置）
-- ------------------------------------------------------------
CREATE TABLE `fee_template_crypto_withdrawal_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `fee_template_id` BIGINT UNSIGNED NOT NULL,
  `currency` VARCHAR(10) NOT NULL,
  `chain` VARCHAR(20) NOT NULL,
  `fee_rate` DECIMAL(10,6) NOT NULL DEFAULT 0,
  `fixed_fee` DECIMAL(28,8) NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_crypto_withdrawal_items_template_chain` (`fee_template_id`, `currency`, `chain`),
  CONSTRAINT `fk_crypto_withdrawal_items_template`
    FOREIGN KEY (`fee_template_id`) REFERENCES `fee_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 5. fee_template_fiat_withdrawal_items（法币提现手续费配置）
-- ------------------------------------------------------------
CREATE TABLE `fee_template_fiat_withdrawal_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `fee_template_id` BIGINT UNSIGNED NOT NULL,
  `currency` VARCHAR(10) NOT NULL,
  `transfer_type` VARCHAR(10) NOT NULL,
  `fee_rate` DECIMAL(10,6) NOT NULL DEFAULT 0,
  `fixed_fee` DECIMAL(28,8) NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fiat_withdrawal_items_template_type` (`fee_template_id`, `currency`, `transfer_type`),
  CONSTRAINT `fk_fiat_withdrawal_items_template`
    FOREIGN KEY (`fee_template_id`) REFERENCES `fee_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 6. merchant_wallets（商户钱包余额）
-- ------------------------------------------------------------
CREATE TABLE `merchant_wallets` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `account_type` VARCHAR(10) NOT NULL COMMENT 'FUNDING / TRADING',
  `currency` VARCHAR(10) NOT NULL,
  `balance` DECIMAL(28,8) NOT NULL DEFAULT 0,
  `frozen_balance` DECIMAL(28,8) NOT NULL DEFAULT 0 COMMENT '冻结金额（提现/兑换处理中）',
  `version` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_merchant_wallets_account` (`merchant_id`, `account_type`, `currency`),
  KEY `idx_merchant_wallets_merchant_id` (`merchant_id`),
  CONSTRAINT `fk_merchant_wallets_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`),
  CONSTRAINT `chk_merchant_wallets_balance` CHECK (`balance` >= 0),
  CONSTRAINT `chk_merchant_wallets_frozen` CHECK (`frozen_balance` >= 0),
  CONSTRAINT `chk_merchant_wallets_available` CHECK (`balance` >= `frozen_balance`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 7. crypto_addresses（数币白名单地址）
-- ------------------------------------------------------------
CREATE TABLE `crypto_addresses` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `currency` VARCHAR(10) NOT NULL,
  `chain` VARCHAR(20) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `alias` VARCHAR(64) NOT NULL,
  `kun_account_id` VARCHAR(128) NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'INIT',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_crypto_addresses_merchant_id` (`merchant_id`),
  KEY `idx_crypto_addresses_status` (`status`),
  KEY `idx_crypto_addresses_deleted_at` (`deleted_at`),
  KEY `idx_crypto_addresses_merchant_currency` (`merchant_id`, `currency`, `chain`),
  CONSTRAINT `fk_crypto_addresses_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 8. bank_accounts（法币提现银行账户）
-- ------------------------------------------------------------
CREATE TABLE `bank_accounts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `kun_account_id` VARCHAR(128) NULL,
  `currency_list` VARCHAR(50) NOT NULL,
  `transfer_type` VARCHAR(10) NOT NULL,
  `account_no` VARCHAR(128) NOT NULL COMMENT 'AES 加密存储',
  `account_name` VARCHAR(128) NOT NULL,
  `bank_name` VARCHAR(128) NULL,
  `bank_code` VARCHAR(20) NULL,
  `swift_code` VARCHAR(20) NULL,
  `payee_country_code` VARCHAR(10) NULL,
  `payee_address` VARCHAR(255) NULL,
  `middle_swift_code` VARCHAR(20) NULL,
  `area` VARCHAR(10) NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_bank_accounts_merchant_id` (`merchant_id`),
  KEY `idx_bank_accounts_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_bank_accounts_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 9. transaction_records（平台资金流水账本 — 纯账本层）
-- ------------------------------------------------------------
CREATE TABLE `transaction_records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `platform_order_id` VARCHAR(64) NOT NULL,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `type` VARCHAR(20) NOT NULL COMMENT 'DEPOSIT/WITHDRAWAL/EXCHANGE/TRANSFER',
  `sub_type` VARCHAR(30) NULL,
  `amount` DECIMAL(28,8) NOT NULL,
  `currency` VARCHAR(10) NOT NULL,
  `platform_fee` DECIMAL(28,8) NOT NULL DEFAULT 0,
  `platform_fee_currency` VARCHAR(10) NULL,
  `actual_amount` DECIMAL(28,8) NULL,
  `remark` VARCHAR(500) NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  `completed_at` DATETIME(3) NULL,
  `failed_reason` VARCHAR(500) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_records_platform_order_id` (`platform_order_id`),
  KEY `idx_transaction_records_merchant_id` (`merchant_id`),
  KEY `idx_transaction_records_type` (`type`),
  KEY `idx_transaction_records_status` (`status`),
  KEY `idx_transaction_records_created_at` (`created_at`),
  KEY `idx_transaction_records_merchant_type` (`merchant_id`, `type`),
  KEY `idx_transaction_records_merchant_status` (`merchant_id`, `status`),
  CONSTRAINT `fk_transaction_records_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 10. admin_users（管理员账号）
-- ------------------------------------------------------------
CREATE TABLE `admin_users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(64) NOT NULL,
  `email` VARCHAR(128) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `role` VARCHAR(20) NOT NULL DEFAULT 'OPERATOR' COMMENT 'SUPER_ADMIN/OPERATOR/FINANCE',
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `last_login_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_users_username` (`username`),
  UNIQUE KEY `uk_admin_users_email` (`email`),
  KEY `idx_admin_users_role` (`role`),
  KEY `idx_admin_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 11. audit_logs（审计日志）
-- ------------------------------------------------------------
CREATE TABLE `audit_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `operator_id` BIGINT UNSIGNED NOT NULL,
  `operator_type` VARCHAR(10) NOT NULL COMMENT 'ADMIN/MERCHANT',
  `action` VARCHAR(64) NOT NULL,
  `target_type` VARCHAR(32) NULL,
  `target_id` VARCHAR(64) NULL,
  `detail` JSON NULL,
  `ip_address` VARCHAR(45) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_audit_logs_operator` (`operator_id`, `operator_type`),
  KEY `idx_audit_logs_action` (`action`),
  KEY `idx_audit_logs_target` (`target_type`, `target_id`),
  KEY `idx_audit_logs_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 12. system_configs（系统配置 KV）
-- ------------------------------------------------------------
CREATE TABLE `system_configs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `config_key` VARCHAR(64) NOT NULL,
  `config_value` TEXT NOT NULL,
  `description` VARCHAR(255) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_system_configs_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 13. webhook_logs（KUN Webhook 回调日志）
-- ------------------------------------------------------------
CREATE TABLE `webhook_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id` VARCHAR(128) NOT NULL,
  `event_topic` VARCHAR(64) NOT NULL,
  `event_type` VARCHAR(20) NOT NULL,
  `customer_no` VARCHAR(64) NULL,
  `raw_data` JSON NOT NULL,
  `process_status` VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  `process_error` TEXT NULL,
  `processed_at` DATETIME(3) NULL,
  `received_at` DATETIME(3) NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_webhook_logs_event` (`event_id`, `event_topic`),
  KEY `idx_webhook_logs_process_status` (`process_status`),
  KEY `idx_webhook_logs_customer_no` (`customer_no`),
  KEY `idx_webhook_logs_received_at` (`received_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 14. system_announcements（系统公告）
-- ------------------------------------------------------------
CREATE TABLE `system_announcements` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(128) NOT NULL,
  `content` TEXT NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'DRAFT' COMMENT 'DRAFT/PUBLISHED',
  `published_at` DATETIME(3) NULL,
  `created_by` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_system_announcements_status` (`status`),
  KEY `idx_system_announcements_published_at` (`published_at`),
  KEY `idx_system_announcements_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_system_announcements_created_by`
    FOREIGN KEY (`created_by`) REFERENCES `admin_users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 15. deposit_orders（充值订单）
-- ------------------------------------------------------------
CREATE TABLE `deposit_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `transaction_record_id` BIGINT UNSIGNED NOT NULL,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `currency` VARCHAR(10) NOT NULL,
  `chain` VARCHAR(20) NOT NULL,
  `to_address` VARCHAR(255) NOT NULL COMMENT '充值目标地址',
  `from_address` VARCHAR(255) NULL COMMENT '发送方地址（Webhook 回调填充）',
  `tx_id` VARCHAR(128) NULL COMMENT '链上交易哈希',
  `kun_order_id` VARCHAR(128) NULL,
  `confirmations` INT UNSIGNED NULL COMMENT '链上确认数',
  `completed_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_deposit_orders_txn_record` (`transaction_record_id`),
  KEY `idx_deposit_orders_merchant_id` (`merchant_id`),
  KEY `idx_deposit_orders_tx_id` (`tx_id`),
  KEY `idx_deposit_orders_kun_order_id` (`kun_order_id`),
  CONSTRAINT `fk_deposit_orders_txn_record`
    FOREIGN KEY (`transaction_record_id`) REFERENCES `transaction_records` (`id`),
  CONSTRAINT `fk_deposit_orders_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 16. withdrawal_orders（提现订单 — 含审核流程）
-- ------------------------------------------------------------
CREATE TABLE `withdrawal_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `transaction_record_id` BIGINT UNSIGNED NOT NULL COMMENT '关联流水',
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `withdrawal_type` VARCHAR(20) NOT NULL COMMENT 'CRYPTO / FIAT',

  -- 目标信息
  `crypto_address_id` BIGINT UNSIGNED NULL COMMENT '关联白名单地址',
  `bank_account_id` BIGINT UNSIGNED NULL COMMENT '关联银行账户',
  `to_address` VARCHAR(255) NULL COMMENT '提现目标地址/银行账号（冗余快照）',

  -- 数币提现字段
  `chain` VARCHAR(20) NULL,
  `tx_id` VARCHAR(128) NULL COMMENT '链上交易哈希',

  -- 法币提现字段
  `transfer_type` VARCHAR(10) NULL COMMENT 'LOCAL/CHATS/TT',
  `purpose` VARCHAR(10) NULL COMMENT '法币提现用途',

  -- 审核流程
  `review_status` VARCHAR(20) NOT NULL DEFAULT 'PENDING_REVIEW',
  `reviewer_id` BIGINT UNSIGNED NULL,
  `reviewer_type` VARCHAR(10) NULL COMMENT 'ADMIN / SYSTEM',
  `reviewed_at` DATETIME(3) NULL,
  `review_remark` VARCHAR(500) NULL,

  -- KUN 对接
  `kun_order_id` VARCHAR(128) NULL,
  `kun_request_no` VARCHAR(64) NULL,
  `kun_fee` DECIMAL(28,8) NULL DEFAULT 0,
  `kun_fee_currency` VARCHAR(10) NULL,
  `kun_submitted_at` DATETIME(3) NULL COMMENT 'KUN 提交时间（审核通过后）',

  `failed_reason` VARCHAR(500) NULL,
  `completed_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_withdrawal_orders_txn_record` (`transaction_record_id`),
  UNIQUE KEY `uk_withdrawal_orders_kun_request_no` (`kun_request_no`),
  KEY `idx_withdrawal_orders_merchant_id` (`merchant_id`),
  KEY `idx_withdrawal_orders_review_status` (`review_status`),
  KEY `idx_withdrawal_orders_reviewer_id` (`reviewer_id`),
  KEY `idx_withdrawal_orders_withdrawal_type` (`withdrawal_type`),
  KEY `idx_withdrawal_orders_created_at` (`created_at`),
  CONSTRAINT `fk_withdrawal_orders_txn_record`
    FOREIGN KEY (`transaction_record_id`) REFERENCES `transaction_records` (`id`),
  CONSTRAINT `fk_withdrawal_orders_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`),
  CONSTRAINT `fk_withdrawal_orders_crypto_address`
    FOREIGN KEY (`crypto_address_id`) REFERENCES `crypto_addresses` (`id`),
  CONSTRAINT `fk_withdrawal_orders_bank_account`
    FOREIGN KEY (`bank_account_id`) REFERENCES `bank_accounts` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 17. exchange_orders（兑换订单）
-- ------------------------------------------------------------
CREATE TABLE `exchange_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `transaction_record_id` BIGINT UNSIGNED NOT NULL,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `exchange_type` VARCHAR(20) NOT NULL COMMENT 'SPOT_EXCHANGE / ONE_TO_ONE',
  `from_currency` VARCHAR(10) NOT NULL,
  `from_amount` DECIMAL(28,8) NOT NULL,
  `to_currency` VARCHAR(10) NOT NULL,
  `to_amount` DECIMAL(28,8) NULL COMMENT '成交到账金额（回调后填充）',
  `exchange_rate` DECIMAL(20,10) NULL,
  `quote_id` VARCHAR(128) NULL COMMENT '询价 ID（报价锁定用）',
  `auto_transfer` VARCHAR(3) NULL DEFAULT 'NO' COMMENT '1:1 交易是否自动划转（YES/NO）',
  `kun_order_id` VARCHAR(128) NULL,
  `kun_request_no` VARCHAR(64) NULL,
  `kun_fee` DECIMAL(28,8) NULL DEFAULT 0,
  `kun_fee_currency` VARCHAR(10) NULL,
  `completed_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_exchange_orders_txn_record` (`transaction_record_id`),
  UNIQUE KEY `uk_exchange_orders_kun_request_no` (`kun_request_no`),
  KEY `idx_exchange_orders_merchant_id` (`merchant_id`),
  KEY `idx_exchange_orders_kun_order_id` (`kun_order_id`),
  KEY `idx_exchange_orders_exchange_type` (`exchange_type`),
  CONSTRAINT `fk_exchange_orders_txn_record`
    FOREIGN KEY (`transaction_record_id`) REFERENCES `transaction_records` (`id`),
  CONSTRAINT `fk_exchange_orders_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- 18. transfer_orders（划转订单）
-- ------------------------------------------------------------
CREATE TABLE `transfer_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `transaction_record_id` BIGINT UNSIGNED NOT NULL,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `from_account_type` VARCHAR(10) NOT NULL COMMENT 'FUNDING / TRADING',
  `to_account_type` VARCHAR(10) NOT NULL COMMENT 'FUNDING / TRADING',
  `kun_order_id` VARCHAR(128) NULL,
  `kun_request_no` VARCHAR(64) NULL,
  `completed_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transfer_orders_txn_record` (`transaction_record_id`),
  UNIQUE KEY `uk_transfer_orders_kun_request_no` (`kun_request_no`),
  KEY `idx_transfer_orders_merchant_id` (`merchant_id`),
  CONSTRAINT `fk_transfer_orders_txn_record`
    FOREIGN KEY (`transaction_record_id`) REFERENCES `transaction_records` (`id`),
  CONSTRAINT `fk_transfer_orders_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
