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

func Test_engProvider_Provide(t *testing.T) {
	db, mock, _ := sqlmock.NewWithDSN(
		"test-db",
		sqlmock.MonitorPingsOption(true),
	)
	defer db.Close()

	drivers := sql.Drivers()
	if !contains(drivers, "sqlite") {
		sql.Register("sqlite", db.Driver())
	}

	tests := []struct {
		name      string
		before    func() (after func())
		runFn     any
		panicWant bool
	}{
		{
			name: "inject cluster slaves slice success && inject cluster master success and inject cluster success",
			before: func() (after func()) {
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_ENABLE", "true")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_MASTER", "{\n\t\"driver-name\":\"sqlite\",\n\t\"dsn\": \"test-db\"}")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_SLAVES", "[{\n\t\"driver-name\":\"sqlite\",\n\t\"dsn\": \"test-db\"},{\n\t\"driver-name\":\"sqlite\",\n\t\"dsn\": \"test-db\"}]")

				mock.ExpectPing()
				mock.ExpectBegin()
				mock.ExpectCommit()

				return func() {
					_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_ENABLE")
					_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_MASTER")
					_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_SLAVES")
				}
			},
			runFn: func(in struct {
				eng0 Engine   `gone:",db=customer"`
				eng1 Engine   `gone:"*,db=customer,master"`
				eng2 []Engine `gone:"xorm,db=customer,slaves"`
				eng3 Engine   `gone:"xorm,db=customer,slave=0"`
				eng4 Engine   `gone:"xorm,db=customer,slave=1"`
			}) {
				in.eng0.SetPolicy(xorm.RandomPolicy())
				err := in.eng1.Ping()
				assert.Nil(t, err)
				err = in.eng1.Transaction(func(session xorm.Interface) error {
					return nil
				})
				assert.Nil(t, err)

				group, ok := in.eng0.GetOriginEngine().(*xorm.EngineGroup)

				assert.Truef(t, ok, "inject cluster instance failed")
				slaves := group.Slaves()
				assert.Equal(t, group.Master(), in.eng1.GetOriginEngine())

				assert.Equal(t, 2, len(slaves))
				assert.Equal(t, slaves[0], in.eng3.GetOriginEngine())
				assert.Equal(t, slaves[0], in.eng2[0].GetOriginEngine())
				assert.Equal(t, slaves[1], in.eng4.GetOriginEngine())
				assert.Equal(t, slaves[1], in.eng2[1].GetOriginEngine())
			},
			panicWant: false,
		},
		{
			name: "group inject error",
			before: func() (after func()) {
				return func() {}
			},
			runFn: func(in struct {
				eng0 []Engine `gone:"xorm,db=customer,slaves"`
			}) {
			},
			panicWant: true,
		},
		{
			name: "inject ont support type",
			before: func() (after func()) {
				return func() {}
			},
			runFn: func(in struct {
				eng0 []xorm.GroupPolicy `gone:"xorm,db=customer,slaves"`
			}) {
			},
			panicWant: true,
		},
		{
			name: "inject cluster error",
			before: func() (after func()) {
				return func() {}
			},
			runFn: func(in struct {
				eng0 Engine `gone:"xorm,db=customer"`
			}) {
			},
			panicWant: true,
		},
		{
			name: "inject cluster master error",
			before: func() (after func()) {
				return func() {}
			},
			runFn: func(in struct {
				eng0 Engine `gone:"xorm,db=customer,master"`
			}) {
			},
			panicWant: true,
		},
		{
			name: "inject cluster slave error",
			before: func() (after func()) {
				return func() {}
			},
			runFn: func(in struct {
				eng0 Engine `gone:"xorm,db=customer,slave"`
			}) {
			},
			panicWant: true,
		},
		{
			name: "inject cluster slave error",
			before: func() (after func()) {
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_ENABLE", "true")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_MASTER", "{\n\t\"driver-name\":\"sqlite\",\n\t\"dsn\": \"test-db\"}")
				_ = os.Setenv("GONE_CUSTOMER_CLUSTER_SLAVES", "[{\n\t\"driver-name\":\"sqlite\",\n\t\"dsn\": \"test-db\"},{\n\t\"driver-name\":\"sqlite\",\n\t\"dsn\": \"test-db\"}]")

				mock.ExpectPing()
				mock.ExpectPing()
				mock.ExpectPing()

				return func() {
					_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_ENABLE")
					_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_MASTER")
					_ = os.Unsetenv("GONE_CUSTOMER_CLUSTER_SLAVES")
				}
			},
			runFn: func(in struct {
				eng0 Engine `gone:"xorm,db=customer,slave=2"`
			}) {
			},
			panicWant: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.before()()
			defer func() {
				a := recover()
				if tt.panicWant && a == nil {
					t.Errorf("panic want but not panic")
				} else if !tt.panicWant && a != nil {
					t.Errorf("panic not want but panic:%v", a)
				}
			}()
			gone.
				NewApp(Load).
				Run(tt.runFn)
		})
	}
}
