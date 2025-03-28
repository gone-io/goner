# gone-redis

本库基于 [redigo](github.com/gomodule/redigo/redis) 集成了 Redis 的基本操作。

## 使用方法

### 0. Redis 服务器配置

本库使用 [gone-config](../config) 进行配置，因此可以在配置文件（config/default.properties，config/${env}.properties）中进行配置。

- redis.server：Redis 服务器地址，例如：`localhost:6379`。
- redis.password：Redis 服务器密码。
- redis.db：要使用的 Redis 数据库索引。
- redis.max-idle：Redis 连接池中的空闲连接数。
- redis.max-active：Redis 连接池中的最大活动连接数。
- redis.cache.prefix：用于隔离不同应用程序的前缀字符串。如果您的 Redis 被多个应用程序使用，建议使用此配置。例如，如果 `redis.cache.prefix=app-x`，那么 `Cache.Set("the-module-cache-key", value)` 将会在 Redis 中设置键值为 `app-x#the-module-cache-key`。

### 1. 使用 Redis 实现分布式缓存

```go
package demo

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
)

func NewService() gone.Goner {
	return &service{}
}

type service struct {
	gone.Flag
	cache redis.Cache `gone:"gone-redis-cache"` //标签标识
}

func (s *service) Use() {

	type ValueStruct struct {
		X int
		Y int
		//...
	}

	var v ValueStruct

	// 设置缓存到 Redis
	err := s.cache.Set("cache-key", v)
	if err != nil {
		//处理错误
	}

	//从缓存获取值
	err = s.cache.Get("cache-key", &v)
	if err != nil {
		//处理错误
	}

	//...
}
```

### 2. 使用 Redis 实现分布式锁

```go
package demo

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
	"time"
)

func NewService() gone.Goner {
	return &service{}
}

type service struct {
	gone.Flag
	locker redis.Locker `gone:"gone-redis-locker"`
}

//UseTryLock 使用 Locker.TryLock
func (s *service) UseTryLock() {
	unlock, err := s.locker.TryLock("a-lock-key-in-redis", 10*time.Second)
	if err != nil {
		// 处理错误
	}
	defer unlock()
	// ... 其他操作
}

//UseLockAndDo 使用 Locker.LockAndDo
func (s *service) UseLockAndDo() {
	err := s.locker.LockAndDo("a-lock-key-in-redis", func() {
		//执行业务逻辑
		//...
		//函数结束时会自动解锁
		//否则，锁会自动续期直到函数执行完成

	}, 10*time.Second, 2*time.Second)

	if err != nil {
		//处理错误
	}
}
```

### 3. Key 操作

```go
package demo

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
	"time"
)

func NewService() gone.Goner {
	return &service{}
}

type service struct {
	gone.Flag
	key redis.Key `gone:"gone-redis-locker"`
}

func (s *service) UseTryLock() {

	// 设置过期时间
	s.key.Expire("the-key-in-redis", 2*time.Second)
	s.key.ExpireAt("the-key-in-redis", time.Now().Add(10*time.Minute))

	// 获取键的 TTL
	s.key.Ttl("the-key-in-redis")

	// 等等
}
```

### 4. Redis 哈希

```go
package demo

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
)

func NewService() gone.Goner {
	return &service{}
}

type service struct {
	gone.Flag
	h redis.Hash `gone:"gone-redis-provider,key-in-redis"` //使用 gone-redis-provider 标签提供一个 redis.Hash 来操作 `key-in-redis` 上的哈希
}

func (s *service) Use() {
	s.h.Set("a-field", "some thing")
	var str string
	s.h.Get("a-field", &str)

	//...
}
```

### 5. Provider

> Provider 可以在应用程序中再次隔离键命名空间。例如，您想为模块 A 使用 `app-x#module-a` 作为 Redis 前缀，为模块 B 使用 `app-x#module-b` 作为前缀。您可以像下面这样使用：

```go
package A

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
)

//在模块 A 中
//...
type service struct {
	gone.Flag
	cache  redis.Cache  `gone:"gone-redis-provider,module-a"` //使用缓存
	key    redis.Key    `gone:"gone-redis-provider,module-a"` //使用键
	locker redis.Locker `gone:"gone-redis-provider,module-a"` //使用锁
}
```

```go
package B

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
)

//在模块 B 中
//...
type service struct {
	gone.Flag
	cache  redis.Cache  `gone:"gone-redis-provider,module-b"` //使用缓存
	key    redis.Key    `gone:"gone-redis-provider,module-b"` //使用键
	locker redis.Locker `gone:"gone-redis-provider,module-b"` //使用锁
}
```

如果键值在配置文件中，您可以使用 `gone:"gone-redis-provider,config=config-file-key,default=default-val"`。

```go
package A

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/redis"
)

//在模块 B 中
//...
type service struct {
	gone.Flag
	cache redis.Cache `gone:"gone-redis-provider,config=app.module-a.redis.prefix"` //使用缓存
}
```

### 5. Redis 连接池

您可以直接使用 `redis.Pool` 来读写 Redis，这是由 [redigo](github.com/gomodule/redigo/redis) 提供的。

```go
package demo

import (
	"github.com/gone-io/goner/redis"
	"github.com/gone-io/gone/v2"
)

type service struct {
	gone.Flag
	pool redis.Pool `gone:"gone-redis-pool"`
}

func (s *service) Use() {
	conn := s.pool.Get()
	defer s.pool.Close(conn)

	//执行一些操作
	//conn.Do(/*...*/)

	//发送命令
	//conn.Send(/*...*/)
}
```

## 测试

> 以下测试脚本依赖于 [Make](https://cmake.org/download/) 和 [Docker](https://www.docker.com/get-started/)，Docker 用于运行 Redis。

```shell
make test
```