package xorm

import "time"

type ClusterNodeConf struct {
	DriverName string `properties:"driver-name,default=" mapstructure:"driver-name" json:"driver-name"`
	DSN        string `properties:"dsn,default=" mapstructure:"dsn,default=yyyy" json:"dsn"`
}

type Conf struct {
	DriverName    string        `properties:"driver-name,default=" mapstructure:"driver-name" json:"driver-name"`
	Dsn           string        `properties:"dsn,default=" mapstructure:"dsn" json:"dsn"`
	MaxIdleCount  int           `properties:"max-idle-count,default=5" mapstructure:"max-idle-count,default=5" json:"max-idle-count"`
	MaxOpen       int           `properties:"max-open,default=20" mapstructure:"max-open,default=20" json:"max-open"`
	MaxLifetime   time.Duration `properties:"max-lifetime,default=10m" mapstructure:"max-lifetime,default=10m" json:"max-lifetime"`
	ShowSql       bool          `properties:"show-sql,default=true" mapstructure:"show-sql,default=true" json:"show-sql"`
	PingAfterInit bool          `properties:"ping-after-init,default=false" mapstructure:"ping-after-init,default=false" json:"ping-after-init"`
}
