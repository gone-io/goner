
# 用cron表达式配置定时任务
定时任务对于Web项目基本上时标配，可以通过Gone的内置组件来实现定时任务，支持`cron`表达式。在Web项目中代码一般都是多节点运行，我们使用了redis作为分布式锁来保证任务每次执行只在一个节点上进行，所以需要先准备redis服务，关于redis相关内容请参考：[利用redis提供分布式锁和分布式缓存](https://goner.fun/zh/guide/redis.html)。另外定时任务还可以和框架“配置注入”的特性结合，将cron表达式放到配置文件中，参考[通过内置Goners支持配置文件](https://goner.fun/zh/guide/config.html)。

## 将相关Goners注册到Gone
```go
	//使用 goner.SchedulePriest 函数，将 定时任务 相关的Goner 注册到Gone
	_ = goner.SchedulePriest(cemetery)
```

## 编写定时任务执行的Job函数
```go
type sch struct {
	gone.Flag
}

func (sch *sch) job1() {
	//todo 定时任务逻辑
}
```


## 设置定时任务
实现`Cron(run schedule.RunFuncOnceAt) `，框架会扫描结构体上的该方法并自动执行，在该方法中设置定时任务。
```go
func (sch *sch) Cron(run schedule.RunFuncOnceAt) {

	//使用 run `RunFuncOnceAt`设置定时任务，
	run(
		"*/5 * * * * *", // cron 表达式，表示每5秒执行一次
		"job1",          //需要设置一个唯一标识，用于 分布式锁加锁
		sch.job1,        // 定时任务逻辑
	)
}
```

完整的demo代码如下，代码可以在[example](./example)中查看：
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
	//todo 定时任务逻辑
	fmt.Println("job1 execute")
}

func (sch *sch) Cron(run schedule.RunFuncOnceAt) {

	//使用 run `RunFuncOnceAt`设置定时任务，
	run(
		"*/5 * * * * *", // cron 表达式，表示每5秒执行一次
		"job1",          //需要设置一个唯一标识，用于 分布式锁加锁
		sch.job1,        // 定时任务逻辑
	)
}

func main() {
	gone.
		Load(&sch{}).
		Loads(
			goner.BaseLoad,
			redis.Load, //使用 redis 实现分布式锁；单机模式下（即配置为schedule.in-cluster=false）时，不需要加载redis，
			schedule.Load,
		).
		Serve()
}

```

上面代码会每隔5s打印：`job1 execute`，是不是很简单？

## 将定时配置放到配置文件中
将定时配置放到配置文件中，代码上只需要做如下3点修改：

1. 将配置文件支持的相关Goner 注册到Gone
2. 注入放到配置文件的定时任务配置
3. 使用从配置文件注入的定时配置设置定时任务

修改后的`sch`代码如下：
```go


type sch struct {
	gone.Flag

	cron string `gone:"config,cron.job1,default=*/5 * * * * *"` //2. 注入放到配置文件的定时任务配置
}

func (sch *sch) job1() {
	//todo 定时任务逻辑
	fmt.Println("job1 execute")
}

func (sch *sch) Cron(run schedule.RunFuncOnceAt) {

	//使用 run `RunFuncOnceAt`设置定时任务，
	run(
		sch.cron, // 3. 使用从配置文件注入的定时配置
		"job1",   //需要设置一个唯一标识，用于 分布式锁加锁
		sch.job1, // 定时任务逻辑
	)
}
```