# 用户表
CREATE TABLE `t_user` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `uid` varchar(20) NOT NULL COMMENT '用户id',
    `gender` tinyint(1) UNSIGNED NOT NULL DEFAULT '0' COMMENT '性别',
    `name` varchar(12) NOT NULL COMMENT '昵称',
    `email` varchar(64) COMMENT '邮箱',
    `tel` int(11)  COMMENT '电话号码',
    `password` varchar(64) NOT NULL COMMENT '密码',
    `birthday` date DEFAULT NULL COMMENT '生日',
    `avatar` varchar(2083) DEFAULT NULL COMMENT '头像地址',
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

# 分类表
CREATE TABLE `t_post_category` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `name` varchar(12) NOT NULL COMMENT '名称',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT  CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='分类表';

# 标签表
CREATE TABLE `t_post_tag` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `name` varchar(12) NOT NULL COMMENT '名称',
     PRIMARY KEY (`id`),
     UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT  CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='标签表';

#文章和标签关联表
CREATE TABLE `t_post_tag_relation` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `pid` bigint(20) UNSIGNED NOT NULL COMMENT '文章id',
    `tid` bigint(20) UNSIGNED NOT NULL COMMENT '标签id',
    PRIMARY KEY (`id`),
    KEY `idx_pid` (`pid`),
    KEY `idx_tid` (`tid`)
) ENGINE=InnoDB DEFAULT  CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='文章标签关系表';

#文章表
CREATE TABLE `t_post` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `uid` varchar(20) NOT NULL COMMENT '用户id',
    `pid` bigint(20) UNSIGNED NOT NULL COMMENT '文章id',
    `category` bigint(20) UNSIGNED NOT NULL COMMENT '分类id',
    `title` varchar(128) NOT NULL COMMENT '标题',
    `cover` varchar(2083) NOT NULL COMMENT '封面地址',
    `bgm` varchar(2083) NOT NULL COMMENT '背景音乐地址',
    `summary` varchar(128) NOT NULL COMMENT '摘要',
    `content` TEXT NOT NULL COMMENT '文章内容',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_pid` (`pid`),
    KEY `idx_uid` (`uid`),
    KEY `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='文章表';

#文章配置项表
CREATE TABLE `t_post_config` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `pid` bigint(20) UNSIGNED NOT NULL COMMENT '文章id',
    `comment_open` tinyint(1) UNSIGNED NOT NULL DEFAULT '1' COMMENT '是否开启评论',
    `display` tinyint(1) UNSIGNED NOT NULL DEFAULT '1' COMMENT '显示隐藏',
    `top_time` timestamp NULL DEFAULT NULL COMMENT '置顶时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY  `idx_pid` (`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='文章配置项表';

#文章阅读量表
CREATE TABLE `t_post_views` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `pid` bigint(20) UNSIGNED NOT NULL COMMENT '文章id',
    `views` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '阅读量',
    PRIMARY KEY (`id`),
    UNIQUE KEY  `idx_pid` (`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='文章阅读量表';


#评论表
CREATE TABLE `t_post_comment` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `cid` bigint(20) UNSIGNED NOT NULL COMMENT '评论id',
    `pid` bigint(20) UNSIGNED NOT NULL COMMENT '文章id',
    `uid` varchar(20) NOT NULL COMMENT '评论用户id',
    `target_cid` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '评论目标评论id',
    `target_uid` varchar(20) NULL DEFAULT NULL COMMENT '评论目标用户id',
    `parent_cid` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT  '父评论id',
    `parent_uid` varchar(20) NULL DEFAULT NULL COMMENT '父评论用户id',
    `content` varchar(1024) NOT NULl COMMENT '评论内容',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `review_status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '审核状态',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_cid` (`cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='文章评论表';

#文章喜欢表
CREATE TABLE `t_post_favorite` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `pid` bigint(20) UNSIGNED NOT NULL COMMENT '文章id',
    `uid` varchar(20) NOT NULL COMMENT '用户id',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '喜欢时间',
    PRIMARY KEY (`id`),
    KEY `idx_pid` (`pid`),
    KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='文章喜欢表';