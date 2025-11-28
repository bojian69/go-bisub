-- 创建数据库
CREATE DATABASE IF NOT EXISTS go_sub DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE go_sub;

-- 创建订阅表
CREATE TABLE IF NOT EXISTS `sub_subscription_theme` (
	`id` bigint unsigned NOT NULL COMMENT '主键Id',
	`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
	`type` char(1) NOT NULL DEFAULT '' COMMENT '订阅类型 [ref:sub_refs] ref_field:SUBSCRIPTION_TYPE',
	`sub_key` varchar(120) NOT NULL DEFAULT '' COMMENT '订阅key',
	`version` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '版本号',
	`title` varchar(240) NOT NULL COMMENT '订阅标题',
	`abstract` tinytext NOT NULL COMMENT '订阅简介',
	`status` char(1) NOT NULL DEFAULT '' COMMENT '状态[ref:sub_refs] ref_field:SUBSCRIPTION_STATUS',
	`created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人ID uhomse_sso[ref:sso_user]',
	`extra_config` json NOT NULL COMMENT '订阅扩展配置{"sql_content":"订阅数据SQL","sql_replace":"SQL替换变量说明","example":"示例说明"}',
	PRIMARY KEY (`id`),
	UNIQUE KEY `uk_type_subkey_version` (`type`,`sub_key`,`version`),
	KEY `idx_title` (`title`)
) DEFAULT CHARACTER SET=utf8 COMMENT='订阅服务主题表';

-- 创建统计表
CREATE TABLE IF NOT EXISTS `sub_logs_bidata_response` (
	`id` bigint unsigned NOT NULL COMMENT '主键ID',
	`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
	`sub_key` varchar(120) NOT NULL DEFAULT '' COMMENT '订阅key',
	`version` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '订阅版本号',
	`execution_duration` mediumint unsigned NOT NULL DEFAULT 0 COMMENT '执行耗时 单位：毫秒',
	`request_url` varchar(1000) NOT NULL DEFAULT '' COMMENT '请求链接',
	`request_response` json NOT NULL COMMENT '请求详情json {"params":"请求参数","instance_sql":"执行实例SQL","instance_source":"实例来源","request_ip":"请求来源IP","version":"版本号"}',
	`instance_source` varchar(120) NOT NULL DEFAULT '' COMMENT '数据实例标识',
	PRIMARY KEY (`id`),
	KEY `idx_subkey_version_instancesource` (`sub_key`,`version`,`instance_source`),
	KEY `idx_createdat` (`created_at`)
) DEFAULT CHARACTER SET=utf8mb4 COMMENT='订阅BI数据响应日志';

-- sub-字段参考表
CREATE TABLE IF NOT EXISTS `sub_refs` (
	`id` bigint unsigned NOT NULL COMMENT '主键ID',
	`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
	`ref_field` varchar(50) NOT NULL DEFAULT '' COMMENT '类型',
	`ref_value` char(1) NOT NULL DEFAULT '' COMMENT '值 A-Z 大写字母',
	`ref_name` varchar(100) NOT NULL DEFAULT '' COMMENT '字段名-中文',
	`ref_name_en` varchar(100) NOT NULL DEFAULT '' COMMENT '字段名-英文',
	`sort` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '排序(越小越靠前)',
	PRIMARY KEY (`id`),
	KEY `idx_reffield_refvalue` (`ref_field`,`ref_value`)
) DEFAULT CHARACTER SET=utf8 COMMENT='sub-字段参考表';

-- 创建操作日志表
CREATE TABLE IF NOT EXISTS `sub_logs_operation` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '操作用户ID',
  `username` varchar(120) NOT NULL DEFAULT '' COMMENT '操作用户名',
  `operation` varchar(50) NOT NULL DEFAULT '' COMMENT '操作类型',
  `resource` varchar(200) NOT NULL DEFAULT '' COMMENT '操作资源',
  `resource_id` varchar(120) DEFAULT '' COMMENT '资源ID',
  `status` varchar(20) NOT NULL DEFAULT '' COMMENT '操作状态',
  `client_ip` varchar(45) NOT NULL DEFAULT '' COMMENT '客户端IP',
  `user_agent` varchar(500) DEFAULT '' COMMENT '用户代理',
  `request_url` varchar(1000) DEFAULT '' COMMENT '请求URL',
  `method` varchar(10) NOT NULL DEFAULT '' COMMENT 'HTTP方法',
  `duration` int unsigned NOT NULL DEFAULT '0' COMMENT '执行耗时(毫秒)',
  `error_msg` text COMMENT '错误信息',
  `request_data` json DEFAULT NULL COMMENT '请求数据',
  `response_data` json DEFAULT NULL COMMENT '响应数据',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_operation` (`operation`),
  KEY `idx_resource` (`resource`),
  KEY `idx_status` (`status`),
  KEY `idx_client_ip` (`client_ip`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 插入测试数据
INSERT INTO `sub_refs` (`id`, `ref_field`, `ref_value`, `ref_name`, `ref_name_en`, `sort`)
VALUES
	('756098401183076352', 'SUBSCRIPTION_TYPE', 'A', '分析数据', 'analysis data', '0'),
	('756098401686392832', 'SUBSCRIPTION_STATUS', 'A', '待生效', 'pending', '1'),
	('756098402231652352', 'SUBSCRIPTION_STATUS', 'B', '生效中', 'activating', '2'),
	('756098402768523264', 'SUBSCRIPTION_STATUS', 'C', '生效中-强制兼容低版本', 'activating force compatible', '3'),
	('756098403276034048', 'SUBSCRIPTION_STATUS', 'D', '已失效', 'expired', '4');

