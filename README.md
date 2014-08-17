# DNAS - Domain Name Analytics System
[![Build Status](https://drone.io/github.com/mephux/dnas/status.png)](https://drone.io/github.com/mephux/dnas/latest)

Eventually this will actually do something besides logging DNS questions and Answers. 
The hope is to record the data and build metrics on usage and for searching. i.e malware blah.exe sent data to blah.org what ips did that resolve to at that time.

## Install

  `go get github.com/mephux/dnas`

  then just run `dnas` if `$GOPATH/bin` is in your path.

## Usage

```
  DNAS (0.1.0) - Domain Name Analytics System

  Usage: dnas [options]

  Options:
    -i, --interface=eth0    Interface to monitor
    -p, --port=53           DNS port (53)
    -f, --filter=*.com      Filter by question
    -d, --daemon            Run DNAS in daemon mode
    -w, --write=FILE        Write JSON output to log file
    -u, --user=USER         Drop privileges to this user
    -v, --version           Show version information

  Help Options:
    -h, --help              Show this help message
```

## STDOUT

  `Example: sudo dnas -i en0 -u mephux`

  ![dnas](https://raw.githubusercontent.com/mephux/dnas/master/screenshot/dnas-output.png)


## JSON Output

  `Example: sudo dnas -i en0 -u mephuxi -w output.txt`

  ```json
  {"dns":{"answers":[{"class":"IN","name":"api.twitter.com.","record":"A","data":"199.16.156.8","ttl":"19"},{"class":"IN","name":"api.twitter.com.","record":"A","data":"199.16.156.199","ttl":"19"},{"class":"IN","name":"api.twitter.com.","record":"A","data":"199.16.156.231","ttl":"19"},{"class":"IN","name":"api.twitter.com.","record":"A","data":"199.16.156.72","ttl":"19"}],"question":"api.twitter.com."},"dstip":"172.16.1.19","protocol":"UDP","srcip":"172.16.1.1","timestamp":"2014-08-07T16:23:16.343281497-04:00"}
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
