package sql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"strings"
)

const Delete nest.Status = -1 // 删除

type Repository struct {
	db          dbs.DB
	dialect     dbs.Dialect
	idGenerator dbs.IdGenerator
	table       string
}

func NewRepository(db dbs.DB, dialect dbs.Dialect, table string) *Repository {
	var r = &Repository{}
	r.db = db
	r.idGenerator = dbs.GetIdGenerator()
	r.dialect = dialect

	table = strings.TrimSpace(table)
	if table == "" {
		table = "nest"
	}
	r.table = table
	return r
}

func (this *Repository) DB() dbs.DB {
	return this.db
}

func (this *Repository) Dialect() dbs.Dialect {
	return this.dialect
}

func (this *Repository) BeginTx() (dbs.TX, *Repository) {
	var tx = dbs.MustTx(this.db)
	var nRepo = *this
	nRepo.db = tx
	return tx, &nRepo
}

func (this *Repository) WithTx(tx dbs.TX) *Repository {
	var nRepo = *this
	nRepo.db = tx
	return &nRepo
}

func (this *Repository) UseIdGenerator(g dbs.IdGenerator) {
	this.idGenerator = g
}

func (this *Repository) IdGenerator() dbs.IdGenerator {
	return this.idGenerator
}

func (this *Repository) Table() string {
	return this.table
}
