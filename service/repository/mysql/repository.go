package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest/service"
)

type nestRepository struct {
	db           dbs.DB
	table        string
	selectFields []string
}

func NewNestRepository(db dbs.DB, table string, selectFields ...string) service.NestRepository {
	var r = &nestRepository{}
	r.db = db
	r.table = table
	if len(selectFields) == 0 {
		r.selectFields = []string{"c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on"}
	} else {
		r.selectFields = selectFields
	}
	return r
}
