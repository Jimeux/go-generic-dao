CREATE TABLE IF NOT EXISTS `user`
(
    `id`         BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `nickname`   VARCHAR(255)    NOT NULL,
    `bio`        TEXT            NULL,
    `created_at` TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS `like`
(
    `id`         BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_id`    BIGINT UNSIGNED NOT NULL REFERENCES `user` (`id`),
    `partner_id` BIGINT UNSIGNED NOT NULL REFERENCES `user` (`id`),
    `created_at` TIMESTAMP       NOT NULL DEFAULT NOW(),
    INDEX (`user_id`, `partner_id`)
);
