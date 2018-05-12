package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/nest"
)

func main() {
	var db, _ = sql.Open("mysql", "root:smok2015@tcp(192.168.192.250:3306)/v3?parseTime=true")
	var m = nest.NewManager(db, "org_department")
	var s = nest.NewService(m)

	var categoryList []*nest.BasicModel
	err := s.GetNodeList(1, 0, 0, &categoryList)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	for _, category := range categoryList {
		fmt.Println(category.Type, category.Id, category.IsLeafNode(), category.Name, category.LeftValue, category.RightValue)
	}
}
