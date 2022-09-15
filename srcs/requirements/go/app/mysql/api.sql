USE `api`;
CREATE TABLE IF NOT EXISTS `api`(
	`id` int unsigned NOT NULL AUTO_INCREMENT,
	`api` text,
	`api_secret` text,
	PRIMARY_KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;
