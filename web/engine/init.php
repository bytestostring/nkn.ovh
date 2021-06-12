<?php
session_start();
require('engine.class.php');
require('templater.class.php');
require('config.php');
require('jrpc.class.php');
$e = new Engine();
$t = new templater;
$e->DB_CONNECT();

require('route.php');
?>