# DNAS - Domain Name Analytics System

Eventually this will actually do something besides logging DNS questions and Answers. 
The hope is to record the data and build metrics on usage and to search.

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

  ![dnas](https://github.com/mephux/dnas/raw/master/screenshot/dnas-output.png)


## JSON Output

  `Example: sudo dnas -i en0 -u mephuxi -w output.txt`

  ```json
  {"dns":{"answers":[{"class":"IN","name":"api.twitter.com.","record":"A","data":"199.16.156.8","ttl":"19"},{"class":"IN","name":"api.twitter.com.","record":"A","data":"199.16.156.199","ttl":"19"},{"class":"IN","name":"api.twitter.com.","record":"A","data":"199.16.156.231","ttl":"19"},{"class":"IN","name":"api.twitter.com.","record":"A","data":"199.16.156.72","ttl":"19"}],"question":"api.twitter.com."},"dstip":"172.16.1.19","protocol":"UDP","srcip":"172.16.1.1","timestamp":"2014-08-07T16:23:16.343281497-04:00"}
  ```
