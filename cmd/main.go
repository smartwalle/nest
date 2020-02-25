package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"github.com/smartwalle/nest/service/repository/mysql"
)

func main() {
	var db, _ = dbs.NewSQL("mysql", "root:yangfeng@tcp(127.0.0.1:3306)/tt?parseTime=true", 10, 1)
	var m = mysql.NewRepository(db, "xd")

	var ctx int64 = 1

	//var nId int64 = 0
	//-----
	//nId, _ = m.AddRoot(ctx, "新闻分类", nest.Enable)
	//sId, _ := m.AddToLast(ctx, nId, "体育", nest.Enable)
	//
	//m.AddToLast(ctx, sId, "足球", nest.Enable)
	//m.AddToLast(ctx, sId, "篮球", nest.Enable)
	//m.AddToLast(ctx, sId, "排球", nest.Enable)
	//
	//m.AddToLast(ctx, nId, "时尚", nest.Enable)
	//m.AddToLast(ctx, nId, "亲子", nest.Enable)
	//m.AddToLast(ctx, nId, "艺术", nest.Enable)
	//m.AddToLast(ctx, nId, "星座", nest.Enable)
	//m.AddToLast(ctx, nId, "教育", nest.Enable)
	//
	////-----
	//nId, _ = m.AddRoot(ctx, "物品分类", nest.Enable)
	//
	//m.AddToFirst(ctx, nId, "电子产品", nest.Enable)
	//m.AddToFirst(ctx, nId, "学习用具", nest.Enable)
	//
	////-----
	//nId, _ = m.AddRoot(ctx, "音乐分类", nest.Enable)
	//
	//m.AddToFirst(ctx, nId, "摇滚", nest.Enable)
	//m.AddToLast(ctx, nId, "流行", nest.Enable)
	//m.AddToLast(ctx, nId, "乡村", nest.Enable)
	//
	////-----
	//nId, _ = m.AddRoot(ctx, "电影分类", nest.Enable)
	//
	//m.AddToLast(ctx, nId, "悬疑", nest.Enable)
	//m.AddToLast(ctx, nId, "动作", nest.Enable)

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

	//fmt.Println(m.GetLastNode(ctx,2))
	//fmt.Println(m.GetFirstNode(ctx,2))

	//fmt.Println(m.MoveToRoot(ctx, 11))
	//fmt.Println(m.MoveToLast(ctx, 11, 0))

	//fmt.Println(m.GetPreviousNode(ctx, 1))
	//fmt.Println(m.GetPreviousNode(ctx, 3))
	//fmt.Println(m.GetPreviousNode(ctx, 8))
	//fmt.Println(m.GetNextNode(ctx, 2))
	//fmt.Println(m.MoveUp(ctx, 9))
	//fmt.Println(m.GetNextNode(ctx, 9))

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

	//m.MoveToLast(ctx, 11, 1)

	var nodeList, err = m.GetNodes(ctx, 0, nest.Enable, 0, "", 0, 0, false)
	fmt.Println(err)

	for _, node := range nodeList {
		for i := 0; i < node.Depth; i++ {
			fmt.Print("-")
		}
		fmt.Println(node.Ctx, node.Id, node.HasChildNodes(), node.Name, node.LeftValue, node.RightValue)
	}
}
