package xorm

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"reflect"
	"strconv"
	"xorm.io/xorm"
)

const clusterKey = "db"
const defaultCluster = "database"

func newProvider(engine *wrappedEngine) gone.Goner {
	var engineMap = make(map[string]*wrappedEngine)
	engineMap[""] = engine
	engineMap[defaultCluster] = engine

	return &provider{
		engineMap: engineMap,
		newFunc:   engine.newFunc,
		unitTest:  engine.unitTest,
	}
}

type provider struct {
	gone.Flag
	engineMap map[string]*wrappedEngine

	//heaven    gone.Heaven    `gone:"*"`
	//cemetery  gone.Cemetery  `gone:"*"`
	configure gone.Configure  `gone:"*"`
	log       gone.Logger     `gone:"*"`
	before    gone.BeforeStop `gone:"*"`

	newFunc  func(driverName string, dataSourceName string) (xorm.EngineInterface, error)
	unitTest bool
}

func (p *provider) GonerName() string {
	return "xorm"
}

var xormInterface = gone.GetInterfaceType(new(XormEngine))
var xormInterfaceSlice = gone.GetInterfaceType(new([]XormEngine))

func (p *provider) Provide(tagConf string, t reflect.Type) (any, error) {
	m, _ := gone.TagStringParse(tagConf)
	clusterName := m[clusterKey]
	if clusterName == "" {
		clusterName = defaultCluster
	}

	db, err := p.getDb(clusterName)
	if err != nil {
		return nil, gone.ToError(err)
	}
	if t == xormInterfaceSlice {
		if !db.enableCluster {
			return nil, gone.NewInnerError(fmt.Sprintf("database(name=%s) is not enable cluster, cannot inject []gone.XormEngine", clusterName), gone.InjectError)
		}

		engines := db.group.Slaves()
		xormEngines := make([]XormEngine, 0, len(engines))
		for _, eng := range engines {
			xormEngines = append(xormEngines, &wrappedEngine{
				EngineInterface: eng,
			})
		}
		return xormEngines, nil
	}

	if t == xormInterface {
		if _, ok := m["master"]; ok {
			if !db.enableCluster {
				return nil, gone.NewInnerError(fmt.Sprintf("database(name=%s) is not enable cluster, cannot inject master into gone.XormEngine", clusterName), gone.InjectError)
			}

			return &wrappedEngine{
				EngineInterface: db.group.Master(),
			}, nil
		}

		if slaveIndex, ok := m["slave"]; ok {
			if !db.enableCluster {
				return nil, gone.NewInnerError(fmt.Sprintf("database(name=%s) is not enable cluster, cannot inject slave into gone.XormEngine", clusterName), gone.InjectError)
			}

			slaves := db.group.Slaves()
			var index int64
			var err error
			if slaveIndex != "" {
				index, err = strconv.ParseInt(slaveIndex, 10, 64)
				if err != nil || index < 0 || index >= int64(len(slaves)) {
					return nil, gone.NewInnerError(fmt.Sprintf("invalid slave index: %s, must be greater than or equal to 0 and less than %d ", slaveIndex, len(slaves)), gone.InjectError)
				}
			}

			return &wrappedEngine{
				EngineInterface: slaves[index],
			}, nil
		}
		return db, nil
	}
	return nil, gone.NewInnerErrorWithParams(
		gone.GonerTypeNotMatch,
		"Cannot find matched value for %q",
		gone.GetTypeName(t),
	)
}

func (p *provider) getDb(clusterName string) (*wrappedEngine, error) {
	db := p.engineMap[clusterName]
	if db == nil {
		var config Conf
		err := p.configure.Get(clusterName, &config, "")
		if err != nil {
			return nil, gone.NewInnerError("failed to get config for cluster: "+clusterName, gone.InjectError)
		}

		var enableCluster bool
		err = p.configure.Get(clusterName+".cluster.enable", &enableCluster, "false")
		if err != nil {
			return nil, gone.NewInnerError("failed to get cluster enable config for cluster: "+clusterName, gone.InjectError)
		}

		var masterConf ClusterNodeConf
		err = p.configure.Get(clusterName+".cluster.master", &masterConf, "")
		if err != nil {
			return nil, gone.NewInnerError("failed to get master config for cluster: "+clusterName, gone.InjectError)
		}

		var slavesConf []*ClusterNodeConf
		err = p.configure.Get(clusterName+".cluster.slaves", &slavesConf, "")
		if err != nil {
			return nil, gone.NewInnerError("failed to get slaves config for cluster: "+clusterName, gone.InjectError)
		}

		db = newWrappedEngine()
		db.conf = config
		db.enableCluster = enableCluster
		db.masterConf = &masterConf
		db.slavesConf = slavesConf

		//for test
		db.newFunc = p.newFunc
		db.unitTest = p.unitTest

		err = db.Start()
		if err != nil {
			return nil, gone.NewInnerError("failed to start xorm engine for cluster: "+clusterName, gone.InjectError)
		}

		p.before(func() {
			err := db.Stop()
			if err != nil {
				p.log.Errorf("failed to stop xorm engine for cluster(name=%s): %v", clusterName, err)
			}
		})

		p.engineMap[clusterName] = db
	}
	return db, nil
}
