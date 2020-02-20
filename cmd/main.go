package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest/service/repository/mysql"
)

func main() {
	var db, _ = dbs.NewSQL("mysql", "root:yangfen@tcp(127.0.0.1:3306)/tt?parseTime=true", 10, 1)
	var m = mysql.NewRepository(db, "xd")

	var ctx int64 = 1
	//m.AddRoot(ctx, "新闻分类", nest.Enable)
	//m.AddToLast(ctx, 1, "体育", nest.Enable)
	//m.AddToLast(ctx, 1, "娱乐", nest.Enable)
	//m.AddToLast(ctx, 1, "亲子", nest.Enable)
	//m.AddToLast(ctx, 1, "时尚", nest.Enable)
	//m.AddToLast(ctx, 1, "艺术", nest.Enable)
	//m.AddToLast(ctx, 1, "星座", nest.Enable)
	//m.AddToLast(ctx, 1, "教育", nest.Enable)
	//m.AddToLast(ctx, 2, "足球", nest.Enable)
	//m.AddToLast(ctx, 2, "篮球", nest.Enable)
	//m.AddToLast(ctx, 2, "田径", nest.Enable)

	//fmt.Println(m.GetPreviousNode(ctx, 1))
	//fmt.Println(m.GetPreviousNode(ctx, 3))
	//fmt.Println(m.GetPreviousNode(ctx, 8))
	fmt.Println(m.GetNextNode(ctx, 2))
	fmt.Println(m.GetNextNode(ctx, 9))

	//ctx = 2
	//m.AddRootWithId(ctx, 12, "商品分类", nest.Enable)
	//m.AddToLastWithId(ctx, 13, 12, "图书", nest.Enable)
	//m.AddToLastWithId(ctx, 14, 12, "家电", nest.Enable)
	//m.AddToLastWithId(ctx, 15, 12, "家具", nest.Enable)
	//m.AddToLastWithId(ctx, 16, 12, "食品", nest.Enable)
	//m.AddToLastWithId(ctx, 17, 12, "服装", nest.Enable)
	//m.AddToLastWithId(ctx, 18, 17, "女装", nest.Enable)
	//m.AddToLastWithId(ctx, 19, 17, "男装", nest.Enable)
	//m.AddToLastWithId(ctx, 20, 17, "童装", nest.Enable)

	//var nodeList, err = m.GetNodeList(ctx, 1, nest.Enable, 0, "", 0, 0, false)
	//fmt.Println(err)
	//
	//for _, node := range nodeList {
	//	fmt.Println(node.Ctx, node.Id, node.IsLeaf(), node.Name, node.LeftValue, node.RightValue)
	//}
}
