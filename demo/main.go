package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
)

func main() {
	var db, _ = dbs.NewSQL("mysql", "root:yangfeng@tcp(192.168.1.111:3306)/v3?parseTime=true", 2, 2)
	var m = nest.NewManager(db, "org_department")

	var nodeList []*nest.Node
	m.GetSubNodeList(0, 0, 1, &nodeList)

	for _, node := range nodeList {
		fmt.Println(node.Ctx, node.Id, node.IsLeaf(), node.Name, node.LeftValue, node.RightValue)
	}
}
