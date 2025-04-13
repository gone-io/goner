package xorm

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"xorm.io/xorm"
)

func contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func Test_xormProvider_Init(t *testing.T) {
	db, mock, _ := sqlmock.NewWithDSN(
		"root@/blog",
		sqlmock.MonitorPingsOption(true),
	)
	defer db.Close()

	drivers := sql.Drivers()
	if !contains(drivers, "mysql") {
		sql.Register("mysql", db.Driver())
	}

	type Tests struct {
		name        string
		before      func()
		after       func()
		exceptPanic bool
		injectFn    any
	}

	tests := []Tests{
		{
			name: "inject default db failed",
			before: func() {
				_ = os.Setenv("GONE_DATABASE", "{\n\t\"driver-name\":\"error\",\n\t\"dsn\": \"root@/blog\",\n\t\"max-idle-count\": 5,\n\t\"max-open\": 10,\n\t\"max-lifetime\": 10000,\n\t\"show-sql\": true,\n\t\"ping-after-init\": true\n}")
			},
			after: func() {
				_ = os.Unsetenv("GONE_DATABASE")
			},
			injectFn: func(in struct {
				db *xorm.Engine `gone:"*"`
			}) {
				assert.NotNil(t, db)
			},
			exceptPanic: true,
		},
		{
			name: "inject default db,ping error",
			before: func() {
				_ = os.Setenv("GONE_DATABASE", "{\n\t\"driver-name\":\"mysql\",\n\t\"dsn\": \"root@/blog\",\n\t\"max-idle-count\": 5,\n\t\"max-open\": 10,\n\t\"max-lifetime\": 10000,\n\t\"show-sql\": true,\n\t\"ping-after-init\": true\n}")
			},
			after: func() {
				_ = os.Unsetenv("GONE_DATABASE")
			},
			injectFn: func(in struct {
				db *xorm.Engine `gone:"*"`
			}) {
				assert.NotNil(t, db)
			},
			exceptPanic: true,
		},
		{
			name: "inject default db",
			before: func() {
				_ = os.Setenv("GONE_DATABASE", "{\n\t\"driver-name\":\"mysql\",\n\t\"dsn\": \"root@/blog\",\n\t\"max-idle-count\": 5,\n\t\"max-open\": 10,\n\t\"max-lifetime\": 10000,\n\t\"show-sql\": true,\n\t\"ping-after-init\": true\n}")
				mock.ExpectPing()
			},
			after: func() {
				_ = os.Unsetenv("GONE_DATABASE")
			},
			injectFn: func(in struct {
				db1      *xorm.Engine         `gone:""`
				db2      *xorm.Engine         `gone:"*"`
				db       xorm.EngineInterface `gone:"*"`
				injector gone.StructInjector  `gone:"*"`
			}) {
				assert.NotNil(t, in.db1)
				assert.NotNil(t, in.db2)
				assert.Equal(t, in.db1, in.db2)
				assert.Equal(t, in.db1, in.db)

				var x struct {
					db *xorm.EngineGroup `gone:"*"`
				}

				err := in.injector.InjectStruct(&x)
				assert.Error(t, err)
			},
			exceptPanic: false,
		},
		{
			name: "cluster master init failed",
			before: func() {
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_ENABLE", "true")
			},
			after: func() {
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_ENABLE")
			},
			injectFn: func(in struct {
				db *xorm.EngineGroup `gone:",db=customer"`
			}) {

			},
			exceptPanic: true,
		},
		{
			name: "cluster slave init failed",
			before: func() {
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_ENABLE", "true")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_MASTER", "{\n\t\"driver-name\":\"mysql\",\n\t\"dsn\": \"root@/blog\"}")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_SLAVES", "[{\n\t\"driver-name\":\"error\",\n\t\"dsn\": \"root@/blog\"}]")
			},
			after: func() {
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_ENABLE")
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_MASTER")
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_SLAVES")
			},
			injectFn: func(in struct {
				db *xorm.EngineGroup `gone:",db=customer"`
			}) {

			},
			exceptPanic: true,
		},
		{
			name: "cluster success",
			before: func() {
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_ENABLE", "true")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_MASTER", "{\n\t\"driver-name\":\"mysql\",\n\t\"dsn\": \"root@/blog\"}")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_SLAVES", "[{\n\t\"driver-name\":\"mysql\",\n\t\"dsn\": \"root@/blog\"}]")
			},
			after: func() {
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_ENABLE")
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_MASTER")
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_SLAVES")
			},
			injectFn: func(in struct {
				db *xorm.EngineGroup `gone:",db=customer"`
			}) {
				assert.NotNil(t, in.db)
			},
			exceptPanic: false,
		},
		{
			name: "use Engine receive EngineGroup",
			before: func() {
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_ENABLE", "true")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_MASTER", "{\n\t\"driver-name\":\"mysql\",\n\t\"dsn\": \"root@/blog\"}")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_SLAVES", "[{\n\t\"driver-name\":\"mysql\",\n\t\"dsn\": \"root@/blog\"}]")
			},
			after: func() {
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_ENABLE")
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_MASTER")
				_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_SLAVES")
			},
			injectFn: func(in struct {
				db *xorm.Engine `gone:",db=customer"`
			}) {
				assert.NotNil(t, in.db)
			},
			exceptPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			defer tt.after()

			func() {
				defer func() {
					err := recover()
					if tt.exceptPanic {
						assert.Error(t, err.(error))
					} else {
						assert.Nil(t, err)
					}
				}()

				gone.
					NewApp().
					Load(&xormProvider{}).
					Load(xormEngineProvider).
					Load(xormGroupProvider).
					Run(tt.injectFn)
			}()
		})
	}

}
