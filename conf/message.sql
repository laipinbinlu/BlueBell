-- Active: 1744179549288@@127.0.0.1@3306@bluebell
-- 创建消息表
CREATE TABLE IF NOT EXISTS `message` (
    `msg_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    `sender_uid` bigint(20) NOT NULL COMMENT '发送者用户ID',
    `receiver_uid` bigint(20) NOT NULL COMMENT '接收者用户ID',
    `content` text NOT NULL COMMENT '消息内容',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读 0-未读 1-已读',
    PRIMARY KEY (`msg_id`),
    KEY `idx_sender_uid` (`sender_uid`),
    KEY `idx_receiver_uid` (`receiver_uid`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='私信消息表'; 