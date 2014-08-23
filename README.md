# DNAS - Domain Name Analytics System
[![Build Status](https://drone.io/github.com/mephux/dnas/status.png)](https://drone.io/github.com/mephux/dnas/latest)

Eventually this will actually do something besides logging DNS questions and Answers and write to an embeded bolt (https://github.com/boltdb/bolt) key/value store. 
The hope is to record the data and build metrics on usage and for searching. i.e malware blah.exe sent data to blah.org what ips did that resolve to at that time.

## Install

  * Note: You will need libpcap-dev before you build DNAS.

  `go get github.com/mephux/dnas`
  
  -- OR --
  
  `git clone https://github.com/mephux/dnas.git`
  
  `cd dnas`
  
  `make`

## Usage

```
  DNAS (0.1.0) - Domain Name Analytics System

  Usage: dnas [options]

  Options:
    -i, --interface=eth0    Interface to monitor
    -p, --port=53           DNS port (53)
    -d, --database=         Database file path
    -f, --filter=*.com      Filter by question
    -D, --daemon            Run DNAS in daemon mode
    -w, --write=FILE        Write JSON output to log file
    -u, --user=USER         Drop privileges to this user
    -H, --hexdump           Show hexdump of DNS packet
    -v, --version           Show version information

  Help Options:
    -h, --help              Show this help message
```

## STDOUT

  `Example: sudo dnas -i en0 -u mephux`

  ![dnas](https://raw.githubusercontent.com/mephux/dnas/master/screenshot/dnas-output.png)


## JSON Output

  `Example: sudo dnas -i en0 -u mephux -w output.txt`

  ```json
  {"dns":{"answers":[{"class":"IN","name":"avatars2.githubusercontent.com.","record":"CNAME","data":"github.map.fastly.net.","ttl":"1099","created_at":"2014-08-17T17:10:38.194959151-04:00","updated_at":"2014-08-17T17:10:38.194959229-04:00","active":true},{"class":"IN","name":"github.map.fastly.net.","record":"A","data":"199.27.76.133","ttl":"4","created_at":"2014-08-17T17:10:38.194963092-04:00","updated_at":"2014-08-17T17:10:38.194963118-04:00","active":true}],"question":"avatars2.githubusercontent.com.","length":150},"dstip":"172.16.1.19","protocol":"UDP","srcip":"172.16.1.1","timestamp":"2014-08-17T17:10:38.19486575-04:00","packet":"i4WBgAABAAIAAAAACGF2YXRhcnMyEWdpdGh1YnVzZXJjb250ZW50A2NvbQAAAQABCGF2YXRhcnMyEWdpdGh1YnVzZXJjb250ZW50A2NvbQAABQABAAAESwAXBmdpdGh1YgNtYXAGZmFzdGx5A25ldAAGZ2l0aHViA21hcAZmYXN0bHkDbmV0AAABAAEAAAAEAATHG0yF"}
  ```

## Self-Promotion

Like DNAS? Follow the repository on
[GitHub](https://github.com/mephux/dnas) and if
you would like to stalk me, follow [mephux](http://dweb.io/) on
[Twitter](http://twitter.com/mephux) and
[GitHub](https://github.com/mephux).

# MIT LICENSE

The MIT License (MIT)

Copyright (c) 2014 Dustin Willis Webber

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
