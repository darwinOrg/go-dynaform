DROP TABLE IF EXISTS `test_entity`;
CREATE TABLE `test_entity`
(
    `id`          BIGINT      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `foo`         VARCHAR(64) NOT NULL,
    `bar`         VARCHAR(64) NOT NULL,
    `data`        TEXT COMMENT '表单',
    `created_at`  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modified_at` DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`)

) ENGINE = InnoDB
  CHARSET = utf8mb4;