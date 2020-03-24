package mysql

import (
	"fmt"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"github.com/smartwalle/nest/service/repository/internal/sql"
	"strings"
)

type repository struct {
	*sql.Repository
}

func NewRepository(db dbs.DB, table string) nest.Repository {
	var r = &repository{}
	r.Repository = sql.NewRepository(db, dbs.DialectMySQL, table)

	if err := r.initTable(); err != nil {
		panic(fmt.Sprintf("初始化数据库 %s 失败, 错误信息为: %v", table, err))
	}
	return r
}

func (this *repository) initTable() error {
	var rawText = "" +
		"CREATE TABLE IF NOT EXISTS `nest` (" +
		"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
		"`ctx` int(11) DEFAULT '0'," +
		"`name` varchar(128) DEFAULT NULL," +
		"`left_value` int(11) DEFAULT NULL," +
		"`right_value` int(11) DEFAULT NULL," +
		"`depth` int(11) DEFAULT NULL," +
		"`status` int(11) DEFAULT '1'," +
		"`created_on` datetime DEFAULT NULL," +
		"`updated_on` datetime DEFAULT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `nest_id_uindex` (`id`)," +
		"KEY `nest_ctx_index` (`ctx`)," +
		"KEY `nest_id_ctx_index` (`id`,`ctx`)," +
		"KEY `nest_ctx_right_value_index` (`ctx`,`right_value`)," +
		"KEY `nest_ctx_left_value_index` (`ctx`,`left_value`)" +
		");"

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
	tx, nRepo.Repository = this.Repository.BeginTx()
	return tx, &nRepo
}

func (this *repository) WithTx(tx dbs.TX) nest.Repository {
	var nRepo = *this
	nRepo.Repository = this.Repository.WithTx(tx)
	return &nRepo
}
