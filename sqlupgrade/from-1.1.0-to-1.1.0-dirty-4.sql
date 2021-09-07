ALTER TABLE `all_nodes` CHANGE `node_id` `NID` VARCHAR(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL;
ALTER TABLE `all_nodes_last` CHANGE `node_id` `NID` VARCHAR(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL;
ALTER TABLE `all_nodes_last` ADD `PublicKey` VARCHAR(64) NULL AFTER `NID`;
ALTER TABLE `all_nodes` ADD `PublicKey` VARCHAR(64) NULL AFTER `NID`;
ALTER TABLE `all_nodes_last` ADD INDEX(`NID`);
ALTER TABLE `all_nodes_last` ADD INDEX(`PublicKey`);
