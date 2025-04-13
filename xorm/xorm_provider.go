package xorm

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"xorm.io/xorm"
)

type xormProvider struct {
	gone.Flag
	configure gone.Configure   `gone:"configure"`
	logger    gone.Logger      `gone:"*"`
	policy    xorm.GroupPolicy `gone:"*" option:"allowNil"`

	dbMap map[string]xorm.EngineInterface
}

func (s *xormProvider) Init() {
	s.dbMap = make(map[string]xorm.EngineInterface)
}

const dbKey = "db"
const defaultDbName = "database"
const masterKey = "master"
const slaveKey = "slave"

func getDbName(tag string) string {
	m, _ := gone.TagStringParse(tag)
	dbName := m[dbKey]
	if dbName == "" {
		dbName = defaultDbName
	}
	return dbName
}

// Provide use tag `gone:""` to get config by key('database') to create xorm.EngineInterface
// use tag `gone:"db=dbname"` to get config by key('dbname') to create xorm.EngineInterface
func (s *xormProvider) Provide(tag string) (eng xorm.EngineInterface, err error) {
	dbName := getDbName(tag)

	if eng = s.dbMap[dbName]; eng == nil {
		eng, err = s.configAndInitDb(dbName)
		if err != nil {
			return nil, gone.ToError(err)
		}
		s.dbMap[dbName] = eng
	}
	return eng, nil
}

func (s *xormProvider) configAndInitDb(dbName string) (eng xorm.EngineInterface, err error) {
	var config Conf
	var enableCluster bool

	_ = s.configure.Get(dbName, &config, "")
	_ = s.configure.Get(dbName+".cluster.enable", &enableCluster, "false")

	if !enableCluster {
		eng, err = xorm.NewEngine(config.DriverName, config.Dsn)
		if err != nil {
			return nil, gone.ToErrorWithMsg(err, "failed to create engine for db: "+dbName)
		}
	} else {
		var masterConf ClusterNodeConf
		var slavesConf []ClusterNodeConf
		_ = s.configure.Get(dbName+".cluster.master", &masterConf, "")
		_ = s.configure.Get(dbName+".cluster.slaves", &slavesConf, "")

		master, err := xorm.NewEngine(masterConf.DriverName, masterConf.DSN)
		if err != nil {
			return nil, gone.ToErrorWithMsg(err, "failed to create master engine for db: "+dbName)
		}

		slaves := make([]*xorm.Engine, 0, len(slavesConf))
		for _, slave := range slavesConf {
			slaveEngine, err := xorm.NewEngine(slave.DriverName, slave.DSN)
			if err != nil {
				return nil, gone.ToErrorWithMsg(err, "failed to create slave engine for db: "+dbName)
			}
			slaves = append(slaves, slaveEngine)
		}
		eng, err = xorm.NewEngineGroup(master, slaves, s.policy)
		if err != nil {
			return nil, gone.ToErrorWithMsg(err, "failed to create engine group for db: "+dbName)
		}
	}

	if config.MaxIdleCount > 0 {
		eng.SetMaxIdleConns(config.MaxIdleCount)
	}
	if config.MaxOpen > 0 {
		eng.SetMaxOpenConns(config.MaxOpen)
	}
	if config.MaxLifetime > 0 {
		eng.SetConnMaxLifetime(config.MaxLifetime)
	}
	eng.ShowSQL(config.ShowSql)
	eng.SetLogger(&dbLogger{Logger: s.logger, showSql: config.ShowSql})
	if config.PingAfterInit {
		if err = eng.Ping(); err != nil {
			return nil, gone.ToErrorWithMsg(err, "failed to ping db: "+dbName)
		}
	}
	return eng, nil
}

func (s *xormProvider) ProvideEngine(tagConf string) (*xorm.Engine, error) {
	e, err := s.Provide(tagConf)
	if err != nil {
		return nil, gone.ToError(err)
	}
	if x, ok := e.(*xorm.Engine); !ok {
		return nil, gone.ToError(fmt.Sprintf("db(%s) enabled cluster mod, try to use `*xorm.EngineGroup` to receive value", getDbName(tagConf)))
	} else {
		return x, nil
	}
}

func (s *xormProvider) ProvideEngineGroup(tagConf string) (*xorm.EngineGroup, error) {
	e, err := s.Provide(tagConf)
	if err != nil {
		return nil, gone.ToError(err)
	}
	if x, ok := e.(*xorm.EngineGroup); !ok {
		return nil, gone.ToError(fmt.Sprintf("db(%s) not enabled cluster mod, try to use `*xorm.Engine` to receive value", getDbName(tagConf)))
	} else {
		return x, nil
	}
}

var xormEngineProvider = gone.WrapFunctionProvider(func(tagConf string, param struct {
	xormProvider *xormProvider `gone:"*"`
}) (*xorm.Engine, error) {
	return param.xormProvider.ProvideEngine(tagConf)
})

var xormGroupProvider = gone.WrapFunctionProvider(func(tagConf string, param struct {
	xormProvider *xormProvider `gone:"*"`
}) (*xorm.EngineGroup, error) {
	return param.xormProvider.ProvideEngineGroup(tagConf)
})
