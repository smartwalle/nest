package postgresql

import (
	"fmt"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"github.com/smartwalle/nest/service/repository/internal/sql"
	"strings"
)

type repository struct {
	sql.Repository
}

func NewRepository(db dbs.DB, table string) nest.Repository {
	var r = &repository{}
	r.Repository = sql.NewRepository(db, dbs.DialectPostgreSQL, table)

	if err := r.initTable(); err != nil {
		panic(fmt.Sprintf("初始化数据库 %s 失败, 错误信息为: %v", table, err))
	}
	return r
}

func (this *repository) initTable() error {
	var rawText = "" +
		"create table if not exists nest(" +
		"id bigint not null constraint nest_pk primary key," +
		"ctx bigint default 0," +
		"name varchar(128)," +
		"left_value bigint," +
		"right_value bigint," +
		"depth integer," +
		"status integer default 1," +
		"created_on timestamp with time zone," +
		"updated_on timestamp with time zone);" +
		"create unique index if not exists nest_id_uindex on nest (id);" +
		"create index if not exists nest_ctx_index on nest (ctx);" +
		"create index if not exists nest_ctx_left_value_index on nest (ctx, left_value);" +
		"create index if not exists nest_ctx_right_value_index on nest (ctx, right_value);" +
		"create index if not exists nest_id_ctx_index on nest (id, ctx);"

	var sql = strings.ReplaceAll(rawText, "nest", this.Table())
	var rb = dbs.NewBuilder(sql)
	if _, err := rb.Exec(this.DB()); err != nil {
		return err
	}
	return nil
}

func (this *repository) BeginTx() (dbs.TX, nest.Repository) {
	var nRepo = *this
	var tx dbs.TX
	tx, nRepo.Repository = this.Repository.ExBeginTx()
	return tx, &nRepo
}

func (this *repository) WithTx(tx dbs.TX) nest.Repository {
	var nRepo = *this
	nRepo.Repository = this.Repository.ExWithTx(tx)
	return &nRepo
}
