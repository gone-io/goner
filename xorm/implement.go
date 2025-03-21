package xorm

import (
	"github.com/gone-io/gone/v2"
	"io"
	"time"
	"xorm.io/xorm"
)

func newWrappedEngine() *wrappedEngine {
	return &wrappedEngine{
		newFunc:    newEngine,
		newSession: newSession,
	}
}

func newEngine(driverName string, dataSourceName string) (xorm.EngineInterface, error) {
	return xorm.NewEngine(driverName, dataSourceName)
}

func newSession(eng xorm.EngineInterface) XInterface {
	return eng.NewSession()
}

type ClusterNodeConf struct {
	DriverName string `properties:"driver-name,default=" mapstructure:"driver-name"`
	DSN        string `properties:"dsn,default=" mapstructure:"dsn,default=yyyy"`
}

type Conf struct {
	DriverName   string        `properties:"driver-name,default=" mapstructure:"driver-name"`
	Dsn          string        `properties:"dsn,default=" mapstructure:"dsn"`
	MaxIdleCount int           `properties:"max-idle-count,default=5" mapstructure:"max-idle-count,default=5"`
	MaxOpen      int           `properties:"max-open,default=20" mapstructure:"max-open,default=20"`
	MaxLifetime  time.Duration `properties:"max-lifetime,default=10m" mapstructure:"max-lifetime,default=10m"`
	ShowSql      bool          `properties:"show-sql,default=true" mapstructure:"show-sql,default=true"`
}

//go:generate mockgen -package xorm -destination=./engine_mock_test.go xorm.io/xorm EngineInterface
type wrappedEngine struct {
	gone.Flag
	xorm.EngineInterface
	group *xorm.EngineGroup

	newFunc    func(driverName string, dataSourceName string) (xorm.EngineInterface, error)
	newSession func(xorm.EngineInterface) XInterface

	log           gone.Logger        `gone:"gone-logger"`
	conf          Conf               `gone:"config,database"`
	enableCluster bool               `gone:"config,database.cluster.enable,default=false"`
	masterConf    *ClusterNodeConf   `gone:"config,database.cluster.master"`
	slavesConf    []*ClusterNodeConf `gone:"config,database.cluster.slaves"`

	policy   xorm.GroupPolicy
	unitTest bool
}

func (e *wrappedEngine) GetOriginEngine() xorm.EngineInterface {
	return e.EngineInterface
}

func (e *wrappedEngine) SetPolicy(policy xorm.GroupPolicy) {
	e.policy = policy
	if e.group != nil {
		e.group.SetPolicy(policy)
	}
}

func (e *wrappedEngine) Start() error {
	err := e.create()
	if err != nil {
		return err
	}
	if e.unitTest {
		return nil
	}
	e.config()
	return e.Ping()
}
func (e *wrappedEngine) create() error {
	if e.EngineInterface != nil {
		return gone.NewInnerError("duplicate call Start()", gone.StartError)
	}

	if e.enableCluster {
		if e.masterConf == nil {
			return gone.NewInnerError("master config(database.cluster.master) is nil", gone.StartError)
		}

		if len(e.slavesConf) == 0 {
			return gone.NewInnerError("slaves config(database.cluster.slaves) is nil", gone.StartError)
		}

		master, err := e.newFunc(e.masterConf.DriverName, e.masterConf.DSN)
		if err != nil {
			return gone.NewInnerError(err.Error(), gone.StartError)
		}

		slaves := make([]*xorm.Engine, 0, len(e.slavesConf))
		for _, slave := range e.slavesConf {
			slaveEngine, err := e.newFunc(slave.DriverName, slave.DSN)
			if err != nil {
				return gone.NewInnerError(err.Error(), gone.StartError)
			}
			slaves = append(slaves, slaveEngine.(*xorm.Engine))
		}

		e.group, err = xorm.NewEngineGroup(master, slaves, e.policy)
		if err != nil {
			return gone.NewInnerError(err.Error(), gone.StartError)
		}
		e.EngineInterface = e.group
	} else {
		var err error
		e.EngineInterface, err = e.newFunc(e.conf.DriverName, e.conf.Dsn)
		if err != nil {
			return gone.NewInnerError(err.Error(), gone.StartError)
		}
	}
	return nil
}

func (e *wrappedEngine) config() {
	e.SetConnMaxLifetime(e.conf.MaxLifetime)
	e.SetMaxOpenConns(e.conf.MaxOpen)
	e.SetMaxIdleConns(e.conf.MaxIdleCount)
	e.SetLogger(&dbLogger{Logger: e.log, showSql: e.conf.ShowSql})
}

func (e *wrappedEngine) Stop() error {
	if e.unitTest {
		return nil
	}
	return e.EngineInterface.(io.Closer).Close()
}

func (e *wrappedEngine) Sqlx(sql string, args ...any) *xorm.Session {
	sql, args = sqlDeal(sql, args...)
	return e.SQL(sql, args...)
}
