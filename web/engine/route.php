<?php

function uriParse() {
	$uri_array = explode('/', substr($_SERVER['REQUEST_URI'], 1));
	$c = count($uri_array);
	for ($i=0; $i<$c; $i++) {
		if (trim($uri_array[$i]) == "") {
			unset($uri_array[$i]);
		}
	}
	if (count($uri_array) == 0) {
		return [];
	}
	$default_keys = ['do', 'func', 'some'];

	$ret_array = [];
	$i = 0;
	$c = count($default_keys);
	foreach ($uri_array as $val) {
		$ret_array[$default_keys[$i]] = $val;
		$i++;
		if ($i >= $c) {
			break;
		}
	}
	return $ret_array;
}

function forceHttps() {
	$location = "https://{$_SERVER['HTTP_HOST']}{$_SERVER['REQUEST_URI']}";
	if (array_key_exists("HTTPS", $_SERVER)) {
		if ($_SERVER['HTTPS'] != "on") {
			header("HTTP/1.1 301 Moved Permanently");
			header("Location: {$location}");
			exit;
		}
	} else {
			header("HTTP/1.1 301 Moved Permanently");
			header("Location: {$location}");
			exit;
	}
	return;
}
forceHttps();
$input = uriParse();
if (count($input) != 0) {
	foreach ($input as $key => $val) {
		if (!array_key_exists($key, $e->allow_f)) {
			$e->defaultJsonError();
		}
		if (preg_match("/".$e->allow_f[$key]."/", $val, $matches) === 0) {
			$e->defaultJsonError();
		}
		$input[$key] = $val;
	}
}

//Handling important post parameters
if (isset($_POST['hash'])) {
	if (preg_match("/".$e->allow_f['hash']."/", $_POST['hash'], $matches) === 0) {
			$e->defaultJsonError();
	}
}

//Index
if (count($input) === 0) {
	//default params
	$page = "index";
	$t->get_page("header", "main");
	$t->setVariable("header", "style_etag", $t->get_etag(__DIR__ . "/../stat/css/nknc.css"));
	$t->setVariable("header", "func_etag", $t->get_etag(__DIR__ . "/../stat/js/func.js"));
	$t->get_page($page, "pages");
	$t->get_page("footer", "main");
	print($t->output());
	exit(0);
}

$authByHash = function($hash) use ($e, $t) {
	if (strlen($hash) != 64) {
		$e->defaultJsonError();
	}
	$check = json_decode($e->api_auth($hash, true), true);
	if ($check['Code'] !== 0) {
		print_r($check);
		exit((int)$check['Code']);
	}
	$page = "view";
	$t->get_page("header", "main");
	$t->setVariable("header", "style_etag", $t->get_etag(__DIR__ . "/../stat/css/nknc.css"));
	$t->setVariable("header", "func_etag", $t->get_etag(__DIR__ . "/../stat/js/func.js"));
	$t->get_page($page, "pages");
	$t->setVariable($page, "hash", $hash);
	$t->get_page("footer", "main");
	print($t->output());
	exit(0);
};

if (isset($input['do'])) {
	switch ($input['do']) {
		case "api":
			if (!isset($input['func'])) {
				$e->defaultJsonError();
			}
			$params = [];
			if (!isset($_POST['hash'])) {
				$params[0] = "";
			} else {
				$params[0] = $_POST['hash'];
			}
			die($e->{$input['func']}(...$params));
			break;
		case "login":
			if (!isset($input['func'])) {
				$e->defaultJsonError();
			}
			switch ($input['func']) {
				case "auth":
					if (isset($input['some'])) {
						$authByHash($input['some']);
					}
					$e->defaultJsonError();
			break;
			}
		break;
	}
	$e->defaultJsonError();
}

?>
