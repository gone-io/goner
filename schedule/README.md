<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/schedule component, and Configure Scheduled Tasks with Cron Expressions

Scheduled tasks are essential for web projects. You can implement scheduled tasks using Gone's built-in components, which support `cron` expressions. In web projects where code typically runs on multiple nodes, we use Redis as a distributed lock to ensure that each task executes on only one node at a time. Therefore, you need to set up Redis first. For more information about Redis, please refer to: [Using Redis for Distributed Locks and Caching](https://goner.fun/guide/redis.html). Additionally, scheduled tasks can be integrated with the framework's "configuration injection" feature to store cron expressions in configuration files. For more details, see [Support Configuration Files with Built-in Goners](https://goner.fun/guide/config.html).

## Register Related Goners to Gone
```go
    //Use the goner.SchedulePriest function to register schedule-related Goners to Gone
    _ = goner.SchedulePriest(cemetery)
```

## Write Job Functions for Scheduled Tasks
```go
type sch struct {
    gone.Flag
}

func (sch *sch) job1() {
    //todo scheduled task logic
}
```

## Configure Scheduled Tasks
Implement `Cron(run schedule.RunFuncOnceAt)`. The framework will scan this method on the struct and execute it automatically. Set up scheduled tasks within this method.
```go
func (sch *sch) Cron(run schedule.RunFuncOnceAt) {

    //Use run `RunFuncOnceAt` to set up scheduled tasks
    run(
        "*/5 * * * * *", // cron expression, executes every 5 seconds
        "job1",          //unique identifier for distributed lock
        sch.job1,        // scheduled task logic
    )
}
```

Here's the complete demo code, which can be found in [example](./example):
```go
package main

import (
    "fmt"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner"
    "github.com/gone-io/goner/redis"
    "github.com/gone-io/goner/schedule"
)

type sch struct {
    gone.Flag
}

func (sch *sch) job1() {
    //todo scheduled task logic
    fmt.Println("job1 execute")
}

func (sch *sch) Cron(run schedule.RunFuncOnceAt) {

    //Use run `RunFuncOnceAt` to set up scheduled tasks
    run(
        "*/5 * * * * *", // cron expression, executes every 5 seconds
        "job1",          //unique identifier for distributed lock
        sch.job1,        // scheduled task logic
    )
}

func main() {
    gone.
        Load(&sch{}).
        Loads(
            goner.BaseLoad,
            redis.Load, //Use Redis for distributed locks; not needed in single-node mode (when schedule.in-cluster=false)
            schedule.Load,
        ).
        Serve()
}
```

The above code will print `job1 execute` every 5 seconds. Simple, isn't it?

## Store Scheduling Configuration in Configuration Files
To store scheduling configuration in configuration files, you only need to make the following three changes to your code:

1. Register the configuration file-related Goners to Gone
2. Inject the scheduled task configuration from the configuration file
3. Use the injected configuration to set up scheduled tasks

Here's the modified `sch` code:
```go
type sch struct {
    gone.Flag

    cron string `gone:"config,cron.job1,default=*/5 * * * * *"` //2. Inject scheduled task configuration from config file
}

func (sch *sch) job1() {
    //todo scheduled task logic
    fmt.Println("job1 execute")
}

func (sch *sch) Cron(run schedule.RunFuncOnceAt) {

    //Use run `RunFuncOnceAt` to set up scheduled tasks
    run(
        sch.cron, // 3. Use configuration injected from config file
        "job1",   //unique identifier for distributed lock
        sch.job1, // scheduled task logic
    )
}
```