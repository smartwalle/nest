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
	ub.Where("status != ?", kDelete)
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
// 		0、更新当前节点的状态，如果有子节点，则不能设置为无效；
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
			ub.Where("status != ?", kDelete)
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
			ubChild.Where("ctx = ? AND status != ? AND left_value > ? AND right_value < ?", node.Ctx, kDelete, node.LeftValue, node.RightValue)
			if _, err := ubChild.Exec(this.db); err != nil {
				return err
			}
		} else if status == nest.Enable {
			var ub = dbs.NewUpdateBuilder()
			ub.Table(this.table)
			ub.SET("status", status)
			ub.SET("updated_on", now)
			ub.Where("id = ?", id)
			ub.Where("status != ?", kDelete)
			ub.Limit(1)
			if _, err := ub.Exec(this.db); err != nil {
				return err
			}
		}
	case 1:
		var ub = dbs.NewUpdateBuilder()
		ub.Table(this.table)
		ub.SET("status", status)
		ub.SET("updated_on", now)
		if status == nest.Disable {
			ub.Where("ctx = ? AND status != ? AND left_value >= ? AND right_value <= ?", node.Ctx, kDelete, node.LeftValue, node.RightValue)
		} else {
			ub.Where("ctx = ? AND id = ? AND status != ?", node.Ctx, node.Id, kDelete)
			ub.Limit(1)
		}
		if _, err := ub.Exec(this.db); err != nil {
			return err
		}
	case 0:
		if status == nest.Disable && node.IsLeaf() == false {
			return nest.ErrNotLeafNode
		}
		var ub = dbs.NewUpdateBuilder()
		ub.Table(this.table)
		ub.SET("status", status)
		ub.SET("updated_on", now)
		ub.Where("id = ?", id)
		ub.Where("status != ?", kDelete)
		if status == nest.Disable {
			ub.Where("right_value - left_value = 1")
		}
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

	// 计算参照节点
	var rNode *nest.Node

	switch position {
	case nest.Root: // 移动到顶级节点
		// 如果已经是顶级节点，则直接返回
		if node.Depth == 1 {
			return nil
		}
		// 如果是顶级节点，那么参照节点为顶级节点列表中的最后一个节点
		if rNode, err = this.getLastNode(node.Ctx, 0); err != nil {
			return err
		}
		if rNode != nil && rNode.Id == node.Id {
			return nil
		}
	case nest.Left: // 移动到节点左边
		if rId <= 0 {
			// 如果参照节点小于等于 0，则将节点向前移动一位，即向左移动一位
			if rNode, err = this.getPreviousNode(ctx, id); err != nil {
				return err
			}
			// 如果没有获取到上一个节点，则表示当前节点已经是最左边的节点
			if rNode == nil {
				return nil
			}
		} else {
			if rNode, err = this.getNodeWithId(ctx, rId); err != nil {
				return err
			}
		}
	case nest.Right: // 移动到节点右边
		if rId <= 0 {
			// 如果参照节点小于等于 0，则将节点向后移动一位，即向右移动一位
			if rNode, err = this.getNextNode(ctx, id); err != nil {
				return err
			}
			// 如果没有获取到下一个节点，则表示当前节点已经是最右边的节点
			if rNode == nil {
				return nil
			}
		} else {
			if rNode, err = this.getNodeWithId(ctx, rId); err != nil {
				return err
			}
		}
	case nest.First: // 移动到节点列表头部
		if rId <= 0 {
			// 如果参照节点小于等于 0，则将该节点移动到当前所在节点列表的头部
			// 获取其当前的父节点
			rNode, err = this.getParent(ctx, id, nest.Unknown)
			if err != nil {
				return err
			}

			// 如果父节点不存在，则表示当前节点为顶级节点，则找到顶级节点列表中的第一个节点
			if rNode == nil {
				rNode, err = this.getFirstNode(ctx, 0)
				if err != nil {
					return err
				}
				if rNode == nil {
					return nil
				}
				// 并且改变移动的位置类型为移动到指定节点的左边
				position = nest.Left
			}
		} else {
			if rNode, err = this.getNodeWithId(ctx, rId); err != nil {
				return err
			}
		}
	case nest.Last: // 移到到节点列表尾部
		if rId <= 0 {
			// 如果参照节点小于等于 0，则将该节点移动到当前所在节点列表的尾部
			// 获取其当前的父节点
			rNode, err = this.getParent(ctx, id, nest.Unknown)
			if err != nil {
				return err
			}

			// 如果父节点不存在，则表示当前节点为顶级节点，则找到顶级节点列表中的最后一个节点
			if rNode == nil {
				rNode, err = this.getLastNode(ctx, 0)
				if err != nil {
					return err
				}
				if rNode == nil {
					return nil
				}
				// 并且改变移动的位置类型为移动到指定节点的右边
				position = nest.Right
			}
		} else {
			if rNode, err = this.getNodeWithId(ctx, rId); err != nil {
				return err
			}
		}
	}

	if rNode == nil {
		return nest.ErrParentNotExist
	}

	if id == rNode.Id {
		return nil
	}

	// 判断被移动节点和目标参照节点是否属于同一 Ctx
	if rNode.Ctx != node.Ctx {
		return nest.ErrParentNotAllowed
	}

	// 循环连接问题，即 参照节点 是 被移动节点 的子节点
	if rNode.LeftValue > node.LeftValue && rNode.RightValue < node.RightValue {
		return nest.ErrParentNotAllowed
	}

	// 判断是否已经是子节点
	//if refer.LeftValue < node.LeftValue && refer.RightValue > node.RightValue && node.Depth - 1 == refer.Depth {
	//	tx.Rollback()
	//	return ErrParentNotAllowed
	//}

	// 查询出被移动节点的所有子节点
	//children, err := this.getNodeList(node.Id, 0, 0)
	children, err := this.getNodeList(node.Ctx, node.Id, nest.Unknown, 0, "", 0, 0, true)
	if err != nil {
		return err
	}

	var updateIdList []int64
	updateIdList = append(updateIdList, node.Id)
	for _, c := range children {
		updateIdList = append(updateIdList, c.Id)
	}

	if err = this.moveNodeWithPosition(position, node, rNode, updateIdList); err != nil {
		return err
	}
	return nil
}

func (this *nestRepository) moveNodeWithPosition(position nest.Position, node, rNode *nest.Node, updateIdList []int64) (err error) {
	var nodeLen = node.RightValue - node.LeftValue + 1
	var now = time.Now()

	// 把要移动的节点及其子节点从原树中删除掉
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value - ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND status != ? AND left_value > ?", node.Ctx, kDelete, node.RightValue)
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}
	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value - ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND status != ? AND right_value > ?", node.Ctx, kDelete, node.RightValue)
	if _, err = ubTreeRight.Exec(this.db); err != nil {
		return err
	}

	if rNode.LeftValue > node.RightValue {
		rNode.LeftValue -= nodeLen
	}
	if rNode.RightValue > node.RightValue {
		rNode.RightValue -= nodeLen
	}

	switch position {
	case nest.Root:
		return this.moveToRight(node, rNode, updateIdList, nodeLen)
	case nest.First:
		return this.moveToFirst(node, rNode, updateIdList, nodeLen)
	case nest.Last:
		return this.moveToLast(node, rNode, updateIdList, nodeLen)
	case nest.Left:
		return this.moveToLeft(node, rNode, updateIdList, nodeLen)
	case nest.Right:
		return this.moveToRight(node, rNode, updateIdList, nodeLen)
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
	ubTreeLeft.Where("ctx = ? AND status != ? AND left_value > ?", parent.Ctx, kDelete, parent.LeftValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND status != ? AND right_value > ?", parent.Ctx, kDelete, parent.LeftValue)
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
	ubTree.Where("status != ?", kDelete)
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
	ubTreeLeft.Where("status != ?", kDelete)
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND status != ? AND right_value >= ?", parent.Ctx, kDelete, parent.RightValue)
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
	ubTree.Where("status != ?", kDelete)
	if _, err = ubTree.Exec(this.db); err != nil {
		return err
	}

	return nil
}

func (this *nestRepository) moveToLeft(node, rNode *nest.Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND status != ? AND left_value >= ?", rNode.Ctx, kDelete, rNode.LeftValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND status != ? AND right_value >= ?", rNode.Ctx, kDelete, rNode.LeftValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeRight.Exec(this.db); err != nil {
		return err
	}

	//refer.LeftValue += nodeLen

	// 更新被移动节点的信息
	var diff = node.LeftValue - rNode.LeftValue
	var diffDepth = rNode.Depth - node.Depth
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	ubTree.Where("status != ?", kDelete)
	if _, err = ubTree.Exec(this.db); err != nil {
		return err
	}

	return nil
}

func (this *nestRepository) moveToRight(node, rNode *nest.Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("ctx = ? AND status != ? AND left_value > ?", rNode.Ctx, kDelete, rNode.RightValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeLeft.Exec(this.db); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("ctx = ? AND status != ? AND right_value > ?", rNode.Ctx, kDelete, rNode.RightValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = ubTreeRight.Exec(this.db); err != nil {
		return err
	}

	// 更新被移动节点的信息
	var diff = node.LeftValue - rNode.RightValue - 1
	var diffDepth = rNode.Depth - node.Depth
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	ubTree.Where("status != ?", kDelete)
	if _, err = ubTree.Exec(this.db); err != nil {
		return err
	}

	return nil
}
