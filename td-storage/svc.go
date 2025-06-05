package main

import (
	"database/sql"
)

type ServiceContext struct {
	Config Config
	// 时序数据库
	Taos *sql.DB
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Taos:   c.TDengineConfig.NewTDengineManager(),
	}
}
