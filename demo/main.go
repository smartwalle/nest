package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
)

func main() {
	var db, _ = dbs.NewSQL("mysql", "", 10, 1)
	var m = nest.NewManager(db.GetSession(), "org_department")

	m.AddRoot(100, "r", 1000)
	m.AddRoot(100, "r", 1000)

	//var nodeList []*nest.Node
	//err := m.GetSubNodeList(2, 0, 1, &nodeList)
	//if err != nil {
	//	fmt.Println("err", err)
	//	return
	//}
	//
	//for _, node := range nodeList {
	//	fmt.Println(node.Type, node.Id, node.IsLeaf(), node.Name, node.LeftValue, node.RightValue)
	//}
}
