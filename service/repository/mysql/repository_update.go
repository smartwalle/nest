package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
	"time"
)

func (this *nestRepository) updateNodeName(ctx, id int64, name string) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.table)
	ub.SET("name", name)
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	if _, err = ub.Exec(this.db); err != nil {
		return nil
	}
	return nil
}

// updateNodeStatus 更新节点状态
// id: 被更新节点的 id
// status: 新的状态
// updateType:
// 		0、只更新当前节点的状态，子节点的状态不会受到影响，并且不会改变父子关系；
// 		1、子节点的状态会一起更新，不会改变父子关系；
// 		2、子节点的状态不会受到影响，并且所有子节点会向上移动一级（只针对把状态设置为 无效 的时候）；
func (this *nestRepository) updateNodeStatus(ctx, id int64, status nest.Status, updateType int) (err error) {
	var node *nest.Node
	if node, err = this.getNodeWithId(ctx, id); err != nil {
		return err
	}

	if node == nil {
		return nest.ErrNodeNotExist
	}

	if node.Status == status {
		return nil
	}

	var now = time.Now()

	switch updateType {
	case 2:
		if status == nest.Disable {
			var ub = dbs.NewUpdateBuilder()
			ub.Table(this.table)
			ub.SET("status", status)
			ub.SET("right_value", dbs.SQL("left_value + 1"))
			ub.SET("updated_on", now)
			ub.Where("id = ?", id)
			ub.Limit(1)
			if _, err := ub.Exec(this.db); err != nil {
				return err
			}

			var ubChild = dbs.NewUpdateBuilder()
			ubChild.Table(this.table)
			ubChild.SET("left_value", dbs.SQL("left_value + 1"))
			ubChild.SET("right_value", dbs.SQL("right_value + 1"))
			ubChild.SET("depth", dbs.SQL("depth-1"))
			ubChild.SET("updated_on", now)
			ubChild.Where("ctx = ? AND left_value > ? AND right_value < ?", node.Ctx, node.LeftValue, node.RightValue)
			if _, err := ubChild.Exec(this.db); err != nil {
				return err
			}
		}
	case 1:
		var ub = dbs.NewUpdateBuilder()
		ub.Table(this.table)
		ub.SET("status", status)
		ub.SET("updated_on", now)
		ub.Where("ctx = ? AND left_value >= ? AND right_value <= ?", node.Ctx, node.LeftValue, node.RightValue)
		if _, err := ub.Exec(this.db); err != nil {
			return err
		}
	case 0:
		var ub = dbs.NewUpdateBuilder()
		ub.Table(this.table)
		ub.SET("status", status)
		ub.SET("updated_on", now)
		ub.Where("id = ?", id)
		ub.Limit(1)
		if _, err := ub.Exec(this.db); err != nil {
			return err
		}
	}
	return nil
}

func (this *nestRepository) moveNode(position nest.Position, ctx, id, rId int64) (err error) {
	if id == rId {
		return nest.ErrParentNotAllowed
	}

	// 判断被移动的节点是否存在
	var node *nest.Node
	if node, err = this.getNodeWithId(ctx, id); err != nil {
		return err
	}
	if node == nil {
		return nest.ErrNodeNotExist
	}

	// 判断参照节点是否存在
	var refer *nest.Node
	if position == nest.Root {
		// 如果是添加顶级节点，那么参照节点为 right value 最大的
		if refer, err = this.getMaxRightNode(node.Ctx); err != nil {
			return err
		}
		if refer != nil && refer.Id == node.Id {
			return nil
		}
	} else {
		if refer, err = this.getNodeWithId(ctx, rId); err != nil {
			return err
		}
	}
	if refer == nil {
		return nest.ErrParentNotExist
	}

	// 判断被移动节点和目标参照节点是否属于同一 Ctx
	if refer.Ctx != node.Ctx {
		return nest.ErrParentNotAllowed
	}

	// 循环连接问题，即 参照节点 是 被移动节点 的子节点
	if refer.LeftValue > node.LeftValue && refer.RightValue < node.RightValue {
		return nest.ErrParentNotAllowed
	}

	// 判断是否已经是子节点
	//if refer.LeftValue < node.LeftValue && refer.RightValue > node.RightValue && node.Depth - 1 == refer.Depth {
	//	tx.Rollback()
	//	return ErrParentNotAllowed
	//}

	// 查询出被移动节点的所有子节点
	//children, err := this.getNodeList(node.Id, 0, 0)
	children, err := this.getNodeList(node.Ctx, node.Id, 0, 0, "", 0, 0, true)
	if err != nil {
		return err
	}

	var updateIdList []int64
	updateIdList = append(updateIdList, node.Id)
	for _, c := range children {
		updateIdList = append(updateIdList, c.Id)
	}

	if err = this.moveNodeWithPosition(position, node, refer, updateIdList); err != nil {
		return err
	}
	return nil
}

func (this *nestRepository) moveNodeWithPosition(position nest.Position, node, refer *nest.Node, updateIdList []int64) (err error) {
	var nodeLen = node.RightValue - node.LeftValue + 1
	var now = time.Now()

	// 把要移动的节点及其子节点从原树中删除掉
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value - ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND left_value > ?", node.Ctx, node.RightValue)
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}
	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value - ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND right_value > ?", node.Ctx, node.RightValue)
	if _, err = ubTreeRight.Exec(this.db); err != nil {
		return err
	}

	if refer.LeftValue > node.RightValue {
		refer.LeftValue -= nodeLen
	}
	if refer.RightValue > node.RightValue {
		refer.RightValue -= nodeLen
	}

	switch position {
	case nest.Root:
		return this.moveToRight(node, refer, updateIdList, nodeLen)
	case nest.First:
		return this.moveToFirst(node, refer, updateIdList, nodeLen)
	case nest.Last:
		return this.moveToLast(node, refer, updateIdList, nodeLen)
	case nest.Left:
		return this.moveToLeft(node, refer, updateIdList, nodeLen)
	case nest.Right:
		return this.moveToRight(node, refer, updateIdList, nodeLen)
	}
	return nest.ErrUnknownPosition
}

func (this *nestRepository) moveToFirst(node, parent *nest.Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND left_value > ?", parent.Ctx, parent.LeftValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND right_value > ?", parent.Ctx, parent.LeftValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeRight.Exec(this.db); err != nil {
		return err
	}

	parent.RightValue += nodeLen

	// 更新被移动节点的信息
	var diff = node.LeftValue - parent.LeftValue - 1
	var diffDepth = parent.Depth - node.Depth + 1
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	if _, err = ubTree.Exec(this.db); err != nil {
		return err
	}

	return nil
}

func (this *nestRepository) moveToLast(node, parent *nest.Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND left_value > ?", parent.Ctx, parent.RightValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND right_value >= ?", parent.Ctx, parent.RightValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeRight.Exec(this.db); err != nil {
		return err
	}

	parent.RightValue += nodeLen

	// 更新被移动节点的信息
	var diff = node.RightValue - parent.RightValue + 1
	var diffDepth = parent.Depth - node.Depth + 1
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	if _, err = ubTree.Exec(this.db); err != nil {
		return err
	}

	return nil
}

func (this *nestRepository) moveToLeft(node, refer *nest.Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND left_value >= ?", refer.Ctx, refer.LeftValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND right_value >= ?", refer.Ctx, refer.LeftValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeRight.Exec(this.db); err != nil {
		return err
	}

	//refer.LeftValue += nodeLen

	// 更新被移动节点的信息
	var diff = node.LeftValue - refer.LeftValue
	var diffDepth = refer.Depth - node.Depth
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	if _, err = ubTree.Exec(this.db); err != nil {
		return err
	}

	return nil
}

func (this *nestRepository) moveToRight(node, refer *nest.Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND left_value > ?", refer.Ctx, refer.RightValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND right_value > ?", refer.Ctx, refer.RightValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeRight.Exec(this.db); err != nil {
		return err
	}

	// 更新被移动节点的信息
	var diff = node.LeftValue - refer.RightValue - 1
	var diffDepth = refer.Depth - node.Depth
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	if _, err = ubTree.Exec(this.db); err != nil {
		return err
	}

	return nil
}
