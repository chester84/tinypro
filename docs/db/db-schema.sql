DROP TABLE IF EXISTS account_token;
CREATE TABLE `account_token` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'pkid',
  `account_id` bigint(20) NOT NULL,
  `access_token` char(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `token_ip` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `expires` bigint(20) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `platform` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ctime` bigint(20) NOT NULL,
  `utime` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `account_id` (`account_id`),
  KEY `access_token` (`access_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='access token';



