go-crossfeed
============================

Playing around with fgms stuff and crossfeed for fun in golang

replay/
-------------------
Replays a Geoff cf_raw.log file, and transmits to UDP socket,
ie a simulated a fgms crossfeed, but from file instead of network.

An example of the log file truncated at 100k is at stuff/cf_test.log

To monitor the UDP download and compile crossfeed from
https://gitorious.org/fgtools/crossfeed (in c/c++)

```bash
cd replay/
# start with 2hz
go run main.go -l ../stuff/cf_test.log -z 2
```
