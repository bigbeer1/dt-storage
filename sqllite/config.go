package main

import (
	"dt-storage/common/tdenginex"
)

type Config struct {
	TDengineConfig tdenginex.TDengineConfig

	LoggerNumber int64 `json:"LoggerNumber"` // 0不打印日志 1打印日志

	Limit int `json:"Limit"` // 协程并发数

	AllLimit int `json:"AllLimit"` // 总协程数

	TimeNext int64 `json:"TimeNext"` // 间隔毫秒
}
