<?php
class templater {

	private $catalog;
	private $pages;
	function __construct() {
		$this->catalog = "templates";
		$this->pages = [];
	}
	
	public function get_page(string $page, string $inside_path) {
		$file = "{$this->catalog}/{$inside_path}/{$page}.tpl";
		if (!file_exists($file) || !is_file($file)) {
			return false;
		}
			$this->pages[$page] = file_get_contents($file);
			$this->handling_includes($page);		
			return true;
	}

	private function handling_includes(string $page) {
		$page_arr = explode("\n", $this->pages[$page]);
		$c = count($page_arr);
		for ($i = 0; $i < $c; $i++) {
			if (preg_match('/{{{include "(.+)\/(.+)\.tpl"}}}\|m/', $page_arr[$i], $matches) != 0) {
				if ($this->get_page($matches[2], $matches[1]) === false) {
					exit("Cannot get the template file {$matches[1]}/{$matches[2]}.tpl");
				}
				$page_arr[$i] = str_replace("{{{include \"{$matches[1]}/{$matches[2]}.tpl\"}}}|m", $this->pages[$matches[2]], $page_arr[$i]);
				$this->pages[$page] = implode("\n",$page_arr);
				unset($this->pages[$matches[2]]);
			}
		}
		return true;
	}
	
	public function setVariable(string $page, string $marker, $change) {
		if (array_key_exists($page, $this->pages)) {
			$this->pages[$page] = str_replace("{{{$marker}}}|m", $change, $this->pages[$page]);
			return true;
		}
		return false;
	}
	
	private function clear_markers() {
		foreach ($this->pages as $key => $val):
			$this->pages[$key] = preg_replace('/{{.*}}\|m/', '', $val);
		endforeach;
		return true;
	}

	public function output() {
		$this->clear_markers();
		$output = "";
		foreach($this->pages as $val):
			$output .= $val;
		endforeach;
		return $output;
	}

	public function get_etag(string $file) {
		if (!is_file($file) || !file_exists($file)) {
			return false;
		}
		return filemtime($file);
	}
}

?>