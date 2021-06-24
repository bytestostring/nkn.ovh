<?php

class Engine {

	private $db;
	public $stmt;
	public $allow_f;
	private $hash_id;
	private $method;
	private $params;

	function __construct() {
		$this->allow_f = ['do' => '^(api|login|id)$', 'func' => '^(genId|auth|getnetstatus|getmynodesstat|addnodes|rmnodes|getmywallets|getfullstack|savemysettings)$', 'hash' => '^([A-Za-z0-9]{64}+)$', 'some' => '^([A-Za-z0-9]+)$', 'watcher' => '^([A-Za-z0-9]+)$'];
	}

	
	public function __call(string $name, array $args = []) {
		$this->method = "api_{$name}";
		//$params = count($args) == 1 && is_array($args[0]) ? $args[0] : $args;
		if (!method_exists($this, $this->method)) {
			return false;
		}
		return $this->{$this->method}(...$args);
	}


	public function DB_connect() {
		require('db_config.php');
		try {
			$dbh = new PDO("mysql:host={$db['host']};dbname={$db['db']};charset={$db['encoding']}", $db['login'], $db['password']);
			$dbh->setAttribute(PDO::ATTR_EMULATE_PREPARES, false);
			$this->db = $dbh;
			$this->db_prepare();
		} catch (PDOException $e) {
			print("Cannot connect to a database. Sorry, please try again later");
			$this->syslog($e->getMessage());
			die();
		}
		return true;
	}

	private function db_prepare() {
		try {
			$stmt = [
					 'CheckIPCreator' => 'SELECT id FROM uniq WHERE ip_creator = ? AND created_by >= NOW() - INTERVAL 30 MINUTE',
					 'CreateUniq' => 'INSERT INTO uniq(hash,ip_creator) VALUES(?,?)',
					 'SelectUniqByHash' => 'SELECT * FROM uniq WHERE hash = ?',
					 'UpdateUniqWatch' => 'UPDATE uniq SET latest_watch=CURRENT_TIMESTAMP() WHERE id = ?',
					 'getNetStatus' => 'SELECT * FROM all_nodes_stats WHERE id=(SELECT max(id) FROM all_nodes_stats)',
					 'getMyNodes' => 'SELECT * FROM nodes WHERE hash_id = ?',
					 'insertNode' => 'INSERT IGNORE INTO nodes(hash_id,name,ip) VALUES(?,?,?)',
					 'countNodesByHash' => 'SELECT count(id) as cnt FROM nodes WHERE hash_id = ?',
					 'rmNodes' => 'DELETE FROM nodes WHERE hash_id = ? && id = ?',
					 'getMyWallets' => 'SELECT * FROM wallets WHERE hash_id = ? ORDER BY id ASC',
					 'getWalletByAddress' => 'SELECT id FROM wallets WHERE hash_id = ? AND nkn_wallet = ?',
					 'rmAllWalletsByHash' => 'DELETE FROM wallets WHERE hash_id = ?',
					 'addWallet' => 'INSERT INTO wallets(hash_id,nkn_wallet,balance) VALUES(?,?,-100)',
					 'getPrices' => 'SELECT * FROM prices',
					 'getDaemon' => 'SELECT * FROM daemon',
					 'getMyNodeLastInfo' => 'SELECT * FROM nodes_last WHERE node_id = ?'
					];
			foreach ($stmt as $key => $val) {
				$this->stmt[$key] = $this->db->prepare($val); 
			}
			$this->stmt = (object) $this->stmt;
		} catch (PDOException $e) {
			print("Cannot prepare queries. Sorry, please try again later");
			$this->syslog($e->getMessage());
			die();
		}
		return true;
	}

	public function defaultJsonError() {
		http_response_code(500);
		print(json_encode(["Err" => true, "Value" => "Bad parameters"]));
		die();
		return;
	}

	private function syslog(string $t, string $logfile = "") {
		$dir = "logs";
		if ($logfile == "") {
			$logfile = "main.log";
		}
		$info = "[".date('d/m/Y:H:i:s.').substr(explode(".", explode(" ", microtime())[0])[1], 0, 3)."] [".getmypid()."] {$t}";
		$deb = debug_backtrace();
		$info .= " [ {$deb[0]['file']}:{$deb[0]['line']} ({$deb[1]['function']}) ]\n";
		if (!is_dir($dir)) {
			mkdir($dir, 0755);
		}
		$f = fopen("{$dir}/$logfile", 'a+');
		flock($f, LOCK_EX);
		fwrite($f, $info);
		flock($f, LOCK_UN);
		fclose($f);
		return true;
	}

	private function bind_sql(string $prepare, array $params) {
		$n = 0; $t = 1;
		foreach ($params[0] as $key => $value):
			$type = $params[1][$n];
			if ($this->stmt->$prepare->bindValue($t, $value, $type) === false) {
				$this->syslog("BindValue error:");
				return false;
			}
			$n++;
			$t++;
		endforeach;
		return true;
	}

	public function api_genId() {

		$pre_work = &$this->stmt->CheckIPCreator;
		$ip = $_SERVER['REMOTE_ADDR'];
		$params = [ [ $ip ], [PDO::PARAM_STR]];
		if ($this->bind_sql('CheckIPCreator', $params) === false) {
			$this->syslog("bind_sql failed");
			return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
		}
		try {
			$pre_work->execute();
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
		}

		if ($pre_work->rowCount() >= 3) {
			return json_encode(["Code" => 3, "Err" => true, "Value" => "You have created at least 3 ID for the latest 30 minutes"]);
		}
		$gen = hash('sha256', random_bytes(256));
		$work = &$this->stmt->CreateUniq;
		$params = [ [ $gen, $ip ], [PDO::PARAM_STR, PDO::PARAM_STR]];
		if ($this->bind_sql('CreateUniq', $params) === false) {
			$this->syslog("bind_sql failed");
			return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
		}
		try {
			$work->execute();
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
		}
		//fetch params
		return json_encode(["Code" => 0, "Value" => $gen]);
	}


	public function api_auth(string $hash, bool $check_pass = false) {
		if (strlen($hash) != 64) {
			return json_encode(["Code" => 254, "Err" => true, "Value" => "Incorrect ID length"]);
		}
		$work = &$this->stmt->SelectUniqByHash;
		$params = [ [ $hash ], [ PDO::PARAM_STR ] ];
			if ($this->bind_sql('SelectUniqByHash', $params) === false) {
			$this->syslog("bind_sql failed");
			return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
		}
		try {
			$work->execute();
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
		}
		if ($work->rowCount() == 0) {
			return json_encode(["Code" => 255, "Err" => true, "Value" => "Passed ID is not found"]);
		}
		$row = $work->fetch(PDO::FETCH_ASSOC);
		if (is_null($row['pass'])) {
			$up = $this->UpdateUniqWatch($row['id']);
			if ($up !== true) {
				return $up;
			}
		}

		$this->hash_id = $row['id'];

		if ($check_pass === false) {
			return json_encode(["Code" => 0, "Value" => $row['hash']]);
		}
		if (!is_null($row['pass'])) {
			return json_encode(["Code" => 10, "Value" => "Need password"]);
		}


		return json_encode(["Code" => 0, "Value" => $row['hash']]);
	}

	public function api_rmnodes(string $hash) {
		$a = json_decode($this->api_auth($hash, false), true);
		if ($a['Code'] !== 0) {
			return json_encode($a);
		}
		if (!isset($_POST['ids'])) {
			return json_encode(["Code" => 3, "Err" => true, "Value" => "Incorrect POST query"]);
		}
		$work = &$this->stmt->rmNodes;
		foreach ($_POST['ids'] as $val) {
			if (!is_numeric($val)) {
				continue;
			}
			$params = [ [ $this->hash_id, $val ], [ PDO::PARAM_INT, PDO::PARAM_INT] ];
			if ($this->bind_sql('rmNodes', $params) === false) {
				$this->syslog("bind_sql failed");
				return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
			}
			try {
				$work->execute();
			} catch (Exception $ex) {
				$this->syslog("Cannot execute query: ".$ex->getMessage());
				return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);			
			}
		}
		return json_encode(["Code" => 0, "Value" => "Your nodes removed"]);
	}

	public function api_savemysettings(string $hash) {
		$a = json_decode($this->api_auth($hash, false), true);
		if ($a['Code'] !== 0) {
			return json_encode($a);
		}

		$work = &$this->stmt->getMyWallets;
		$params = [ [ $this->hash_id ], [ PDO::PARAM_INT ] ];

		if ($this->bind_sql('getMyWallets', $params) === false) {
			$this->syslog("bind_sql failed");
			return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
		}

		try {
			$work->execute();
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);			
		}
		$db_wallets = [];
		if ($work->rowCount() != 0) {
			while ($w = $work->fetch(PDO::FETCH_ASSOC)) {
				$db_wallets[] = $w['nkn_wallet'];
			}
		}

		//Parse wallets
		$wallets = [];
		if (isset($_POST['wallets'])) {
			$c = is_array($_POST['wallets']) ? count($_POST['wallets']): 0;
		} else {
			$c = 0;
		}
		if ($c > 3) {
			return json_encode(['Code' => 4, 'Err' => true, 'Value' => "Incorrect array length!"]);
		}
		if ($c > 0) {
			foreach (array_unique($_POST['wallets']) as $val) {
				$address = trim(strip_tags($val));
				if (strlen($address) != 36) {
					return json_encode(['Code' => 5, 'Err' => true, 'Value' => "Incorrect address length!"]);
				}
				if (preg_match('/^NKN([A-Za-z0-9]{33}+)$/', $address) !== 1) {
					return json_encode(['Code' => 6, 'Err' => true, 'Value' => "Incorrect NKN address"]);
				}
				$wallets[] = $address;
			}
			$cnt_dbw = count($db_wallets);
			$check = $cnt_dbw;
			$new_found = false;
			foreach ($wallets as $val) {
				if ($cnt_dbw > 0) {
					if (in_array($val, $db_wallets)) {
						continue;
					}
				}
				$new_found = true;
				break;
			}

			//Remove old and insert new
			if ($new_found == true || $cnt_dbw > $c) {
				$work = &$this->stmt->rmAllWalletsByHash;
				$params = [ [ $this->hash_id ], [ PDO::PARAM_INT ] ];

				if ($this->bind_sql('rmAllWalletsByHash', $params) === false) {
					$this->syslog("bind_sql failed");
					return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
				}
				try {
					$work->execute();
				} catch (Exception $ex) {
					$this->syslog("Cannot execute query: ".$ex->getMessage());
					return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);	
				}
				$work = &$this->stmt->addWallet;
				foreach ($wallets as $val) {
					$params = [ [ $this->hash_id, $val], [ PDO::PARAM_INT, PDO::PARAM_STR ] ];
					if ($this->bind_sql('addWallet', $params) === false) {
						$this->syslog("bind_sql failed");
						return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
					}
					try {
						$work->execute();
					} catch (Exception $ex) {
						$this->syslog("Cannot execute query: ".$ex->getMessage());
						return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);	
					}
				}
			}
		} elseif ($c == 0 && count($db_wallets) != 0) {
			$work = &$this->stmt->rmAllWalletsByHash;
			$params = [ [ $this->hash_id ], [ PDO::PARAM_INT ] ];

			if ($this->bind_sql('rmAllWalletsByHash', $params) === false) {
				$this->syslog("bind_sql failed");
				return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
			}
			try {
				$work->execute();
			} catch (Exception $ex) {
				$this->syslog("Cannot execute query: ".$ex->getMessage());
				return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);	
			}
		}

		return json_encode(['Code' => 0, 'Value' => "Settings saved"]);
	}

	public function api_getmywallets(string $hash, bool $ignore_auth = false) {
		if ($ignore_auth === false) {
			$a = json_decode($this->api_auth($hash, false), true);
			if ($a['Code'] !== 0) {
				return json_encode($a);
			}
		}
		$work = &$this->stmt->getMyWallets;
		$params = [ [ $this->hash_id ], [ PDO::PARAM_INT ] ];

		if ($this->bind_sql('getMyWallets', $params) === false) {
			$this->syslog("bind_sql failed");
			return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
		}

		try {
			$work->execute();
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);			
		}
		if ($work->rowCount() == 0) {
			return json_encode(["Code" => 3, "Value" => "No wallets are found"]);
		}
		$wallets = [];

		while ($w = $work->fetch(PDO::FETCH_ASSOC)) {
			$wallets[] = $w;
		}
		return json_encode(["Code" => 0, "Value" => $wallets]);
	}

	public function api_getmynodesstat($hash, $ignore_auth = false) {
		if ($ignore_auth === false) {
			$a = json_decode($this->api_auth($hash, false), true);
			if ($a['Code'] !== 0) {
				return json_encode($a);
			}
		}
		$work = &$this->stmt->getMyNodes;
		$params = [ [ $this->hash_id ], [ PDO::PARAM_INT ] ];
		if ($this->bind_sql('getMyNodes', $params) === false) {
			$this->syslog("bind_sql failed");
			return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
		}
		try {
			$work->execute();
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);			
		}
		if ($work->rowCount() == 0) {
			return json_encode(["Code" => 3, "Value" => "No nodes are found"]);
		}
		$nodes = [];
		$n = -1;
		$tmp = [];
		$nodes_id = [];
		$no_history_nodes = [];
		$c = [];
		while ($node = $work->fetch(PDO::FETCH_ASSOC)) {
			$n++;
			$nodes[$n] = ['ip' => $node['ip'], 'name' => $node['name']];
			$nodes_id[$node['id']] = $n;

			$c[$n] = 0;
			$params = [ [ $node['id'] ], [ PDO::PARAM_INT ] ];

			$work_last = &$this->stmt->getMyNodeLastInfo;
			if ($this->bind_sql('getMyNodeLastInfo', $params) === false) {
				$this->syslog("bind_sql failed");
				return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
			}
			try {
				$work_last->execute();
			} catch (Exception $ex) {
				$this->syslog("Cannot execute query: ".$ex->getMessage());
				return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
			}
			if ($work_last->rowCount() == 0) {
				$nodes[$n]['Err'] = 2;
				$nodes[$n]['node_id'] = $node['id'];
				$nodes[$n]['latest_update'] = date('Y-m-d H:i:s');
				$nodes[$n]['SyncState'] = "Waiting for first update";
				$no_history_nodes[$n] = 1;
				continue;
			}
			$row = $work_last->fetch(PDO::FETCH_ASSOC);
			$nodes[$n]['latest_update'] = $row['latest_update'];
			$nodes[$n]['Uptime'] = $row['Uptime'];
			$nodes[$n]['SyncState'] = $row['SyncState'];
			$nodes[$n]['RelayMessageCount'] = $row['RelayMessageCount'];
			$nodes[$n]['Currtimestamp'] = $row['Currtimestamp'];				
			$nodes[$n]['ProposalSubmitted'] = $row['ProposalSubmitted'];
			if ($row['Uptime'] != 0) {
				$nodes[$n]['RelaysPerHour'] = floor($row['RelayMessageCount']/$row['Uptime']*3600);
			} else {
				$nodes[$n]['RelaysPerHour'] = 0;	
			}
			$nodes[$n]['Height'] = $row['Height'];
			$nodes[$n]['Version'] = strip_tags($row['Version']);
			$nodes[$n]['node_id'] = $row['node_id'];

			if ($row['SyncState'] == "OFFLINE") {
				$nodes[$n]['Err'] = 1;
				$nodes[$n]['SyncState'] = "_OFFLINE_";
				$no_history_nodes[$n] = 1;
				continue;
			} elseif ($row['SyncState'] == "PRUNING DB" || $row['SyncState'] == "GENERATING ID") {
				$no_history_nodes[$n] = 1;
			}
		}
		$in = str_repeat('?,', count($nodes_id) - 1) . '?';
		$sqlHistory = "SELECT * FROM nodes_history WHERE node_id in ($in) ORDER BY id DESC";
		$local_stmt = $this->db->prepare($sqlHistory);
		$keys = array_keys($nodes_id);
		try {
			$local_stmt->execute($keys);
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
		}

		if ($local_stmt->rowCount() != 0) {
			while ($row = $local_stmt->fetch(PDO::FETCH_ASSOC)) {
				$n = $nodes_id[$row['node_id']];
				if (isset($no_history_nodes[$n])) {
					continue;
				}
				$tmp[$n][$c[$n]] = ["Uptime" => $row['Uptime'], "Currtimestamp" => $row['Currtimestamp'], "Relays" => $row['RelayMessageCount']];
					$c[$n]++;
			}
		}

		$opts = [
					["rh" => 600, "max_diff" => 60, "elem" => 1],
					["rh" => 3600, "max_diff" => 60, "elem" => 6],
	
				];
		foreach ($tmp as $key => $val) {
				foreach ($opts as $p) {
					$this->setRelays($val, $p['rh'], $p['max_diff'], $p['elem'], $nodes, $key);
				}
			}
		$retjson = ['Code' => 0, 'Value' => $nodes];
		return json_encode($retjson, true);
	}

	private function setRelays($val, $rh, $max_diff, $elem, &$nodes, $key) {
		switch ($rh) {
			case 600:
			$t = "RelaysPerHour10";
			break;
			case 3600:
			$t = "RelaysPerHour60";
			break;
			default:
			return false;
		}
		$i = 0;
		if (array_key_exists($i+$elem, $val) === true) {
			$diff_uptime = $rh - ($val[$i]['Uptime'] - $val[$i+$elem]['Uptime']);
				$diff_timestamp = $rh - ($val[$i]['Currtimestamp'] - $val[$i+$elem]['Currtimestamp']);
				$diff_general = $diff_timestamp-$diff_uptime;
				if (($diff_general >= 0 && $diff_general <= 10) XOR ($diff_general < 0 && $diff_general >= -10)) {
					if (($diff_uptime <= $max_diff && $diff_uptime >= 0 ) XOR ($diff_uptime < $max_diff && $diff_uptime < 0 && ($max_diff+$diff_uptime) > 0)) {
						if (($val[$i]['Uptime'] > $val[$i+$elem]['Uptime']) && ($val[$i+$elem]['Uptime'] > 0)) {
							$cf = $val[$i]['Uptime']/$val[$i+$elem]['Uptime'];
							$nodes[$key][$t] = ($val[$i]['Relays'] - $val[$i+$elem]['Relays'])/$rh*3600;
							return true;
						} else {
							$nodes[$key][$t] = -1;
							return true;
						}
					} else {
						$nodes[$key][$t] = -1;
						return true;
					}
				}
			}
			$nodes[$key][$t] = -1;
			return;
	}


	private function getPrices() {
		$work = &$this->stmt->getPrices;
		try {
			$work->execute();
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);			
		}
		if ($work->rowCount() > 0) {
			$prices = [];
			while ($r = $work->fetch(PDO::FETCH_ASSOC)) {
				$prices[$r['name']] = $r['price'];
			}
			return json_encode(["Code" => 0, "Value" => $prices]);
		}
		return json_encode(["Code" => 3, "Value" => "No prices found"]);
	}


	public function api_getfullstack(string $hash) {
		$a = json_decode($this->api_auth($hash, false), true);
		if ($a['Code'] !== 0) {
			return json_encode($a);
		}
		$netstatus = json_decode($this->api_getnetstatus($hash, true), true);
		$nodes = json_decode($this->api_getmynodesstat($hash, true), true);
		$wallets = json_decode($this->api_getmywallets($hash, true), true);
		$prices = json_decode($this->getPrices(), true);
		$daemonInfo = json_decode($this->getDaemon(), true);
		if ($netstatus['Code'] !== 0 || ($nodes['Code'] !== 0 && $nodes['Code'] !== 3) || ($wallets['Code'] !== 0 && $wallets['Code'] !== 3) || $prices['Code'] !== 0) {
			return json_encode(['Code' => 4, 'Value' => 'An error occured']);
		}
		return json_encode(['Code' => 0, 'NETSTATUS' => $netstatus, 'NODES' => $nodes, 'WALLETS' => $wallets, 'PRICES' => $prices, 'DAEMON' => $daemonInfo]);
	}


	private function getDaemon() {
		$work = &$this->stmt->getDaemon;
		try {
			$work->execute();
		} catch(Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
		}
		if ($work->rowCount() == 0) {
			return json_encode(["Code" => 3, "Err" => true, "Value" => "No info in a database"]);
		}
		$info = [];
		while ($r = $work->fetch(PDO::FETCH_ASSOC)) {
			$info[$r['name']] = $r['value'];
		}
		return json_encode(["Code" => 0, "Value" => $info]);
	}

	public function api_getnetstatus(string $hash, bool $ignore_auth = false) {
		if ($ignore_auth === false) {
			$a = json_decode($this->api_auth($hash, false), true);
			if ($a['Code'] !== 0) {
				return json_encode($a);
			}
		}
		$work = &$this->stmt->getNetStatus;
		try {
			$work->execute();
		} catch(Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
		}
		if ($work->rowCount() != 1) {
			return json_encode(["Code" => 2, "Err" => true, "Value" => "No info in a database"]);
		}
		$r = $work->fetch(PDO::FETCH_ASSOC);
		unset($r['id']);
		return json_encode(["Code" => 0, "Value" => $r]);
	}

	public function api_addnodes(string $hash) {
		$a = json_decode($this->api_auth($hash, false), true);
		if ($a['Code'] !== 0) {
			return json_encode($a);
		}
		if (!isset($_POST['IP']) || !isset($_POST['Name']) || !isset($_POST['multiple']))  {
			return json_encode(["Code" => 3, "Err" => true, "Value" => "Incorrect POST query"]);
		}
		$ips = trim($_POST['IP']);
		$name = trim($_POST['Name']);
		$multiple = $_POST['multiple'];

		if (strlen($name) > 0) {
			if (strlen($name) > 32) {
				return json_encode(["Code" => 6, "Err" => true, "Value" => "The name too long"]);
			}
			$name = filter_var($name, FILTER_SANITIZE_STRING);
			if ($name === false) {
				$name = "";
			}
		}
		$work = &$this->stmt->insertNode;

		if ($multiple != 0 && $multiple != 1) {
			return json_encode(["Code" => 3, "Err" => true, "Value" => "Incorrect POST query"]);
		}
		if ($multiple == 0) {
			$ip = filter_var($ips, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4 | FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE);
			if ($ip === false) {
				return json_encode(["Code" => 4, "Err" => true, "Value" => "Incorrect IP address"]);
			}
			if (strlen($name) == 0) {
				$work2 = &$this->stmt->countNodesByHash;
				$params = [ [ $this->hash_id], [ PDO::PARAM_INT] ];
				if ($this->bind_sql('countNodesByHash', $params) === false) {
					$this->syslog("bind_sql failed");
					return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
				}
				try {
					$work2->execute();
				} catch (Exception $ex) {
					$this->syslog("Cannot execute query:". $ex->getMessage());
					return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
				}
				$c = $work2->fetch(PDO::FETCH_ASSOC)['cnt']+1;
				$name = "Node_{$c}";
			}
			$params = [ [ $this->hash_id, $name, $ip], [ PDO::PARAM_INT, PDO::PARAM_STR, PDO::PARAM_STR] ];
			if ($this->bind_sql('insertNode', $params) === false) {
				$this->syslog("bind_sql failed");
				return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
			}
			try {
				$work->execute();
			} catch (Exception $ex) {
				$this->syslog("Cannot execute query:". $ex->getMessage());
				return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
			}
			return json_encode(["Code" => 0, "Value" => "Node added"]);
		}
		
		if ($multiple == 1) {
			$ip_array = [];
			if (strpos($ips, ",") !== false) {
				$ex = explode(",", $ips);
			} elseif (strpos($ips, " ")) {
				$ex = explode(" ", $ips);
			} elseif (strpos($ips, "\n")) {
				$ex = explode("\n", $ips);
			} else {
				return json_encode(["Code" => 5, "Err" => true, "Value" => "Incorrect delimiter."]);
			}
			$cnt_nodes = count($ex);
			if ($cnt_nodes <= 1) {
				return json_encode(["Code" => 5, "Err" => true, "Value" => "Incorrect IP addresses. You should use a single form for one address."]);
			}

			if ($cnt_nodes > 5000) {
				return json_encode(["Code" => 7, "Err" => true, "Value" => "Cannot add over 5000 nodes per your id"]);
			}

			$work2 = &$this->stmt->countNodesByHash;
			$params = [ [ $this->hash_id], [ PDO::PARAM_INT] ];
			if ($this->bind_sql('countNodesByHash', $params) === false) {
				$this->syslog("bind_sql failed");
				return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
			}
			try {
				$work2->execute();
			} catch (Exception $ex) {
				$this->syslog("Cannot execute query:". $ex->getMessage());
				return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
			}
			$c = $work2->fetch(PDO::FETCH_ASSOC)['cnt']+1;

			if (($cnt_nodes+$c-1) > 5000) {
				return json_encode(["Code" => 7, "Err" => true, "Value" => "Cannot add over 5000 nodes per your id"]);
			} 

			foreach ($ex as $val) {
				$ip = filter_var(trim($val), FILTER_VALIDATE_IP, FILTER_FLAG_IPV4 | FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE);
				if ($ip === false) {
					return json_encode(["Code" => 4, "Err" => true, "Value" => "Incorrect IP address/addresses in your list"]);
				}
				if (array_search($ip, $ip_array) === false) {
					$ip_array[] = $ip;
				}
			}	

			if (strlen($name) == 0) {
				$name = "Node";
			}
			foreach ($ip_array as $ip) {
				$params = [ [ $this->hash_id, "${name}_{$c}", $ip], [ PDO::PARAM_INT, PDO::PARAM_STR, PDO::PARAM_STR] ];
				if ($this->bind_sql('insertNode', $params) === false) {
					$this->syslog("bind_sql failed");
					return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
				}
				try {
					$work->execute();
				} catch (Exception $ex) {
					$this->syslog("Cannot execute query:". $ex->getMessage());
					return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]); 
				}
				$c++;
			}
			return json_encode(["Code" => 0, "Value" => "All nodes added"]);
		}
	}

	private function UpdateUniqWatch(int $id) {
		$work = &$this->stmt->UpdateUniqWatch;
		$params = [ [ $id ], [ PDO::PARAM_INT ] ];
		if ($this->bind_sql('UpdateUniqWatch', $params) === false) {
			$this->syslog("bind_sql failed");
			return json_encode(["Code" => 2, "Err" => true, "Value" => "Cannot bind params"]);
		}
		try {
			$work->execute();
		} catch (Exception $ex) {
			$this->syslog("Cannot execute query: ".$ex->getMessage());
			return json_encode(["Code" => 1, "Err" => true, "Value" => "Cannot execute an query"]);
		}
		return true;
	}
}
?>
