ALTER TABLE `admin_users`
  ADD COLUMN `must_change_password` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '1 = must change password on next login' AFTER `password_hash`;
