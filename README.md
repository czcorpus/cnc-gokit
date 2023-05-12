# CNC-GoKit

This is a set of functions used by the CNC across different Go-based
applications.

Please note that some convenience functions log warnings and errors
via [zerolog](https://github.com/rs/zerolog).

For a documentation please see https://pkg.go.dev/github.com/czcorpus/cnc-gokit

## Available packages

### collections

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

* `func ConnectAPI(conf *ConnectionConf, errListen <-chan error) *influxDBAdapter`
* `func RunWriteConsumerSync[T Influxable](db *influxDBAdapter, measurement string, incomingData <-chan T)`

### mail

The `mail` package contains functions for sending e-mails with simplicity in mind.
TLS and authentication is supported.

### strnum

The `strnum` package contains functions for converting between numbers, slices of
numbers etc. to strings (and in reverse).


### uniresp

The `uniresp` package contains functions usable for writing HTTP JSON responses.
