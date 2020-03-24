package sql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"strings"
)

const Delete nest.Status = -1 // 删除

type Repository struct {
	DB          dbs.DB
	Table       string
	Dialect     dbs.Dialect
	IdGenerator dbs.IdGenerator
}

func NewRepository(db dbs.DB, dialect dbs.Dialect, table string) *Repository {
	var r = &Repository{}
	r.DB = db

	table = strings.TrimSpace(table)
	if table == "" {
		table = "nest"
	}

	r.Table = table
	r.IdGenerator = dbs.GetIdGenerator()
	r.Dialect = dialect
	return r
}

func (this *Repository) BeginTx() (dbs.TX, nest.Repository) {
	var tx = dbs.MustTx(this.DB)
	var nRepo = *this
	nRepo.DB = tx
	return tx, &nRepo
}

func (this *Repository) WithTx(tx dbs.TX) nest.Repository {
	var nRepo = *this
	nRepo.DB = tx
	return &nRepo
}

func (this *Repository) UseIdGenerator(g dbs.IdGenerator) {
	this.IdGenerator = g
}
