# DNAS - Domain Name Analytics System
[![Build Status](https://drone.io/github.com/mephux/dnas/status.png)](https://drone.io/github.com/mephux/dnas/latest)

Logs all DNS questions and answers for searching and metrics. DNA answers are stored in
a bloom filter for better performance (less questions asked to the embeded key/value store).

DNAS supports logging to an embeded Bolt (https://github.com/boltdb/bolt) key / value store 
(-d database.db) or to a flat file as json (-w filename.txt).

The plan is to continue to scale DNAS for powerful searches and metrics. 
i.e malware blah.exe sent data to blah.org what ips did that resolve to at that time.

## Install

  1. Go Get

    * Note: You will need libpcap-dev before you build DNAS.
    * `go get github.com/mephux/dnas`

  2. Git
  
    * Note: You will need libpcap-dev before you build DNAS.
    * `git clone https://github.com/mephux/dnas.git`
    * `cd dnas`
    * `make`

  3. Vagrant & Docker

    * `vagrant up`

## OUTPUT

  `Example: sudo dnas -i en0 -H -u mephux`

  ![dnas](https://raw.githubusercontent.com/mephux/dnas/master/dnas.gif)


## Usage

```
  DNAS (0.1.0) - Domain Name Analytics System

    Usage: dnas [options]

    Options:
      -i, --interface=eth0          Interface to monitor
      -p, --port=53                 DNS port (53)
      -d, --database=FILE           Database file path (dnas.db)
      -F, --filter=*.com            Filter by question
      -D, --daemon                  Run DNAS in daemon mode
      -w, --write=FILE              Write JSON output to log file
      -u, --user=USER               Drop privileges to this user
      -H, --hexdump                 Show hexdump of DNS packet
      -q, --find-question=STRING    Search for DNS record by question
      -a, --find-answer=STRING      Search for DNS records by answer data
      -l, --list                    List all seen DNS questions
      -v, --version                 Show version information

    Help Options:
      -h, --help                    Show this help message
```

## JSON Output

  `Example: sudo dnas -i en0 -u mephux -w output.txt`

  ```json
  {
    "dns": {
      "answers": [
        {
          "class": "IN",
          "name": "github.com.",
          "record": "A",
          "data": "192.30.252.130",
          "ttl": "24",
          "created_at": "2014-08-24T16:02:20.56537176-04:00",
          "updated_at": "2014-08-24T16:02:20.565371798-04:00",
          "active": true
        }
      ],
      "question": "github.com.",
      "length": 54
    },
    "dstip": "172.16.1.19",
    "protocol": "UDP",
    "srcip": "172.16.1.1",
    "timestamp": "2014-08-24T16:02:20.565137019-04:00",
    "packet": "xBaBgAABAAEAAAAABmdpdGh1YgNjb20AAAEAAQZnaXRodWIDY29tAAABAAEAAAAYAATAHvyC",
    "bloom": "eyJGaWx0ZXJTZXQiOiJBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQ0FBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFnQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBSUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUNBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBZ0FBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUlBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFDQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBPT0iLCJTZXRMb2NzIjo3fQ=="
  }
  ```

## Self-Promotion

Like DNAS? Follow the repository on
[GitHub](https://github.com/mephux/dnas) and if
you would like to stalk me, follow [mephux](http://dweb.io/) on
[Twitter](http://twitter.com/mephux) and
[GitHub](https://github.com/mephux).

# MIT LICENSE

The MIT License (MIT) - [LICENSE](https://github.com/mephux/dnas/blob/master/LICENSE)
