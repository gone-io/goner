package main

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
	"github.com/gone-io/goner/schedule"
	"github.com/gone-io/goner/viper"
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
			viper.Load,
			redis.Load, //使用 redis 实现分布式锁；单机模式下（即配置为schedule.in-cluster=false）时，不需要加载redis，
			schedule.Load,
		).
		Serve()
}
