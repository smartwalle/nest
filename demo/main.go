package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"fmt"
)

func main() {
	var db, _ = dbs.NewSQL("mysql", "", 2, 2)
	var m = nest.NewManager(db.GetSession(), "org_department")

	var nodeList []*nest.Node
	m.GetSubNodeList(2, 0, 1, &nodeList)

	for _, node := range nodeList {
		fmt.Println(node.Type, node.Id, node.IsLeaf(), node.Name, node.LeftValue, node.RightValue)
	}
}
