ALTER TABLE `admin_users`
  ADD COLUMN `totp_secret` VARCHAR(64) NULL COMMENT 'TOTP base32 secret' AFTER `password_hash`,
  ADD COLUMN `totp_enabled` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '1 = 2FA enabled' AFTER `totp_secret`,
  ADD COLUMN `totp_pending_secret` VARCHAR(64) NULL COMMENT 'Pending TOTP secret during rebind' AFTER `totp_enabled`;
