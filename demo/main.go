package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"fmt"
)

func main() {
	var db, _ = dbs.NewSQL("mysql", "root:smok2015@tcp(192.168.192.250:3306)/titan_dev?parseTime=true", 10, 1)
	var m = nest.NewManager(db.GetSession(), "category")

	var nodeList []*nest.Node
	err := m.GetSubNodeList(2, 0, 1, &nodeList)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	for _, node := range nodeList {
		fmt.Println(node.Type, node.Id, node.IsLeaf(), node.Name, node.LeftValue, node.RightValue)
	}
}
