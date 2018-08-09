package mysql

import (
	"github.com/smartwalle/dbs"
	"sort"
	"time"
	"github.com/smartwalle/nest"
)

// addNode 添加节点
// cId: 指定节点的 id，如果值小于等于0，则表示自增；
// ctx: 节点组标识；
// position:
// 		1、将新的节点添加到参照节点的子节点列表头部；
// 		2、将新的节点添加到参照节点的子节点列表尾部；
// 		3、将新的节点添加到参照节点的左边；
// 		4、将新的节点添加到参照节点的右边；
// referTo: 参照节点 id，如果值等于 0，则表示添加顶级节点；
// name: 节点名
// status: 节点状态 1000、有效；2000、无效；
// ext: 其它数据；
func (this *nestRepository) addNode(id, ctx int64, position int, rid int64, name string, status int, exts ...map[string]interface{}) (result int64, err error) {
	var tx = dbs.MustTx(this.db)

	// 查询出参照节点的信息
	var referNode *nest.Node

	if position == nest.K_ADD_POSITION_ROOT {
		// 如果是添加顶级节点，那么参照节点为 right value 最大的
		if referNode, err = this._getNodeWithMaxRightValue(tx, ctx); err != nil {
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
		if referNode, err = this._getNodeWithId(tx, rid); err != nil {
			return 0, err
		}
		if referNode == nil {
			tx.Rollback()
			return 0, nest.ErrNodeNotExist
		}
	}

	var ext map[string]interface{}
	if len(exts) > 0 {
		ext = exts[0]
	}

	if result, err = this.addNodeWithPosition(tx, referNode, id, position, name, status, ext); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) addNodeWithPosition(tx dbs.TX, refer *nest.Node, id int64, position int, name string, status int, ext map[string]interface{}) (result int64, err error) {
	switch position {
	case nest.K_ADD_POSITION_ROOT:
		return this.insertNodeToRoot(tx, refer, id, name, status, ext)
	case nest.K_ADD_POSITION_FIRST:
		return this.insertNodeToFirst(tx, refer, id, name, status, ext)
	case nest.K_ADD_POSITION_LAST:
		return this.insertNodeToLast(tx, refer, id, name, status, ext)
	case nest.K_ADD_POSITION_LEFT:
		return this.insertNodeToLeft(tx, refer, id, name, status, ext)
	case nest.K_ADD_POSITION_RIGHT:
		return this.insertNodeToRight(tx, refer, id, name, status, ext)
	}
	tx.Rollback()
	return 0, nest.ErrUnknownPosition
}

func (this *nestRepository) insertNodeToRoot(tx dbs.TX, refer *nest.Node, id int64, name string, status int, ext map[string]interface{}) (result int64, err error) {
	var ctx = refer.Ctx
	var leftValue = refer.RightValue + 1
	var rightValue = refer.RightValue + 2
	var depth = refer.Depth
	if result, err = this.insertNode(tx, id, ctx, name, leftValue, rightValue, depth, status, ext); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) insertNodeToFirst(tx dbs.TX, refer *nest.Node, id int64, name string, status int, ext map[string]interface{}) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", refer.Ctx, refer.LeftValue)
	if _, err = ubLeft.Exec(tx); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value > ?", refer.Ctx, refer.LeftValue)
	if _, err = ubRight.Exec(tx); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(tx, id, refer.Ctx, name, refer.LeftValue+1, refer.LeftValue+2, refer.Depth+1, status, ext); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) insertNodeToLast(tx dbs.TX, refer *nest.Node, id int64, name string, status int, ext map[string]interface{}) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", refer.Ctx, refer.RightValue)
	if _, err = ubLeft.Exec(tx); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value >= ?", refer.Ctx, refer.RightValue)
	if _, err = ubRight.Exec(tx); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(tx, id, refer.Ctx, name, refer.RightValue, refer.RightValue+1, refer.Depth+1, status, ext); err != nil {
		return 0, err
	}

	return result, nil
}

func (this *nestRepository) insertNodeToLeft(tx dbs.TX, refer *nest.Node, id int64, name string, status int, ext map[string]interface{}) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value >= ?", refer.Ctx, refer.LeftValue)
	if _, err = ubLeft.Exec(tx); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value >= ?", refer.Ctx, refer.LeftValue)
	if _, err = ubRight.Exec(tx); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(tx, id, refer.Ctx, name, refer.LeftValue, refer.LeftValue+1, refer.Depth, status, ext); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) insertNodeToRight(tx dbs.TX, refer *nest.Node, id int64, name string, status int, ext map[string]interface{}) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.table)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", refer.Ctx, refer.RightValue)
	if _, err = ubLeft.Exec(tx); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.table)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value > ?", refer.Ctx, refer.RightValue)
	if _, err = ubRight.Exec(tx); err != nil {
		return 0, err
	}

	if result, err = this.insertNode(tx, id, refer.Ctx, name, refer.RightValue+1, refer.RightValue+2, refer.Depth, status, ext); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *nestRepository) insertNode(tx dbs.TX, id, ctx int64, name string, leftValue, rightValue, depth, status int, ext map[string]interface{}) (result int64, err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.table)

	if ext == nil {
		ext = make(map[string]interface{})
	}

	ext["id"] = id
	ext["ctx"] = ctx
	ext["name"] = name
	ext["left_value"] = leftValue
	ext["right_value"] = rightValue
	ext["depth"] = depth
	ext["status"] = status
	ext["created_on"] = now
	ext["updated_on"] = now

	var keys = make([]string, 0, len(ext))
	for key := range ext {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ib.SET(key, ext[key])
	}

	if sResult, err := ib.Exec(tx); err != nil {
		return 0, err
	} else {
		result, _ = sResult.LastInsertId()
	}
	return result, err
}
