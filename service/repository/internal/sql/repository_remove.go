package sql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"time"
)

func (this *Repository) removeNode(ctx, id int64) (err error) {
	var node *nest.Node
	if node, err = this.getNodeWithId(ctx, id); err != nil {
		return err
	}

	if node == nil {
		return nest.ErrNodeNotExist
	}

	var now = time.Now()

	// 更新当前节点及其子节点的状态为删除状态
	var ub = dbs.NewUpdateBuilder()
	ub.UseDialect(this.Dialect)
	ub.Table(this.Table)
	ub.SET("status", Delete)
	ub.SET("updated_on", now)
	ub.Where("ctx = ? AND status != ? AND left_value >= ? AND right_value <= ?", node.Ctx, Delete, node.LeftValue, node.RightValue)
	if _, err := ub.Exec(this.DB); err != nil {
		return err
	}

	// 把被删除的节点及其子节点占用的空间从原树中删除掉
	var nodeLen = node.RightValue - node.LeftValue + 1
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.UseDialect(this.Dialect)
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value - ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND status != ? AND left_value > ?", node.Ctx, Delete, node.RightValue)
	if _, err = ubTreeLeft.Exec(this.DB); err != nil {
		return err
	}
	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.UseDialect(this.Dialect)
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value - ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND status != ? AND right_value > ?", node.Ctx, Delete, node.RightValue)
	if _, err = ubTreeRight.Exec(this.DB); err != nil {
		return err
	}

	return nil
}
