# log
log library with multiple output and level
```
go get github.com/moorada/log
```

## Example
```
config := log.FormatConfigBasic
config.Format = "{level:name} {message}"
_ = log.AddOutput("logs_test1.log", log.DEBUG, config, true)

config2 := log.FormatConfigBasic
config2.Format = "{time} {level:color}{level:name}{reset} {message}"
_ = log.AddOutput("", log.ERROR, config2, true) // Stdout

log.Debug("Hello world! %s", "Hello")
log.Error("Hello world! %s", "ciao1")

_ = log.RemoveOutput("logs_test1.log")
```

This library is an extension of [isLazy/log by evilsocket](https://github.com/evilsocket/islazy/tree/master/log).
