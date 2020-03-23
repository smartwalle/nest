package sql

import (
	"fmt"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"strings"
)

const kDelete nest.Status = -1 // 删除

type nestRepository struct {
	db    dbs.DB
	table string
}

func NewRepository(db dbs.DB, table string) nest.Repository {
	var r = &nestRepository{}
	r.db = db
	r.table = table

	if err := r.initTable(); err != nil {
		panic(fmt.Sprintf("创建 %s 失败, 错误信息为: %v", table, err))
	}
	return r
}

func (this *nestRepository) BeginTx() (dbs.TX, nest.Repository) {
	var tx = dbs.MustTx(this.db)
	var nRepo = *this
	nRepo.db = tx
	return tx, &nRepo
}

func (this *nestRepository) WithTx(tx dbs.TX) nest.Repository {
	var nRepo = *this
	nRepo.db = tx
	return &nRepo
}

func (this *nestRepository) initTable() error {
	var tx = dbs.MustTx(this.db)

	var cb = dbs.NewBuilder("")

	if cb.GetDialect() == dbs.DialectPostgreSQL {
		cb.Format(strings.ReplaceAll(initPostgreSQLTable, "nest", this.table))
	} else {
		cb.Format(initMySQLTable, this.table, this.table, this.table, this.table, this.table, this.table)
	}
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	tx.Commit()
	return nil
}
