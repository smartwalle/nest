package nest

import (
	"github.com/smartwalle/dbs"
	"time"
	"sort"
)

const (
	K_MOVE_POSITION_ROOT  = 0 // 顶级节点
	K_MOVE_POSITION_FIRST = 1 // 列表头部 (子节点)
	K_MOVE_POSITION_LAST  = 2 // 列表尾部 (子节点)
	K_MOVE_POSITION_LEFT  = 3 // 左边 (兄弟节点)
	K_MOVE_POSITION_RIGHT = 4 // 右边 (兄弟节点)
)

func (this *Manager) updateNode(id int64, updateInfo map[string]interface{}) (err error) {
	if updateInfo == nil || len(updateInfo) == 0 {
		return nil
	}

	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.Table)

	delete(updateInfo, "id")
	delete(updateInfo, "type")
	delete(updateInfo, "left_value")
	delete(updateInfo, "right_value")
	delete(updateInfo, "depth")
	delete(updateInfo, "status")
	delete(updateInfo, "created_on")

	updateInfo["updated_on"] = time.Now()

	var keys = make([]string, 0, len(updateInfo) + 1)
	for key := range updateInfo {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ub.SET(key, updateInfo[key])
	}

	ub.Where("id = ?", id)
	ub.Limit(1)
	if _, err = ub.Exec(this.DB); err != nil {
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
func (this *Manager) updateNodeStatus(id int64, status, updateType int) (err error) {
	var tx = dbs.MustTx(this.DB)
	var node *Node
	if node, err = this._getNodeWithId(tx, id); err != nil {
		return err
	}

	if node == nil {
		tx.Rollback()
		return ErrNodeNotExist
	}

	if node.Status == status {
		tx.Rollback()
		return nil
	}

	var now = time.Now()

	switch updateType {
	case 2:
		if status == K_STATUS_DISABLE {
			var ub = dbs.NewUpdateBuilder()
			ub.Table(this.Table)
			ub.SET("status", status)
			ub.SET("right_value", dbs.SQL("left_value + 1"))
			ub.SET("updated_on", now)
			ub.Where("id = ?", id)
			ub.Limit(1)
			//if _, err := tx.ExecUpdateBuilder(ub); err != nil {
			//	return err
			//}
			if _, err := ub.ExecTx(tx); err != nil {
				return err
			}

			var ubChild = dbs.NewUpdateBuilder()
			ubChild.Table(this.Table)
			ubChild.SET("left_value", dbs.SQL("left_value + 1"))
			ubChild.SET("right_value", dbs.SQL("right_value + 1"))
			ubChild.SET("depth", dbs.SQL("depth-1"))
			ubChild.SET("updated_on", now)
			ubChild.Where("type = ? AND left_value > ? AND right_value < ?", node.Type, node.LeftValue, node.RightValue)
			//if _, err := tx.ExecUpdateBuilder(ubChild); err != nil {
			//	return err
			//}
			if _, err := ubChild.ExecTx(tx); err != nil {
				return err
			}
		}
	case 1:
		var ub = dbs.NewUpdateBuilder()
		ub.Table(this.Table)
		ub.SET("status", status)
		ub.SET("updated_on", now)
		ub.Where("type = ? AND left_value >= ? AND right_value <= ?", node.Type, node.LeftValue, node.RightValue)
		//if _, err := tx.ExecUpdateBuilder(ub); err != nil {
		//	return err
		//}
		if _, err := ub.ExecTx(tx); err != nil {
			return err
		}
	case 0:
		var ub = dbs.NewUpdateBuilder()
		ub.Table(this.Table)
		ub.SET("status", status)
		ub.SET("updated_on", now)
		ub.Where("id = ?", id)
		ub.Limit(1)
		//if _, err := tx.ExecUpdateBuilder(ub); err != nil {
		//	return err
		//}
		if _, err := ub.ExecTx(tx); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (this *Manager) moveNode(position int, id, rid int64) (err error) {
	if id == rid {
		return ErrParentNotAllowed
	}

	var tx = dbs.MustTx(this.DB)

	// 判断被移动的节点是否存在
	var node *Node
	if node, err = this._getNodeWithId(tx, id); err != nil {
		return err
	}
	if node == nil {
		tx.Rollback()
		return ErrNodeNotExist
	}

	// 判断参照节点是否存在
	var refer *Node
	if position == K_MOVE_POSITION_ROOT {
		// 如果是添加顶级节点，那么参照节点为 right value 最大的
		if refer, err = this._getNodeWithMaxRightValue(tx, node.Type); err != nil {
			return err
		}
		if refer != nil && refer.Id == node.Id {
			tx.Rollback()
			return nil
		}
	} else {
		if refer, err = this._getNodeWithId(tx, rid); err != nil {
			return err
		}
	}
	if refer == nil {
		tx.Rollback()
		return ErrParentNotExist
	}

	// 判断被移动节点和目标参照节点是否属于同一 type
	if refer.Type != node.Type {
		tx.Rollback()
		return ErrParentNotAllowed
	}

	// 循环连接问题，即 参照节点 是 被移动节点 的子节点
	if refer.LeftValue > node.LeftValue && refer.RightValue < node.RightValue {
		tx.Rollback()
		return ErrParentNotAllowed
	}

	// 判断是否已经是子节点
	//if refer.LeftValue < node.LeftValue && refer.RightValue > node.RightValue && node.Depth - 1 == refer.Depth {
	//	tx.Rollback()
	//	return ErrParentNotAllowed
	//}

	// 查询出被移动节点的所有子节点
	//children, err := this.getNodeList(node.Id, 0, 0)
	var children []*Node
	if err = this.getNodeList(node.Id, 0, 0, 0, "", 0, true, &children); err != nil {
		return err
	}

	var updateIdList []int64
	updateIdList = append(updateIdList, node.Id)
	for _, c := range children {
		updateIdList = append(updateIdList, c.Id)
	}

	if err = this.moveNodeWithPosition(tx, position, node, refer, updateIdList); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (this *Manager) moveNodeWithPosition(tx dbs.TX, position int, node, refer *Node, updateIdList []int64) (err error) {
	var nodeLen = node.RightValue - node.LeftValue + 1
	var now = time.Now()

	// 把要移动的节点及其子节点从原树中删除掉
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value - ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value > ?", node.Type, node.RightValue)
	//if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
	//	return err
	//}
	if _, err = ubTreeLeft.ExecTx(tx); err != nil {
		return err
	}
	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value - ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value > ?", node.Type, node.RightValue)
	//if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
	//	return err
	//}
	if _, err = ubTreeRight.ExecTx(tx); err != nil {
		return err
	}

	if refer.LeftValue > node.RightValue {
		refer.LeftValue -= nodeLen
	}
	if refer.RightValue > node.RightValue {
		refer.RightValue -= nodeLen
	}

	switch position {
	case K_MOVE_POSITION_ROOT:
		return this.moveToRight(tx, node, refer, updateIdList, nodeLen)
	case K_MOVE_POSITION_FIRST:
		return this.moveToFirst(tx, node, refer, updateIdList, nodeLen)
	case K_MOVE_POSITION_LAST:
		return this.moveToLast(tx, node, refer, updateIdList, nodeLen)
	case K_MOVE_POSITION_LEFT:
		return this.moveToLeft(tx, node, refer, updateIdList, nodeLen)
	case K_MOVE_POSITION_RIGHT:
		return this.moveToRight(tx, node, refer, updateIdList, nodeLen)
	}
	tx.Rollback()
	return ErrUnknownPosition
}

func (this *Manager) moveToFirst(tx dbs.TX, node, parent *Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value > ?", parent.Type, parent.LeftValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
	//	return err
	//}
	if _, err = ubTreeLeft.ExecTx(tx); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value > ?", parent.Type, parent.LeftValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
	//	return err
	//}
	if _, err = ubTreeRight.ExecTx(tx); err != nil {
		return err
	}

	parent.RightValue += nodeLen

	// 更新被移动节点的信息
	var diff = node.LeftValue - parent.LeftValue - 1
	var diffDepth = parent.Depth - node.Depth + 1
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.Table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTree); err != nil {
	//	return err
	//}
	if _, err = ubTree.ExecTx(tx); err != nil {
		return err
	}

	return nil
}

func (this *Manager) moveToLast(tx dbs.TX, node, parent *Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value > ?", parent.Type, parent.RightValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
	//	return err
	//}
	if _, err = ubTreeLeft.ExecTx(tx); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value >= ?", parent.Type, parent.RightValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
	//	return err
	//}
	if _, err = ubTreeRight.ExecTx(tx); err != nil {
		return err
	}

	parent.RightValue += nodeLen

	// 更新被移动节点的信息
	var diff = node.RightValue - parent.RightValue + 1
	var diffDepth = parent.Depth - node.Depth + 1
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.Table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTree); err != nil {
	//	return err
	//}
	if _, err = ubTree.ExecTx(tx); err != nil {
		return err
	}

	return nil
}

func (this *Manager) moveToLeft(tx dbs.TX, node, refer *Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value >= ?", refer.Type, refer.LeftValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
	//	return err
	//}
	if _, err = ubTreeLeft.ExecTx(tx); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value >= ?", refer.Type, refer.LeftValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
	//	return err
	//}
	if _, err = ubTreeRight.ExecTx(tx); err != nil {
		return err
	}

	//refer.LeftValue += nodeLen

	// 更新被移动节点的信息
	var diff = node.LeftValue - refer.LeftValue
	var diffDepth = refer.Depth - node.Depth
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.Table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTree); err != nil {
	//	return err
	//}
	if _, err = ubTree.ExecTx(tx); err != nil {
		return err
	}

	return nil
}

func (this *Manager) moveToRight(tx dbs.TX, node, refer *Node, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value > ?", refer.Type, refer.RightValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
	//	return err
	//}
	if _, err = ubTreeLeft.ExecTx(tx); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value > ?", refer.Type, refer.RightValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
	//	return err
	//}
	if _, err = ubTreeRight.ExecTx(tx); err != nil {
		return err
	}

	// 更新被移动节点的信息
	var diff = node.LeftValue - refer.RightValue - 1
	var diffDepth = refer.Depth - node.Depth
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.Table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	//if _, err = tx.ExecUpdateBuilder(ubTree); err != nil {
	//	return err
	//}
	if _, err = ubTree.ExecTx(tx); err != nil {
		return err
	}

	return nil
}
