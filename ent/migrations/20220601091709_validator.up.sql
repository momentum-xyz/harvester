-- create "validators" table
CREATE TABLE `validators` (`id` bigint NOT NULL AUTO_INCREMENT, `account_id` varchar(255) NOT NULL, `name` varchar(255) NOT NULL, `commission` double NOT NULL, `status` varchar(255) NOT NULL, `balance` varchar(255) NOT NULL, `reserved` varchar(255) NOT NULL, `locked` json NOT NULL, `own_stake` varchar(255) NOT NULL, `total_stake` varchar(255) NOT NULL, `identity` json NOT NULL, `nominators` json NOT NULL, `parent` json NOT NULL, `children` json NOT NULL, `hash` varchar(255) NOT NULL, `chain` varchar(255) NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
