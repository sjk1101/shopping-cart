DROP TABLE IF EXISTS `product`;
CREATE TABLE IF NOT EXISTS `product` (
    `id` INT NOT NULL COMMENT 'PK',
    `name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '產品名稱',
    `amount` INT NOT NULL  DEFAULT 0 COMMENT '產品價格',
    `inventory` INT NOT NULL COMMENT '庫存',
    `image` VARCHAR(255) NOT NULL COMMENT '圖片路徑',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '建立時間',
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新時間',
    PRIMARY KEY (`id`)
    ) COMMENT = '產品';