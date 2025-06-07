package main

import (
	"database/sql"
	"dt-storage/common/tdenginex"
	"dt-storage/td-storage/model"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/conf"
	"math/rand"
	"strconv"
	"time"
)

var configFile = flag.String("f", "td-storage.yaml", "the config file")

func main() {
	// 读取配置文件
	var c Config
	conf.MustLoad(*configFile, &c)

	fmt.Println(c)

	svcCtx := NewServiceContext(c)

	if c.AllLimit > 0 {
		svcCtx.DataAntsPool, _ = ants.NewPoolWithFunc(c.AllLimit, func(req interface{}) {
			data, _ := req.(int)
			svcCtx.Idata(data)
		})
	}

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	createTable(db)

	insertData(db, 111, time.Now())
	defer db.Close()
}

func createTable(db *sql.DB) {
	sqlStmt := `
		CREATE TABLE monitor (
    		id INTEGER PRIMARY KEY,
    		monitor_data TEXT,
    		monitor_ts INTEGER
		);
    `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}
}

func insertData(db *sql.DB, data float64, ts time.Time) {
	sqlStmt := `INSERT INTO monitor (monitor_data, monitor_ts) VALUES (?, ?)`
	_, err := db.Exec(sqlStmt, data, ts)
	if err != nil {
		panic(err)
	}
}

// 随机 写入数据
func (l ServiceContext) Idata(i int) {

	// 随机数据
	numberRandom := RandomDecimal("0", 300)
	monitor := &model.TdMonitor{
		Ts:   time.Now(),
		Data: numberRandom,
	}

	// td 插入数据库名称
	tddb := &tdenginex.TdDb{
		DbName:    "chint.d" + fmt.Sprint(i),
		TableName: "chint.monitor_point",
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := monitor.Insert(ctx, l.Taos, tddb)

	if l.Config.LoggerNumber == 1 {
		fmt.Println(fmt.Sprintf("发送一条数据:%s ", monitor))
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("数据错误:%s ", err.Error()))
	}

}

func RandomDecimal(bit string, multiple float64) float64 {
	bit = "%." + bit + "f"
	data := rand.Float64() * multiple
	data, _ = strconv.ParseFloat(fmt.Sprintf(bit, data), 64)
	return data
}
