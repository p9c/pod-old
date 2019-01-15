# clog

Clog is a logger that works by receiving strings on channels and triggering an output through a configurable output function. The log level can be changed on the fly and immediately takes effect. There is colour available for tags if desired.

## Automatic

As soon as you import it into your application it is instantly activated by an automatic init routine.

## Channel-based

Because the loggers are channel based, the overhead of logging function calls, which can be quite extensive, is completely eliminated.

## Subsystems

One simply calls `clog.NewSubSystem` with a name and level and you get back a struct containing named channels that sit waiting to prepend your prescribed name string on whatever you put in the channel, and then forward it to the appropriate level of the root clog logger (accessible through `clog.L`)

```go
ss := clog.NewSubSystem("TEST", clog.Ndbg)
ss.Info <- "this will print an info message"
```

Output will be like this:

```
2019-01-15 11:59:58.155324 UTC [INF] this will print an info message
```