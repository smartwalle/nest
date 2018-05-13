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

	var categoryList []*nest.BaseModel
	err := m.GetNodePathList(1, 0, 1, &categoryList)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	for _, category := range categoryList {
		fmt.Println(category.Type, category.Id, category.IsLeafNode(), category.Name, category.LeftValue, category.RightValue)
	}
}
