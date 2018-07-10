package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
)

func main() {
	var db, _ = dbs.NewSQL("mysql", "root:yangfeng@tcp(192.168.1.99:3306)/test?parseTime=true", 2, 2)
	var m = nest.NewManager(db, "org_department")

	//m.AddRootWithId(1, 1001, "新闻分类", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(2, 1, "体育", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(3, 1, "娱乐", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(4, 1, "亲子", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(5, 1, "时尚", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(6, 1, "艺术", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(7, 1, "星座", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(8, 1, "教育", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(9, 2, "足球", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(10, 2, "篮球", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(11, 2, "田径", nest.K_STATUS_ENABLE)
	//
	//m.AddRootWithId(12, 1002, "商品分类", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(13, 12, "图书", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(14, 12, "家电", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(15, 12, "家具", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(16, 12, "食品", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(17, 12, "服装", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(18, 17, "女装", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(19, 17, "男装", nest.K_STATUS_ENABLE)
	//m.AddToLastWithId(20, 17, "童装", nest.K_STATUS_ENABLE)


	var nodeList []*nest.Node
	m.GetSubNodeList(1001, 2, 0, 0, &nodeList)

	for _, node := range nodeList {
		fmt.Println(node.Ctx, node.Id, node.IsLeaf(), node.Name, node.LeftValue, node.RightValue)
	}
}
