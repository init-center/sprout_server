# 用户表
CREATE TABLE `t_user` (
                          `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
                          `uid` varchar(20) NOT NULL COMMENT '用户id',
                          `gender` tinyint(1) UNSIGNED NOT NULL DEFAULT "0" COMMENT '性别',
                          `name` varchar(12) NOT NULL COMMENT '昵称',
                          `email` varchar(64) COMMENT '邮箱',
                          `tel` int(11)  COMMENT '电话号码',
                          `password` varchar(64) NOT NULL COMMENT '密码',
                          `birthday` date DEFAULT NULL COMMENT '生日',
                          `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
                          `group` tinyint(3) UNSIGNED NOT NULL DEFAULT '2' COMMENT '用户组',
                          `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
                          `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          `delete_time` timestamp NULL DEFAULT NULL COMMENT '注销时间',
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `idx_name` (`name`),
                          UNIQUE KEY `idx_uid` (`uid`),
                          UNIQUE KEY `idx_tel` (`tel`),
                          UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';