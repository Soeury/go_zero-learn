CREATE TABLE `order` (
    `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'primary_key',
    `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create_time',
    `create_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'create_by',
    `update_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update_time',
    `update_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'update_user',
    `version` SMALLINT(5) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'opt_version',
    `is_del` TINYINT(4) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'delete or not: 0-no , 1-yes',

    `user_id` BIGINT(20) UNSIGNED NOT NULL COMMENT 'user_id',
    `order_id` BIGINT(20) UNSIGNED NOT NULL COMMENT 'order_id',
    `trade_id` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'trade_id',
    `pay_channel` TINYINT(4) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'pay_way',
    `status` INT UNSIGNED NOT NULL DEFAULT '0' COMMENT 'order_status: 100-create_order , 200-already_pay , 300-trade_close',
    `pay_amount` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'pay_money(yuan)',

    UNIQUE KEY `idx_order_id` (`order_id`) USING BTREE,
    INDEX (trade_id),
    INDEX (user_id),
    INDEX (is_del)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT = 'order table';
