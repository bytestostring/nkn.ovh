# nkn.ovh
nkn.ovh - Open source monitoring for the NKN nodes

________

## System requirements.

- Starting at version 1.1 the programm can working in standalone mode (without frontend proxy), but only for HTTP protocol.
- Network bandwidth 30 mbps+.
- MySQL 5.6+ / MariaDB 10+ with a InnoDB support.
- As least 512MB RAM
- (Optional) Any frontend server with WebSocket proxy support for HTTPS protocol access.
- For build the nknovh daemon you need Golang 1.15 or higher


## Build from source

1. Get the package and build it:

```
git clone https://github.com/bytestostring/nkn.ovh.git
cd nkn.ovh
# Compile main daemon
go build cmd/nknovh/nknovh.go
# Compile WebAssembly (optionally)
GOOS=js GOARCH=wasm go build -ldflags=-s -o web/static/lib.wasm cmd/wasm/wasm.go
```

Note. If you have compiled WebAssembly (wasm.go), then you must copy **wasm_exec.js** from your golang distribution to **web/static/js/** directory.
As example, for Go version 1.15 that file can be found here:
https://github.com/golang/go/blob/dev.boringcrypto.go1.15/misc/wasm/wasm_exec.js

2. Create a database and import the sql file like this:

```
mysql -uroot -p
CREATE DATABASE nknovh;
quit
mysql -uroot -p nknovh < struct.sql
```

3. Copy the configuration file:

```
cp conf.json.example conf.json

```

4. Edit DB settings in the configuration file **conf.json**, Also if you use proxy server, you must add your proxy server IP into **TrustedProxies** json array.

6. Run daemon

```
./nknovh
```

7. Optionally you can use systemd script 

You can check journal files in the **logs** directory
