# CNC-GoKit

This is a set of functions used by the CNC across different Go-based
applications.

Please note that some convenience functions log warnings and errors
via [zerolog](https://github.com/rs/zerolog).

For a documentation please see https://pkg.go.dev/github.com/czcorpus/cnc-gokit

## Available packages

### collections

- `CircularList`
- `ConcurrentMap`
- `MultiDict`
- `Set`

### datetime

- mostly ISO 8601 related functions

### fs

The `fs` package contains miscellaneous functions for dealing with
filesystems (obtaining file information, testing existence,...).

### influx

The `influx` serves as a wrapper for the InfluxDB client v2 offering a convenient
way of storing data in an InfluxDB database.

* `func ConnectAPI(conf *ConnectionConf, errListen <-chan error) *InfluxDBAdapter`
* `func RunWriteConsumerSync[T Influxable](db *InfluxDBAdapter, measurement string, incomingData <-chan T)`

### logging


The `logging` package contains functions for a service logging setup based
on ZeroLog.

* `type LogLevel string`
* `func SetupLogging(path string, level LogLevel)`
* `func GinMiddleware() gin.HandlerFunc`


### mail

The `mail` package contains functions for sending e-mails with simplicity in mind.
TLS and authentication is supported.

### maths

The `maths` package contains few useful functions for working with
numbers (`Max`, `Min`, `RoundToN`) and statistics (`OnlineMean`)

### strnum

The `strnum` package contains functions for converting between numbers, slices of
numbers etc. to strings (and in reverse).

### unireq

The `unireq` package contains helper functions for working with an HTTP request.

* `func CheckSuperfluousURLArgs(req *http.Request, allowedArgs []string)`
* `func ClientIP(req *http.Request) net.IP`

### uniresp

The `uniresp` package contains functions usable for writing HTTP JSON responses.

### util

* `Max`
* `Min`
