package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/nest"
)

func main() {
	var db, _ = sql.Open("mysql", "")
	var m = nest.NewManager(db, "category")

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
