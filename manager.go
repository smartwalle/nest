package nest

import (
	"github.com/smartwalle/dbs"
)

type Manager struct {
	db           dbs.DB
	table        string
	SelectFields []string
}

func NewManager(db dbs.DB, table string) *Manager {
	var m = &Manager{}
	m.db = db
	m.table = table
	m.SelectFields = []string{"c.id", "c.type", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on"}
	return m
}

type IManager interface {
	AddCategory(cId int64, cType, position int, referTo int64, name string, status int, exts ...map[string]interface{}) (result int64, err error)

	GetCategory(id int64, result interface{}) (err error)
	GetCategoryWithName(cType int, name string, result interface{}) (err error)

	GetCategoryList(parentId int64, cType, status, depth int, name string, limit uint64, includeParent bool, result interface{}) (err error)
	GetIdList(parentId int64, status, depth int, includeParent bool) (result []int64, err error)
	GetPathList(id int64, status int, includeLastNode bool, result interface{}) (err error)

	UpdateCategory(id int64, updateInfo map[string]interface{}) (err error)
	UpdateCategoryStatus(id int64, status, updateType int) (err error)
	MoveCategory(position int, id, rid int64) (err error)
}
