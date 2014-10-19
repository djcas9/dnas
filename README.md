# DNAS - Domain Name Analytics System
[![Build Status](https://drone.io/github.com/mephux/dnas/status.png)](https://drone.io/github.com/mephux/dnas/latest)

Logs all DNS questions and answers for searching and metrics. Supports logging to mysql, postgres, sqlite3 and json.

## Install

  * Using Go Get

    * Note: You will need libpcap-dev before you build DNAS.
    * `go get github.com/mephux/dnas`

  * Using Git & Go Build
  
    * Note: You will need libpcap-dev before you build DNAS.
    * `git clone https://github.com/mephux/dnas.git`
    * `cd dnas`
    * `make`

  * Using Vagrant & Docker

    * `vagrant up`

## Usage

  `Example: sudo dnas -i en0 -H -u mephux`

```
  DNAS (0.2.0) - Domain Name Analytics System

  Usage: dnas [options]

  Options:
    -i, --interface=eth0          Interface to monitor
    -p, --port=53                 DNS port (53)
    -D, --daemon                  Run DNAS in daemon mode
    -w, --write=FILE              Write JSON output to log file
    -u, --user=USER               Drop privileges to this user
    -H, --hexdump                 Show hexdump of DNS packet
        --mysql                   Enable Mysql Output Support
        --postgres                Enable Postgres Output Support
        --sqlite3                 Enable Sqlite3 Output Support
        --db-user=root            Database User (root)
        --db-password=PASSWORD    Database Password
        --db-database=dnas        Database Database (dnas)
        --db-host=127.0.0.1       Database Host
        --db-port=3306            Database Port
        --db-path=~/.dnas.db      Path to Database on disk. (sqlite3 only)
        --db-tls                  Enable TLS / SSL encrypted connection to the database. (mysql/postgres only)
        --db-skip-verify          Allow Self-signed or invalid certificate (mysql/postgres only)
        --db-verbose              Show database logs in STDOUT
    -q, --quiet                   Suppress DNAS output
    -v, --version                 Show version information

  Help Options:
    -h, --help                    Show this help message
```

## OUTPUT Support

    * sqlite3
    * Mysql
    * Postgres
    * Json


## Self-Promotion

Like DNAS? Follow the repository on
[GitHub](https://github.com/mephux/dnas) and if
you would like to stalk me, follow [mephux](http://dweb.io/) on
[Twitter](http://twitter.com/mephux) and
[GitHub](https://github.com/mephux).

# MIT LICENSE

The MIT License (MIT) - [LICENSE](https://github.com/mephux/dnas/blob/master/LICENSE)
