# go-log [![GoDoc](https://godoc.org/gopkg.in/src-d/go-log.v0?status.svg)](https://godoc.org/github.com/src-d/go-log) [![Build Status](https://travis-ci.org/src-d/go-log.svg)](https://travis-ci.org/src-d/go-log) [![Build status](https://ci.appveyor.com/api/projects/status/15cdr1nk890qpk7g?svg=true)](https://ci.appveyor.com/project/mcuadros/go-log) [![codecov.io](https://codecov.io/github/src-d/go-log/coverage.svg)](https://codecov.io/github/src-d/go-log) [![Go Report Card](https://goreportcard.com/badge/github.com/src-d/go-log)](https://goreportcard.com/report/github.com/src-d/go-log)

Log is a generic logging library based on logrus (this may change in the
future), that minimize the exposure of src-d projects to logrus or any other
logging library, as well defines the standard for configuration and usage of the
logging libraries in the organization.

Installation
------------

The recommended way to install *go-log* is:

```
go get -u gopkg.in/src-d/go-log.v0/...
```

Usage
-----

### Logger instantiation

A basic example that instantiates a default `Logger`, by default the logger
will be configured a Info level and text format, if a TTY doesn't using the
default format will be JSON.

```go
logger, _ := log.New()
logger.Infof("The answer to life, the universe and everything is %d", 42)
// INFO The answer to life, the universe and everything is 42
```

Also a new `Logger` can be created from other `Logger` in order to have
contextual, information.

```go
logger, _ := log.New()

bookLogger := logger.New(log.Field{"book": "Hitchhiker's Guide To The Galaxy"})
bookLogger.Infof("The answer to life, the universe and everything is %d", 42)
// INFO The answer to life, the universe and everything is 42 book=Hitchhiker's Guide To The Galaxy
```

License
-------
Apache License Version 2.0, see [LICENSE](LICENSE)