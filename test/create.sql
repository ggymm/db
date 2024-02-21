CREATE TABLE `user`
(
    `user_id`      INT NOT NULL,
    `username`     VARCHAR      DEFAULT NULL,
    `nickname`     VARCHAR      DEFAULT NULL,
    `email`        VARCHAR      DEFAULT NULL,
    `phone_number` VARCHAR      DEFAULT NULL,
    `account`      VARCHAR      DEFAULT NULL,
    `password`     VARCHAR      DEFAULT NULL,
    `status`       INT          DEFAULT NULL,
    `extras`       VARCHAR      DEFAULT NULL,
    `create_time`  VARCHAR      DEFAULT NULL,
    `create_id`    INT NOT NULL DEFAULT '1',
    `update_time`  VARCHAR      DEFAULT NULL,
    `update_id`    INT NOT NULL DEFAULT '1',
    `del_flag`     INT NOT NULL DEFAULT '1',
    PRIMARY KEY (`user_id`),
    INDEX          account_index (`account`)
);