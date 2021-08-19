<!DOCTYPE html>
<html lang="en"><head><title>NKN.OVH - A simple monitoring for your NKN nodes</title>	
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width"><meta name="theme-color" content="#148AE3" />
	<link href="/static/css/nknc.css?up={{style_etag}}|m" rel="stylesheet" property="stylesheet" title="" />
	<link rel="icon" href="/static/favs/favicon.ico" type="image/x-icon" />
	<link rel="shortcut icon" href="/static/favs/favicon.ico" type="image/x-icon">
	<script src="/static/js/wasm_exec.js?{{wexec_etag}}|m"></script>
	<script>
	if (typeof WebAssembly === "object") {
		const go = new Go();
	        fetch("/static/lib.wasm?{{wasm_etag}}|m").then(response =>
                        response.arrayBuffer()
                ).then(bytes =>
                        WebAssembly.instantiate(bytes, go.importObject)
                ).then(result => {
                        go.run(result.instance);
                });
	} else {
		alert("Your browser is not support WebAssembly");
	}
	</script>
</head>
<body>
<div id="baseloading"><div><img src="/static/images/baseloading.gif" alt="Loading"></div></div>
<div style="clear:  both;"></div>
<div class="wrapper" id="wrapper">
