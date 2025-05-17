
# `goner/xorm`使用说明

## import 和 加载

- import `goner` 包
```go
"github.com/gone-io/goner/xorm"
```
- 按需导入数据库驱动包
  - mysql 驱动
  ```go
  _ "github.com/go-sql-driver/mysql"
  ```
  - sqlite3 驱动
  ```go
  _ "github.com/mattn/go-sqlite3"
  ```
  - postgres 驱动
  ```go
  _ "github.com/lib/pq"
  ```
  - oracle 驱动
  ```go
  _ "github.com/mattn/go-oci8"
  ```
  - mssql 驱动
  ```go
  _ "github.com/denisenkom/go-mssqldb"
  ```


## 配置项

| 配置项                                     |         必须         | 默认值 | 描述                                                            |
| :----------------------------------------- | :------------------: | ------ | --------------------------------------------------------------- |
| **database.cluster.enable**                |          否          | false  | 是否启用集群模式                                                |
| **database.driver-name**                   | 否（非集群模式必填） | -      | 数据库驱动名称，支持mysql、sqlite3、postgres、oracle、mssql     |
| **database.dsn**                           | 否（非集群模式必填） | -      | 数据库连接字符串                                                |
| **database.max-idle-count**                |          否          | 5      | 连接池最大空闲连接数                                            |
| **database.max-open**                      |          否          | 20     | 连接池最大连接数                                                |
| **database.max-lifetime**                  |          否          | 10m    | 连接池连接最大存活时间                                          |
| **database.show-sql**                      |          否          | true   | 是否打印SQL日志                                                 |
| **database.cluster.master.driver-name**    |  否（集群模式必填）  | -      | 主库数据库驱动名称，支持mysql、sqlite3、postgres、oracle、mssql |
| **database.cluster.master.dsn**            |  否（集群模式必填）  | -      | 主库数据库连接字符串                                            |
| **database.cluster.slaves[n].driver-name** |  否（集群模式必填）  | -      | 从库数据库驱动名称，支持mysql、sqlite3、postgres、oracle、mssql |
| **database.cluster.slaves[n].dsn**         |  否（集群模式必填）  | -      | 从库数据库连接字符串                                            |

### 非集群模式例子
```ini
database.driver-name=mysql
database.dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```
### 集群模式例子
```ini
database.cluster.enable=true

# 主数据库配置
database.cluster.master.driver-name=mysql
database.cluster.master.dsn=root:123456@tcp(master-db-host:3306)/test?charset=utf8mb4&parseTime=True&loc=Local

# 从数据库配置
database.cluster.slaves[0].driver-name=mysql
database.cluster.slaves[0].dsn=root:123456@tcp(slave-db-0-host:3306)/test?charset=utf8mb4&parseTime=True&loc=Local

database.cluster.slaves[1].driver-name=mysql
database.cluster.slaves[1].dsn=root:123456@tcp(slave-db-1-host:3306)/test?charset=utf8mb4&parseTime=True&loc=Local

database.cluster.slaves[2].driver-name=mysql
database.cluster.slaves[2].dsn=root:123456@tcp(slave-db-1-host:3306)/test?charset=utf8mb4&parseTime=True&loc=Local

# ... 更多从数据库
```
## 在代码中注入数据库引擎
```go
import "github.com/gone-io/gone/v2"

type struct dbUser struct {
	gone.Flag
	db gone.XormEngine `gone:"*"` //注入数据库引擎【集群模式，注入的是数据库集群，通过该方式获取的引擎查询数据时会随机到各从数据库，写入数据时会使用主数据库】
	masterDb gone.XormEngine `gone:"xorm,master"` // 注入主数据库，集群模式下有效
	slaveDb0 gone.XormEngine `gone:"xorm,slave=0"` //注入从数据库0，集群模式下有效
	slaveDb1 gone.XormEngine `gone:"xorm,slave=1"` //注入从数据库1，集群模式下有效
	slaveDbs []gone.XormEngine `gone:"xorm"`       //主入从数据库Slice，集群模式下有效
}

type Book struct {
	Id int64
	Title string
}

func (d *dbUser) GetBookById(id int64) (book *Book, err error) {
	book = new(Book)
	has, err := d.db.Where("id=?", id).Get(book)
	if err != nil {
		return nil, gone.ToError(err)
	}
	if !has {
		return nil, gone.NewParameterError("book not found", 404)
	}
	return book, nil
}
```
## 对多数据库的支持
### 1. 配置多个数据库
多数据库配置前缀为`database.{数据库名称}`，例如：`database.db1`、`database.db2`。

| 配置项                                             |         必须         | 默认值 | 描述                                                            |
| :------------------------------------------------- | :------------------: | ------ | --------------------------------------------------------------- |
| **{数据库配置前缀}.cluster.enable**                |          否          | false  | 是否启用集群模式                                                |
| **{数据库配置前缀}.driver-name**                   | 否（非集群模式必填） | -      | 数据库驱动名称，支持mysql、sqlite3、postgres、oracle、mssql     |
| **{数据库配置前缀}.dsn**                           | 否（非集群模式必填） | -      | 数据库连接字符串                                                |
| **{数据库配置前缀}.max-idle-count**                |          否          | 5      | 连接池最大空闲连接数                                            |
| **{数据库配置前缀}.max-open**                      |          否          | 20     | 连接池最大连接数                                                |
| **{数据库配置前缀}.max-lifetime**                  |          否          | 10m    | 连接池连接最大存活时间                                          |
| **{数据库配置前缀}.show-sql**                      |          否          | true   | 是否打印SQL日志                                                 |
| **{数据库配置前缀}.cluster.master.driver-name**    |  否（集群模式必填）  | -      | 主库数据库驱动名称，支持mysql、sqlite3、postgres、oracle、mssql |
| **{数据库配置前缀}.cluster.master.dsn**            |  否（集群模式必填）  | -      | 主库数据库连接字符串                                            |
| **{数据库配置前缀}.cluster.slaves[n].driver-name** |  否（集群模式必填）  | -      | 从库数据库驱动名称，支持mysql、sqlite3、postgres、oracle、mssql |
| **{数据库配置前缀}.cluster.slaves[n].dsn**         |  否（集群模式必填）  | -      | 从库数据库连接字符串                                            |

### 2. 多数据库注入
```go
import "github.com/gone-io/gone/v2"

type struct dbUser struct {
	gone.Flag

	// 注入`{数据库配置前缀}`为`database.db1`的数据库引擎【集群模式，注入的是数据库集群，通过该方式获取的引擎查询数据时会随机到各从数据库，写入数据时会使用主数据库】
	db1 gone.XormEngine `gone:"xorm,db=database.db1"`

	// 注入`{数据库配置前缀}`为`database.db1`的主数据库，集群模式下有效
	masterDb gone.XormEngine `gone:"xorm,db=database.db1,master"`

	// 注入`{数据库配置前缀}`为`database.db1`的从数据库0，集群模式下有效
	slaveDb0 gone.XormEngine `gone:"xorm,db=database.db1,slave=0"`

	// 注入`{数据库配置前缀}`为`database.db1`的从数据库1，集群模式下有效
	slaveDb1 gone.XormEngine `gone:"xorm,db=database.db1,slave=1"`

	// 主入`{数据库配置前缀}`为`database.db1`的从数据库Slice，集群模式下有效
	slaveDbs []gone.XormEngine `gone:"xorm,db=database.db1"`
}
```

## Gone对Xorm的增强

### 1. 自动事务
使用`Transaction`函数包裹的函数，执行前会自动开启事务，返回`error`或者发生`panic`自动完成事务回滚，不返回`error`则自动提交事务。

> 注意：在`Transaction`函数包裹的数据库操作函数需要使用`session xorm.Interface`执行数据库操作

```go
type db struct {
	gone.Flag
	gone.XormEngine `gone:"gone-xorm"` //注入数据库引擎
}

func (d *db) updateUser(user *entity.User) error {

    // 使用Transaction包裹的函数，执行前会自动开启事务
	return d.Transaction(func(session xorm.Interface) error {

        //注意：使用的session进行数据库操作
		_, err := session.ID(user.Id).Update(user)
		return gone.ToError(err)
	})
}
```

### 2. 事务自动传递
嵌套使用`Transaction`函数包裹的函数，只会开启一个事务，嵌套事务会自动传递，嵌套事务的`error`或者发生`panic`自动完成事务回滚，不返回`error`则自动提交事务。

这样带来一个好处，让我们编写的函数在组合时能够自动合并到一个事务中。

请看下面代码，如果`updateUser`、`updateFriends`函数单独使用，会分别开启事务；将他们嵌套在`DoUpdate`的`Transaction`的函数中，则会合并到一个事务中。

```go
type db struct {
	gone.Flag
	gone.XormEngine `gone:"gone-xorm"` //注入数据库引擎
}

func (d *db) updateUser(user *entity.User) error {

    // 使用Transaction包裹的函数，执行前会自动开启事务
	return d.Transaction(func(session xorm.Interface) error {

        //注意：使用的session进行数据库操作
		_, err := session.ID(user.Id).Update(user)
		return gone.ToError(err)
	})
}

func (d *db) updateFriends(userId int64, friedns []*entity.Friend) error {
	return d.Transaction(func(session xorm.Interface) error {
		//todo: 更新朋友的相关操作

		return nil
	})
}

func (d *db) DoUpdate(user *entity.User, friedns []*entity.Friend) error {
	return d.Transaction(func(session xorm.Interface) error {
		err := d.updateUser(user)
		if err != nil {
			return err
		}

		return d.updateFriends(user.Id, friedns)
	})
}
```

### 3. SQL支持名字参数

```go
	sql, args := xorm.MustNamed(`
		update user
		set
		    status = :status,
			avatar = :avatar
		where
		    id = :id`,
		map[string]any{
			"id":     1,
			"status": 1,
			"avatar": "https://wwww....",
		},
	)
```
通过`xorm.MustNamed`处理后的sql为：
```sql
update user
set
    status = ?,
    avatar = ?
where
    id = ?
```

`args`为`[]any`类型，值为：`1,1,"https://wwww...."`。

### 4. SQL查询增强

使用`Sqlx`方法可以直接执行原生SQL语句，支持命名参数和IN查询：

```go
// 使用命名参数
session := engine.Sqlx("select * from user where id = :id and status = :status",
    map[string]any{
        "id":     1,
        "status": 1,
    })

// 支持IN查询
session := engine.Sqlx("select * from user where id in (?)", []int64{1, 2, 3})
```

### 5. 性能优化建议

1. 合理配置连接池参数
   - `max-idle-count`: 根据业务负载调整空闲连接数
   - `max-open`: 避免设置过大，防止数据库连接过载
   - `max-lifetime`: 设置合适的连接存活时间，避免连接长期占用

2. 使用集群模式时的优化
   - 读写分离：查询操作自动分发到从库
   - 合理分配主从库的连接池配置
   - 使用事务时注意主从延迟的影响

3. SQL优化建议
   - 使用`show-sql`参数在开发环境查看SQL执行情况
   - 合理使用索引
   - 避免大事务，合理控制事务粒度
   - 使用批量操作替代循环单条操作
