ALTER TABLE `exchange_orders`
  ADD COLUMN `fail_reason` TEXT NULL COMMENT '兑换失败原因' AFTER `completed_at`;
