# nkn.ovh
nkn.ovh - Open source monitoring for the NKN nodes

________

## System requirements.

- High network bandwidth 100 mbps+. We are strongly recommend to use a server with 250+ mbps bandwidth for the default options.
- As least 3 GB RAM
- PHP 8.0.x
- (Optional) Apache 2.2 or 2.4 with mod_rewrite and a php support module: mod_php, mod_lsapi or mod_fcgid. Also you can use nginx + php-fpm.
- For build the nknovh daemon you need Golang 1.15 or higher
- MySQL 5.6+ / MariaDB 10+ with a InnoDB support.


## Build from source

1. Get the package and build it:

```
git clone https://github.com/bytestostring/nkn.ovh.git
cd nkn.ovh
go build cmd/cmd/nknovh.go
```
2. Create a database and import the sql file like this:

```
mysql -uroot -p
CREATE DATABASE nknovh;
quit
mysql -uroot -p nknovh < struct.sql
```

3. Copy the two configuration files:

```
cp conf.json.example conf.json
cp web/engine/db_config.php.example web/engine/db_config.php

```

4. Edit DB settings in the configuration files **conf.json** and **web/engine/db_config.php**.

5. Move contain of **web** folder to your actually Apache/NGINX path. Like /srv/www/%domain%/
```
mv web/* /srv/www/%domain%/
```
6. Run daemon:

```
./nknovh
```

You can check journal files in the **logs** and **web/logs** folders.
