package main

import (
	"database/sql"
	"github.com/panjf2000/ants/v2"
)

type ServiceContext struct {
	Config Config
	// 时序数据库
	Taos *sql.DB
	// 数据协程池
	DataAntsPool *ants.PoolWithFunc
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Taos:   c.TDengineConfig.NewTDengineManager(),
	}
}
