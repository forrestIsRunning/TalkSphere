CREATE TABLE `user` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `user_id` bigint(20) NOT NULL,
    `username` varchar(64) NOT NULL,
    `password` varchar(64) NOT NULL,
    `email` varchar(128) DEFAULT NULL,
    `avatar` varchar(255) DEFAULT NULL,
    `bio` varchar(500) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_email` (`email`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;