package sql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"time"
)

// addNode 添加节点
// ctx: 节点组标识；
// position:
// 		1、将新的节点添加到参照节点的子节点列表头部；
// 		2、将新的节点添加到参照节点的子节点列表尾部；
// 		3、将新的节点添加到参照节点的左边；
// 		4、将新的节点添加到参照节点的右边；
// rId: 参照节点 id；
// name: 节点名
// status: 节点状态
func (this *Repository) addNode(ctx int64, position nest.Position, rId int64, name string, status nest.Status) (result int64, err error) {
	// 查询出参照节点的信息
	var rNode *nest.Node

	if position == nest.Root {
		// 如果是添加顶级节点，那么参照节点为顶级节点列表中的最后一个节点
		if rNode, err = this.getLastNode(ctx, 0); err != nil {
			return 0, err
		}

		// 如果参照节点为 nil，则创建一个虚拟的
		if rNode == nil {
			rNode = &nest.Node{}
			rNode.Id = -1
			rNode.Ctx = ctx
			rNode.LeftValue = 0
			rNode.RightValue = 0
			rNode.Depth = 1
		}
	} else {
		if rNode, err = this.getNodeWithId(ctx, rId); err != nil {
			return 0, err
		}
		if rNode == nil {
			return 0, nest.ErrNodeNotExist
		}
	}

	if result, err = this.addNodeWithPosition(rNode, position, name, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *Repository) addNodeWithPosition(rNode *nest.Node, position nest.Position, name string, status nest.Status) (result int64, err error) {
	switch position {
	case nest.Root:
		return this.insertNodeToRoot(rNode, name, status)
	case nest.First:
		return this.insertNodeToFirst(rNode, name, status)
	case nest.Last:
		return this.insertNodeToLast(rNode, name, status)
	case nest.Left:
		return this.insertNodeToLeft(rNode, name, status)
	case nest.Right:
		return this.insertNodeToRight(rNode, name, status)
	}
	return 0, nest.ErrUnknownPosition
}

func (this *Repository) insertNodeToRoot(rNode *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ctx = rNode.Ctx
	var leftValue = rNode.RightValue + 1
	var rightValue = rNode.RightValue + 2
	var depth = rNode.Depth
	if result, err = this.insertNode(ctx, name, leftValue, rightValue, depth, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *Repository) insertNodeToFirst(rNode *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.UseDialect(this.dialect)
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", rNode.Ctx, rNode.LeftValue)
	ubLeft.Where("status != ?", Delete)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.UseDialect(this.dialect)
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value > ?", rNode.Ctx, rNode.LeftValue)
	ubRight.Where("status != ?", Delete)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(rNode.Ctx, name, rNode.LeftValue+1, rNode.LeftValue+2, rNode.Depth+1, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *Repository) insertNodeToLast(rNode *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.UseDialect(this.dialect)
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", rNode.Ctx, rNode.RightValue)
	ubLeft.Where("status != ?", Delete)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.UseDialect(this.dialect)
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value >= ?", rNode.Ctx, rNode.RightValue)
	ubRight.Where("status != ?", Delete)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(rNode.Ctx, name, rNode.RightValue, rNode.RightValue+1, rNode.Depth+1, status); err != nil {
		return 0, err
	}

	return result, nil
}

func (this *Repository) insertNodeToLeft(rNode *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.UseDialect(this.dialect)
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value >= ?", rNode.Ctx, rNode.LeftValue)
	ubLeft.Where("status != ?", Delete)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.UseDialect(this.dialect)
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value >= ?", rNode.Ctx, rNode.LeftValue)
	ubRight.Where("status != ?", Delete)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(rNode.Ctx, name, rNode.LeftValue, rNode.LeftValue+1, rNode.Depth, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *Repository) insertNodeToRight(rNode *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.UseDialect(this.dialect)
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", rNode.Ctx, rNode.RightValue)
	ubLeft.Where("status != ?", Delete)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.UseDialect(this.dialect)
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value > ?", rNode.Ctx, rNode.RightValue)
	ubRight.Where("status != ?", Delete)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(rNode.Ctx, name, rNode.RightValue+1, rNode.RightValue+2, rNode.Depth, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *Repository) insertNode(ctx int64, name string, leftValue, rightValue int64, depth int, status nest.Status) (result int64, err error) {
	var now = time.Now()
	var nId = this.idGenerator.Next()
	var ib = dbs.NewInsertBuilder()
	ib.UseDialect(this.dialect)
	ib.Table(this.table)
	ib.SET("id", nId)
	ib.SET("ctx", ctx)
	ib.SET("name", name)
	ib.SET("left_value", leftValue)
	ib.SET("right_value", rightValue)
	ib.SET("depth", depth)
	ib.SET("status", status)
	ib.SET("created_on", now)
	ib.SET("updated_on", now)
	_, err = ib.Exec(this.db)
	if err != nil {
		return 0, err
	}
	return nId, err
}
