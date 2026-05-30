-- merchant_kyc_submissions: local copy of KUN sub-merchant onboarding authentication payloads
CREATE TABLE `merchant_kyc_submissions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `merchant_id` BIGINT UNSIGNED NOT NULL,
  `kun_request_no` VARCHAR(64) NOT NULL,
  `kun_auth_id` VARCHAR(128) NULL,
  `payload` JSON NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'SUBMITTED',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_merchant_kyc_submissions_request_no` (`kun_request_no`),
  KEY `idx_merchant_kyc_submissions_merchant_id` (`merchant_id`),
  CONSTRAINT `fk_merchant_kyc_submissions_merchant`
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
