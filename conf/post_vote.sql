-- 创建帖子投票记录表
CREATE TABLE IF NOT EXISTS `post_vote` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `post_id` bigint(20) NOT NULL COMMENT '帖子ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户ID',
    `vote_type` tinyint(4) NOT NULL COMMENT '投票类型：1赞成 0取消 -1反对',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_user` (`post_id`, `user_id`),
    KEY `idx_post_id` (`post_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子投票记录表';

-- 创建帖子投票统计表
CREATE TABLE IF NOT EXISTS `post_vote_count` (
    `post_id` bigint(20) NOT NULL COMMENT '帖子ID',
    `up_vote_count` int(11) NOT NULL DEFAULT '0' COMMENT '赞成票数',
    `down_vote_count` int(11) NOT NULL DEFAULT '0' COMMENT '反对票数',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子投票统计表'; 