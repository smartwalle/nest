package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"time"
)

// addNode 添加节点
// id: 指定节点的 id，如果值小于等于0，则表示自增；
// ctx: 节点组标识；
// position:
// 		1、将新的节点添加到参照节点的子节点列表头部；
// 		2、将新的节点添加到参照节点的子节点列表尾部；
// 		3、将新的节点添加到参照节点的左边；
// 		4、将新的节点添加到参照节点的右边；
// rId: 参照节点 id，如果值等于 0，则表示添加顶级节点；
// name: 节点名
// status: 节点状态
func (this *nestRepository) addNode(ctx int64, position nest.Position, rId int64, name string, status nest.Status) (result int64, err error) {
	// 查询出参照节点的信息
	var referNode *nest.Node

	if position == nest.Root {
		// 如果是添加顶级节点，那么参照节点为 right value 最大的
		if referNode, err = this.getMaxRightNode(ctx); err != nil {
			return 0, err
		}

		// 如果参照节点为 nil，则创建一个虚拟的
		if referNode == nil {
			referNode = &nest.Node{}
			referNode.Id = -1
			referNode.Ctx = ctx
			referNode.LeftValue = 0
			referNode.RightValue = 0
			referNode.Depth = 1
		}
	} else {
		if referNode, err = this.getNodeWithId(ctx, rId); err != nil {
			return 0, err
		}
		if referNode == nil {
			return 0, nest.ErrNodeNotExist
		}
	}

	if result, err = this.addNodeWithPosition(referNode, position, name, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) addNodeWithPosition(refer *nest.Node, position nest.Position, name string, status nest.Status) (result int64, err error) {
	switch position {
	case nest.Root:
		return this.insertNodeToRoot(refer, name, status)
	case nest.First:
		return this.insertNodeToFirst(refer, name, status)
	case nest.Last:
		return this.insertNodeToLast(refer, name, status)
	case nest.Left:
		return this.insertNodeToLeft(refer, name, status)
	case nest.Right:
		return this.insertNodeToRight(refer, name, status)
	}
	return 0, nest.ErrUnknownPosition
}

func (this *nestRepository) insertNodeToRoot(refer *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ctx = refer.Ctx
	var leftValue = refer.RightValue + 1
	var rightValue = refer.RightValue + 2
	var depth = refer.Depth
	if result, err = this.insertNode(ctx, name, leftValue, rightValue, depth, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) insertNodeToFirst(refer *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", refer.Ctx, refer.LeftValue)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value > ?", refer.Ctx, refer.LeftValue)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(refer.Ctx, name, refer.LeftValue+1, refer.LeftValue+2, refer.Depth+1, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) insertNodeToLast(refer *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", refer.Ctx, refer.RightValue)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value >= ?", refer.Ctx, refer.RightValue)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(refer.Ctx, name, refer.RightValue, refer.RightValue+1, refer.Depth+1, status); err != nil {
		return 0, err
	}

	return result, nil
}

func (this *nestRepository) insertNodeToLeft(refer *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value >= ?", refer.Ctx, refer.LeftValue)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value >= ?", refer.Ctx, refer.LeftValue)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(refer.Ctx, name, refer.LeftValue, refer.LeftValue+1, refer.Depth, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) insertNodeToRight(refer *nest.Node, name string, status nest.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", refer.Ctx, refer.RightValue)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value > ?", refer.Ctx, refer.RightValue)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(refer.Ctx, name, refer.RightValue+1, refer.RightValue+2, refer.Depth, status); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) insertNode(ctx int64, name string, leftValue, rightValue, depth int, status nest.Status) (result int64, err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.table)
	ib.SET("ctx", ctx)
	ib.SET("name", name)
	ib.SET("left_value", leftValue)
	ib.SET("right_value", rightValue)
	ib.SET("depth", depth)
	ib.SET("status", status)
	ib.SET("created_on", now)
	ib.SET("updated_on", now)
	sResult, err := ib.Exec(this.db)
	if err != nil {
		return 0, err
	}
	result, _ = sResult.LastInsertId()
	return result, err
}
