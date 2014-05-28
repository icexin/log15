![obligatory xkcd](http://imgs.xkcd.com/comics/standards.png)

# log15

Package log15 provides an opinionated, simple toolkit for best-practice logging that is both human and machine readable. It is modeled after the standard library's io and net/http packages.

## Features
- A simple, easy-to-understand API
- Promotes structured logging by encouraging use of key/value pairs
- Child loggers which inherit and add their own private context
- Lazy evaluation of expensive operations
- Simple Handler interface allowing for construction of flexible, custom logging configurations with a tiny API.
- Color terminal support
- Built-in support for logging to files, streams, syslog, and the network
- Support for forking records to multiple handlers, buffering records for output, failing over from failed handler writes, + more

## Documentation

The package documentation is extensive and complete. Browse on godoc:

#### [log15 API Documentation](https://godoc.org/github.com/inconshreveable/log15)

## Versioning
The API of the master branch of log15 should always be considered unstable. Using a stable version
of the log15 package is supported by gopkg.in. Include your dependency like so:

    import log "gopkg.in/inconshreveable/log15.v1"

You can also vendor log15 with a tool like Godep.

## Examples

    // all loggers can have key/value context
    srvlog := log.New("module", "app/server")

    // all log messages can have key/value context 
    srvlog.Warn("abnormal conn rate", "rate", curRate, "low", lowRate, "high", highRate)

    // child loggers with inherited context
    connlog := srvlog.New("raddr", c.RemoteAddr())
    connlog.Info("connection open")

    // lazy evaluation
    connlog.Debug("ping remote", "latency", log.Lazy(pingRemote))

    // flexible configuration
    srvlog.SetHandler(log.MultiHandler(
        log.StreamHandler(os.Stderr, log.LogfmtFormat()),
        log.LvlFilterHandler(
            log.LvlError,
            log.Must.FileHandler("errors.json", log.JsonHandler())))

## FAQ

### The varargs style is brittle and error prone! Can I have type saftey please?
Yes. Use log.Ctx:

    srvlog := log.New(log.Ctx{"module": "app/server"})
    srvlog.Warn("abnormal conn rate", log.Ctx{"rate": curRate, "low": lowRate, "high": highRate})

## License
Apache
