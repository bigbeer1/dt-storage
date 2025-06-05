package model

import (
	"context"
	"database/sql"
	"dt-storage/common/tdenginex"
	"fmt"
	"time"
)

type TdMonitor struct {
	Ts   time.Time `json:"ts"`   // 创建时间
	Data float64   `json:"data"` // 数据
}

func (t *TdMonitor) Insert(ctx context.Context, taos *sql.DB, tddb *tdenginex.TdDb) error {

	// 拼接请求数据库和表
	dbName := fmt.Sprintf("INSERT INTO %s USING %s ", tddb.DbName, tddb.TableName)

	// 拼接参数
	tableData := fmt.Sprintf(" Tags('%v') (`ts`,`data`)", 1)

	value := fmt.Sprintf(" values ('%v',%v);", t.Ts.Format(time.RFC3339Nano), t.Data)

	sqlx := dbName + tableData + value

	_, err := taos.ExecContext(ctx, sqlx)
	if err != nil {
		return err
	}
	return nil
}
