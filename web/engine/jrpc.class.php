<?php
class jrpc {
	const VERSION = '2.0';
	protected $url;
	private $method;
	private $params;

	public function __construct(string $url) {
		$this->url = $url;
	}

	public function __call(string $name, array $args = []) {
		$this->method = $name;
		$this->params = count($args) == 1 && is_array($args[0]) ? $args[0] : $args;
		return $this->_curl();
	}

	private function _curl()
	{
		$curl = curl_init();
		curl_setopt($curl, CURLOPT_URL, $this->url);
		curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($curl, CURLOPT_POST, true);
		curl_setopt($curl, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);
		curl_setopt($curl, CURLOPT_POSTFIELDS, json_encode([
			'method' => $this->method,
			'params' => $this->params,
			'id' => microtime(),
			'jsonrpc' => self::VERSION
		]));
		$result = curl_exec($curl);

		curl_close($curl);

		return new jrpc_response($result);
	}
}
class jrpc_response {
	protected $response;
	protected $success;
	public function __construct($response){
		$this->response = json_decode($response, false, 512, JSON_THROW_ON_ERROR);
		$this->success = !isset($this->response->error);
	}
	public function isSuccess() {
		return $this->success;
	}
	public function getResult() {
		return $this->response->result;
	}
	public function getErrorCode() {
		return $this->response->error->code;
	}
	public function getErrorMessage() {
		return $this->response->error->message;
	}
}