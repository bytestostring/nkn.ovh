CREATE TABLE `all_nodes` (
  `id` int(10) UNSIGNED NOT NULL,
  `ip` varchar(15) NOT NULL,
  `addr` varchar(32) NOT NULL,
  `node_id` varchar(65) DEFAULT NULL,
  `syncState` varchar(32) NOT NULL,
  `uptime` int(11) DEFAULT NULL,
  `proposalSubmitted` int(11) DEFAULT NULL,
  `relayMessageCount` int(11) DEFAULT NULL,
  `height` int(11) NOT NULL,
  `currtimestamp` bigint(20) UNSIGNED DEFAULT NULL,
  `version` varchar(64) DEFAULT NULL,
  `latest_update` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `all_nodes_stats` (
  `id` int(10) UNSIGNED NOT NULL,
  `relays` bigint(20) NOT NULL,
  `average_uptime` int(11) NOT NULL,
  `average_relays` int(11) NOT NULL,
  `relays_per_hour` bigint(20) UNSIGNED NOT NULL,
  `proposalSubmitted` int(10) UNSIGNED NOT NULL,
  `persist_nodes_count` int(10) UNSIGNED NOT NULL,
  `nodes_count` int(10) UNSIGNED NOT NULL,
  `last_height` int(11) NOT NULL,
  `last_timestamp` bigint(20) NOT NULL,
  `average_blockTime` float NOT NULL,
  `average_blocksPerDay` float NOT NULL,
  `latest_update` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `daemon` (
  `id` int(10) UNSIGNED NOT NULL,
  `name` varchar(64) NOT NULL,
  `value` varchar(1024) NOT NULL,
  `updated` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `nodes` (
  `id` int(10) UNSIGNED NOT NULL,
  `hash_id` int(10) UNSIGNED NOT NULL,
  `name` varchar(64) NOT NULL,
  `ip` varchar(20) NOT NULL,
  `created` datetime NOT NULL DEFAULT current_timestamp(),
  `dirty` tinyint(1) NOT NULL DEFAULT 1,
  `dirty_fcnt` int(10) UNSIGNED NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `nodes_history` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `node_id` int(10) UNSIGNED NOT NULL,
  `NID` varchar(128) NOT NULL,
  `Currtimestamp` bigint(10) UNSIGNED NOT NULL,
  `Height` int(10) UNSIGNED NOT NULL,
  `ProposalSubmitted` int(10) UNSIGNED NOT NULL,
  `ProtocolVersion` int(11) NOT NULL,
  `RelayMessageCount` bigint(20) UNSIGNED NOT NULL,
  `SyncState` varchar(64) NOT NULL,
  `Uptime` int(11) UNSIGNED NOT NULL,
  `Version` varchar(64) NOT NULL,
  `latest_update` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `nodes_last` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `node_id` int(10) UNSIGNED NOT NULL,
  `NID` varchar(128) NOT NULL,
  `Currtimestamp` bigint(10) UNSIGNED NOT NULL,
  `Height` int(10) UNSIGNED NOT NULL,
  `ProposalSubmitted` int(10) UNSIGNED NOT NULL,
  `ProtocolVersion` int(11) NOT NULL,
  `RelayMessageCount` bigint(20) UNSIGNED NOT NULL,
  `SyncState` varchar(64) NOT NULL,
  `Uptime` int(11) UNSIGNED NOT NULL,
  `Version` varchar(64) NOT NULL,
  `failcnt` int(11) NOT NULL DEFAULT 0,
  `firsttime_failed` tinyint(1) NOT NULL DEFAULT 0,
  `latest_update` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `prices` (
  `id` int(10) UNSIGNED NOT NULL,
  `name` varchar(32) CHARACTER SET latin1 NOT NULL,
  `price` decimal(16,8) NOT NULL,
  `last_update` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `uniq` (
  `id` int(10) UNSIGNED NOT NULL,
  `hash` varchar(64) NOT NULL,
  `pass` varchar(64) DEFAULT NULL,
  `ip_creator` varchar(65) NOT NULL,
  `created_by` datetime NOT NULL DEFAULT current_timestamp(),
  `latest_watch` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `wallets` (
  `id` int(10) UNSIGNED NOT NULL,
  `hash_id` int(10) UNSIGNED NOT NULL,
  `nkn_wallet` varchar(48) NOT NULL,
  `balance` decimal(16,8) NOT NULL,
  `created` datetime NOT NULL DEFAULT current_timestamp(),
  `updated` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `watchers` (
  `id` int(10) UNSIGNED NOT NULL,
  `hash_id` int(10) UNSIGNED NOT NULL,
  `ro_hash` varchar(64) NOT NULL,
  `hide_names` tinyint(1) NOT NULL,
  `hide_ip` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


ALTER TABLE `all_nodes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `addr` (`addr`),
  ADD KEY `syncState` (`syncState`),
  ADD KEY `ip` (`ip`);

ALTER TABLE `all_nodes_stats`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `daemon`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

ALTER TABLE `nodes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `hash_id` (`hash_id`,`ip`),
  ADD KEY `dirty` (`dirty`),
  ADD KEY `dirty_fcnt` (`dirty_fcnt`);

ALTER TABLE `nodes_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `node_id` (`node_id`);

ALTER TABLE `nodes_last`
  ADD PRIMARY KEY (`id`),
  ADD KEY `node_id` (`node_id`),
  ADD KEY `failcnt` (`failcnt`),
  ADD KEY `firsttime_failed` (`firsttime_failed`);

ALTER TABLE `prices`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

ALTER TABLE `uniq`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `hash` (`hash`),
  ADD KEY `ip_creator` (`ip_creator`),
  ADD KEY `created_by` (`created_by`);

ALTER TABLE `wallets`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `hash_id_2` (`hash_id`,`nkn_wallet`),
  ADD KEY `hash_id` (`hash_id`);

ALTER TABLE `watchers`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `hash_id` (`hash_id`),
  ADD UNIQUE KEY `ro_hash` (`ro_hash`);


ALTER TABLE `all_nodes`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `all_nodes_stats`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `daemon`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `nodes`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `nodes_history`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `nodes_last`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `prices`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `uniq`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `wallets`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

ALTER TABLE `watchers`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;


ALTER TABLE `nodes`
  ADD CONSTRAINT `nodes_ibfk_1` FOREIGN KEY (`hash_id`) REFERENCES `uniq` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE `nodes_history`
  ADD CONSTRAINT `nodes_history_ibfk_1` FOREIGN KEY (`node_id`) REFERENCES `nodes` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE `nodes_last`
  ADD CONSTRAINT `nodes_last_ibfk_1` FOREIGN KEY (`node_id`) REFERENCES `nodes` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE `wallets`
  ADD CONSTRAINT `wallets_ibfk_1` FOREIGN KEY (`hash_id`) REFERENCES `uniq` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE `watchers`
  ADD CONSTRAINT `watchers_ibfk_1` FOREIGN KEY (`hash_id`) REFERENCES `uniq` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION;


