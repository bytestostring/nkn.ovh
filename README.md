# nkn.ovh
nkn.ovh - Open source monitoring for NKN nodes

## API

[Documentation](../v1.1/API.md)

________

## System requirements.

- Since version 1.1 the programm can work in standalone mode (without frontend proxy), but only for HTTP protocol.
- Network bandwidth 30 mbps+.
- MySQL 5.6+ / MariaDB 10+ with a InnoDB support.
- At least 512MB RAM
- (Optional) Any frontend server with WebSocket proxy support for HTTPS protocol access.
- For building the nknovh daemon, you need Golang 1.15 or higher


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

Note. If you have compiled WebAssembly (wasm.go), make sure to copy **wasm_exec.js** from your golang distribution to **web/static/js/** directory.
As example, for Go version 1.15, this file can be found here:
https://github.com/golang/go/blob/dev.boringcrypto.go1.15/misc/wasm/wasm_exec.js

2. Create a database and import the sql file this way:

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

4. Edit DB settings in the configuration file **conf.json**, Also if you use proxy server, add your proxy server IP into **TrustedProxies** json array.

6. Run daemon

```
./nknovh
```

7. Optionally you can use systemd script 

You can check journal files in the **logs** directory


## Upgrade from previous version (for building from sources)

0. Stop your nknovh: kill -9 / systemctl stop nknovh / etc 

1. Check the version of your nknovh

```
cat conf.json | grep "Version"
```

2. Upgrade using git, and recompile main daemon.

```
git pull
go build cmd/nknovh/nknovh.go
```

3. Optionally! ONLY If you want to recompile .wasm binary from source: 

```
GOOS=js GOARCH=wasm go build -ldflags=-s -o web/static/lib.wasm cmd/wasm/wasm.go
```

Note, if you have compiled lib.wasm, make sure to copy wasm_exec.js from your golang distribution to web/static/js directory.


4. Before restarting the nknovh, you must check for sql changes into "sqlupgrade" directory and import the changes if it exists.
Example, your previous version was "1.1.0" and you need upgrade your database structure to the latest version "1.1.0-dirty-4":


```
mysql -uroot -p nknovh < sqlupgrade/from-1.1.0-to-1.1.0-dirty-4.sql
```

5. Start the nknovh:

```
./nknovh
# Or if you are using the systemctl script
systemctl start nknovh
```

## Donation

If you want to support this project:  

- Ethereum: [0xD5305428401C9295401c89ff14CB8f6588A34F20](https://etherscan.io/address/0xD5305428401C9295401c89ff14CB8f6588A34F20)
- NKN: [NKNZKKF9u1MUQWnK272YoFiMTn5tjZh7uRQE](https://explorer.nkn.org/detail/address/NKNZKKF9u1MUQWnK272YoFiMTn5tjZh7uRQE/1)
