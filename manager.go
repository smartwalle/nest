package nest

import (
	"github.com/smartwalle/dbs"
)

type Manager struct {
	DB           dbs.DB
	Table        string
	SelectFields []string
}

func NewManager(db dbs.DB, table string) *Manager {
	var m = &Manager{}
	m.DB = db
	m.Table = table
	m.SelectFields = []string{"c.id", "c.type", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on"}
	return m
}

func (this *Manager) lockTable() {
	this.DB.ExecRaw("LOCK TABLES `" + this.Table + "` WRITE;")
}

func (this *Manager) unlockTable() {
	this.DB.ExecRaw("UNLOCK TABLES;")
}