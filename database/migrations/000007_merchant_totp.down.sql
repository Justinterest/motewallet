ALTER TABLE `merchants`
  DROP COLUMN `totp_pending_secret`,
  DROP COLUMN `totp_enabled`,
  DROP COLUMN `totp_secret`;
