-- create "blocks" table
CREATE TABLE `blocks` (`id` bigint NOT NULL AUTO_INCREMENT, `number` int unsigned NOT NULL, `author_id` varchar(255) NOT NULL, `finalized` bool NOT NULL, `extrinsics_count` bigint NOT NULL, `extrinsics` json NULL, `chain` varchar(255) NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
