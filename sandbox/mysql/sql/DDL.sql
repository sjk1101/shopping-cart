DROP TABLE IF EXISTS `product`;
CREATE TABLE `product`
(
    `name`       varchar(512)   NOT NULL COMMENT '產品名稱',
    `number`     varchar(16) DEFAULT '' COMMENT '產品貨號',
    `cost`       decimal(10, 4) NOT NULL COMMENT '商品成本',
    `quantity`   int(11) COMMENT '庫存',
    `image`      varchar(255)   NOT NULL COMMENT '圖片路徑',
    `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP (6) COMMENT '建立時間',
    `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP (6) COMMENT '更新時間',
    KEY          `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='產品';

DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user`
(
    `id`         bigint(20) NOT NULL COMMENT 'ID',
    `account`    varchar(20)  NOT NULL DEFAULT '' COMMENT '管理者account',
    `name`       varchar(100) NOT NULL DEFAULT '' COMMENT '管理者名稱',
    `password`   varchar(128) NOT NULL COMMENT '管理者密碼',
    `salt`       varchar(36)  NOT NULL COMMENT '加密金鑰',
    `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP (6) COMMENT '建立時間',
    `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP (6) COMMENT '更新時間',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理者';

DROP TABLE IF EXISTS `shopee_order_detail`;
CREATE TABLE `shopee_order_detail`
(
    `order_id`           varchar(255)   NOT NULL COMMENT '訂單編號',
    `order_created_at`   datetime(6) NOT NULL COMMENT '訂單建立時間',
    `is_established`     tinyint(1) NOT NULL COMMENT '訂單是否成立',
    `order_completed_at` datetime(6) DEFAULT NULL COMMENT '訂單完成時間',
    `product`            varchar(512)   NOT NULL COMMENT '商品名稱',
    `product_price`      decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '商品金額',
    `product_cost`       decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '商品成本',
    `quantity`           int(11) NOT NULL DEFAULT '0.00' COMMENT '商品數量',
    `coupon_discount`    decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '賣場優惠券',
    `deal_fee`           decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '成交手續費',
    `activity_fee`       decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '活動服務費',
    `cash_flow_cost`     decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '金流服務費',
    `created_at`         datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP (6) COMMENT '建立時間',
    `updated_at`         datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP (6) COMMENT '更新時間',
    KEY                  `idx_order_completed_at` (`order_completed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='蝦皮訂單明細';

DROP TABLE IF EXISTS `shopee_order`;
CREATE TABLE `shopee_order`
(
    `order_id`            varchar(255)   NOT NULL COMMENT '訂單編號',
    `order_created_at`    datetime(6) NOT NULL COMMENT '訂單建立時間',
    `is_established`      tinyint(1) NOT NULL COMMENT '訂單是否成立',
    `order_completed_at`  datetime(6) DEFAULT NULL COMMENT '訂單完成時間',
    `allocate_at`         datetime(6) DEFAULT NULL COMMENT '撥款日',
    `coupon_discount`     decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '賣場優惠券',
    `deal_fee`            decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '成交手續費',
    `activity_fee`        decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '活動服務費',
    `cash_flow_cost`      decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '金流服務費',
    `total_product_price` decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '總商品金額',
    `total_product_cost`  decimal(10, 4) NOT NULL DEFAULT '0.0000' COMMENT '總商品成本',
    `created_at`          datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP (6) COMMENT '建立時間',
    `updated_at`          datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP (6) COMMENT '更新時間',
    PRIMARY KEY (`order_id`),
    KEY                   `idx_order_completed_at` (`order_completed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='蝦皮訂單';
