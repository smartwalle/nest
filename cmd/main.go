package main

import (
	"fmt"
	"github.com/smartwalle/nest/service/repository/mysql"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
)

// 如需要测试，在 mod 中添加以下依赖
//github.com/go-sql-driver/mysql v1.4.1
//github.com/lib/pq v1.3.0
//
// import 中添加以下包导入
// _ "github.com/go-sql-driver/mysql"
// _ "github.com/lib/pq"

func main() {
	//var pdb, _ = dbs.NewSQL("postgres", "host=localhost port=5432 user=yang password=111 dbname=yang sslmode=disable", 10, 1)
	//var p = postgresql.NewRepository(pdb, "xxx")

	var mdb, _ = dbs.NewSQL("mysql", "root:yangfeng@tcp(127.0.0.1:3306)/tt?parseTime=true", 10, 1)
	var m = mysql.NewRepository(mdb, "xxx")

	var ctx int64 = 1
	var nId int64 = 0

	nId, _ = m.AddRoot(ctx, "新闻分类", nest.Enable)
	sId, _ := m.AddToLast(ctx, nId, "体育", nest.Enable)

	m.AddToLast(ctx, sId, "足球", nest.Enable)
	m.AddToLast(ctx, sId, "篮球", nest.Enable)
	m.AddToLast(ctx, sId, "排球", nest.Enable)

	m.AddToLast(ctx, nId, "时尚", nest.Enable)
	m.AddToLast(ctx, nId, "亲子", nest.Enable)
	m.AddToLast(ctx, nId, "艺术", nest.Enable)
	m.AddToLast(ctx, nId, "星座", nest.Enable)
	m.AddToLast(ctx, nId, "教育", nest.Enable)

	nId, _ = m.AddRoot(ctx, "音乐分类", nest.Enable)
	m.AddToFirst(ctx, nId, "摇滚", nest.Enable)
	m.AddToLast(ctx, nId, "流行", nest.Enable)
	m.AddToLast(ctx, nId, "乡村", nest.Enable)

	fmt.Println(m.GetLastNode(ctx, nId))
}
