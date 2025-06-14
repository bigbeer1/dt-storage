package tdenginex

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/taosdata/driver-go/v3/taosRestful"
	_ "github.com/taosdata/driver-go/v3/taosWS"
	"reflect"
	"time"
)

var (
	ErrNotFoundTable = "[0x2662] Fail to get table info, error: Table does not exist"
)

type TDengineConfig struct {
	Host     string `json:"Host"`
	Port     string `json:"Prot"`
	UserName string `json:"UserName"`
	Pass     string `json:"Pass"`
	ConFun   int64  `json:"ConFun"`
}

type TdDb struct {
	TableName string // 超级表名称
	DbName    string // 子数据库名称
}

func (t TDengineConfig) NewTDengineManager() (taos *sql.DB) {

	switch t.ConFun {
	case 1: //websocket 连接
		taosUri := fmt.Sprintf("%s:%v@ws(%v:%v)/", t.UserName, t.Pass, t.Host, t.Port)
		taos, err := sql.Open("taosWS", taosUri)
		if err != nil {
			fmt.Println("TDengine连接失败:" + err.Error())
			return nil
		}
		return taos
	case 2: // 2:http 连接
		taosUri := fmt.Sprintf("%s:%v@http(%v:%v)/", t.UserName, t.Pass, t.Host, t.Port)
		taos, err := sql.Open("taosRestful", taosUri)
		if err != nil {
			fmt.Println("TDengine连接失败:" + err.Error())
			return nil
		}
		return taos
	case 3:

	default:
		fmt.Println("td连接类型错误不是1:websocket 2:http  3:sql/tcp:")
		return nil

	}
	return nil
}

func Scan(rows *sql.Rows, Dest any) error {
	columns, _ := rows.Columns()
	values := make([]any, len(columns))

	switch dest := Dest.(type) {
	case map[string]any, *map[string]any:
		if rows.Next() {
			columnTypes, _ := rows.ColumnTypes()
			prepareValues(values, columnTypes, columns)
			if err := rows.Scan(values...); err != nil {
				return err
			}

			mapValue, ok := dest.(map[string]any)
			if !ok {
				if v, ok := dest.(*map[string]any); ok {
					mapValue = *v
				}
			}
			scanIntoMap(mapValue, values, columns)
		}
	case *[]map[string]any:
		columnTypes, _ := rows.ColumnTypes()
		for rows.Next() {
			prepareValues(values, columnTypes, columns)
			if err := rows.Scan(values...); err != nil {
				return err
			}

			mapValue := map[string]any{}
			scanIntoMap(mapValue, values, columns)
			*dest = append(*dest, mapValue)
		}
	case *int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64, *uintptr,
		*float32, *float64,
		*bool, *string, *time.Time,
		*sql.NullInt32, *sql.NullInt64, *sql.NullFloat64,
		*sql.NullBool, *sql.NullString, *sql.NullTime:
		for rows.Next() {
			if err := rows.Scan(dest); err != nil {
				return err
			}
		}
	default:
		return errors.New("not support type")
	}
	return nil
}

func scanIntoMap(mapValue map[string]any, values []any, columns []string) {
	for idx, column := range columns {
		if reflectValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(values[idx]))); reflectValue.IsValid() {
			mapValue[column] = reflectValue.Interface()
			if valuer, ok := mapValue[column].(driver.Valuer); ok {
				mapValue[column], _ = valuer.Value()
			} else if b, ok := mapValue[column].(sql.RawBytes); ok {
				mapValue[column] = string(b)
			}
		} else {
			mapValue[column] = nil
		}
	}
}

func prepareValues(values []any, columnTypes []*sql.ColumnType, columns []string) {
	if len(columnTypes) > 0 {
		for idx, columnType := range columnTypes {
			if columnType.ScanType() != nil {
				values[idx] = reflect.New(reflect.PtrTo(columnType.ScanType())).Interface()
			} else {
				values[idx] = new(any)
			}
		}
	} else {
		for idx := range columns {
			values[idx] = new(any)
		}
	}
}
