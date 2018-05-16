package nest

import (
	"github.com/smartwalle/dbs"
	"time"
)

const (
	K_ADD_POSITION_ROOT  = 0 // 顶级节点
	K_ADD_POSITION_FIRST = 1 // 列表头部 (子节点)
	K_ADD_POSITION_LAST  = 2 // 列表尾部 (子节点)
	K_ADD_POSITION_LEFT  = 3 // 左边 (兄弟节点)
	K_ADD_POSITION_RIGHT = 4 // 右边 (兄弟节点)
)

// addNode 添加节点
// cType: 节点类型（节点组）
// position:
// 		1、将新的节点添加到参照节点的子节点列表头部；
// 		2、将新的节点添加到参照节点的子节点列表尾部；
// 		3、将新的节点添加到参照节点的左边；
// 		4、将新的节点添加到参照节点的右边；
// referTo: 参照节点 id，如果值等于 0，则表示添加顶级节点
// name: 节点名
// status: 节点状态 1000、有效；2000、无效
// ext: 其它数据
func (this *Manager) addNode(cId int64, cType, position int, referTo int64, name string, status int, exts ...map[string]interface{}) (result int64, err error) {
	// 锁表
	this.lockTable()
	// 解锁
	defer func() {
		this.unlockTable()
	}()

	var tx = dbs.MustTx(this.DB)

	// 查询出参照节点的信息
	var referNode *Node

	if position == K_ADD_POSITION_ROOT {
		// 如果是添加顶级节点，那么参照节点为 right value 最大的
		if referNode, err = this._getNodeWithMaxRightValue(tx, cType); err != nil {
			return 0, err
		}

		// 如果参照节点为 nil，则创建一个虚拟的
		if referNode == nil {
			referNode = &Node{}
			referNode.Id = -1
			referNode.Type = cType
			referNode.LeftValue = 0
			referNode.RightValue = 0
			referNode.Depth = 1
		}
	} else {
		if referNode, err = this._getNodeWithId(tx, referTo); err != nil {
			return 0, err
		}
		if referNode == nil {
			tx.Rollback()
			return 0, ErrNodeNotExist
		}
	}

	var ext map[string]interface{}
	if len(exts) > 0 {
		ext = exts[0]
	}

	if result, err = this.addNodeWithPosition(tx, referNode, cId, position, name, status, ext); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *Manager) addNodeWithPosition(tx *dbs.Tx, refer *Node, cId int64, position int, name string, status int, ext map[string]interface{}) (id int64, err error) {
	switch position {
	case K_ADD_POSITION_ROOT:
		return this.insertNodeToRoot(tx, refer, cId, name, status, ext)
	case K_ADD_POSITION_FIRST:
		return this.insertNodeToFirst(tx, refer, cId, name, status, ext)
	case K_ADD_POSITION_LAST:
		return this.insertNodeToLast(tx, refer, cId, name, status, ext)
	case K_ADD_POSITION_LEFT:
		return this.insertNodeToLeft(tx, refer, cId, name, status, ext)
	case K_ADD_POSITION_RIGHT:
		return this.insertNodeToRight(tx, refer, cId, name, status, ext)
	}
	tx.Rollback()
	return 0, ErrUnknownPosition
}

func (this *Manager) insertNodeToRoot(tx *dbs.Tx, refer *Node, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
	var cType = refer.Type
	var leftValue = refer.RightValue + 1
	var rightValue = refer.RightValue + 2
	var depth = refer.Depth
	if id, err = this.insertNode(tx, cId, cType, name, leftValue, rightValue, depth, status, ext); err != nil {
		return 0, err
	}
	return id, nil
}

func (this *Manager) insertNodeToFirst(tx *dbs.Tx, refer *Node, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.Table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("type = ? AND left_value > ?", refer.Type, refer.LeftValue)
	if _, err = tx.ExecUpdateBuilder(ubLeft); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.Table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("type = ? AND right_value > ?", refer.Type, refer.LeftValue)
	if _, err = tx.ExecUpdateBuilder(ubRight); err != nil {
		return 0, err
	}

	if id, err = this.insertNode(tx, cId, refer.Type, name, refer.LeftValue+1, refer.LeftValue+2, refer.Depth+1, status, ext); err != nil {
		return 0, err
	}
	return id, nil
}

func (this *Manager) insertNodeToLast(tx *dbs.Tx, refer *Node, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.Table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("type = ? AND left_value > ?", refer.Type, refer.RightValue)
	if _, err = tx.ExecUpdateBuilder(ubLeft); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.Table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("type = ? AND right_value >= ?", refer.Type, refer.RightValue)
	if _, err = tx.ExecUpdateBuilder(ubRight); err != nil {
		return 0, err
	}

	if id, err = this.insertNode(tx, cId, refer.Type, name, refer.RightValue, refer.RightValue+1, refer.Depth+1, status, ext); err != nil {
		return 0, err
	}

	return id, nil
}

func (this *Manager) insertNodeToLeft(tx *dbs.Tx, refer *Node, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.Table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("type = ? AND left_value >= ?", refer.Type, refer.LeftValue)
	if _, err = tx.ExecUpdateBuilder(ubLeft); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.Table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("type = ? AND right_value >= ?", refer.Type, refer.LeftValue)
	if _, err = tx.ExecUpdateBuilder(ubRight); err != nil {
		return 0, err
	}

	if id, err = this.insertNode(tx, cId, refer.Type, name, refer.LeftValue, refer.LeftValue+1, refer.Depth, status, ext); err != nil {
		return 0, err
	}
	return id, nil
}

func (this *Manager) insertNodeToRight(tx *dbs.Tx, refer *Node, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.Table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("type = ? AND left_value > ?", refer.Type, refer.RightValue)
	if _, err = tx.ExecUpdateBuilder(ubLeft); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.Table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("type = ? AND right_value > ?", refer.Type, refer.RightValue)
	if _, err = tx.ExecUpdateBuilder(ubRight); err != nil {
		return 0, err
	}

	if id, err = this.insertNode(tx, cId, refer.Type, name, refer.RightValue+1, refer.RightValue+2, refer.Depth, status, ext); err != nil {
		return 0, err
	}
	return id, nil
}

func (this *Manager) insertNode(tx *dbs.Tx, cId int64, cType int, name string, leftValue, rightValue, depth, status int, ext map[string]interface{}) (id int64, err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.Table)

	if ext == nil {
		ext = make(map[string]interface{})
	}

	ext["id"] = id
	ext["type"] = cType
	ext["name"] = name
	ext["left_value"] = leftValue
	ext["right_value"] = rightValue
	ext["depth"] = depth
	ext["status"] = status
	ext["created_on"] = now
	ext["updated_on"] = now

	for key, value := range ext {
		ib.SET(key, value)
	}

	if result, err := tx.ExecInsertBuilder(ib); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}
