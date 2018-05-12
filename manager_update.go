package nest

import (
	"github.com/smartwalle/dbs"
	"time"
)

const (
	K_MOVE_POSITION_ROOT  = 0 // 顶级分类
	K_MOVE_POSITION_FIRST = 1 // 列表头部 (子分类)
	K_MOVE_POSITION_LAST  = 2 // 列表尾部 (子分类)
	K_MOVE_POSITION_LEFT  = 3 // 左边 (兄弟分类)
	K_MOVE_POSITION_RIGHT = 4 // 右边 (兄弟分类)
)

func (this *Manager) updateCategory(id int64, updateInfo map[string]interface{}) (err error) {
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
	ub.SetMap(updateInfo)

	ub.Where("id = ?", id)
	ub.Limit(1)
	if _, err = ub.Exec(this.DB); err != nil {
		return nil
	}
	return nil
}

// updateCategoryStatus 更新分类状态
// id: 被更新分类的 id
// status: 新的状态
// updateType:
// 		0、只更新当前分类的状态，子分类的状态不会受到影响，并且不会改变父子关系；
// 		1、子分类的状态会一起更新，不会改变父子关系；
// 		2、子分类的状态不会受到影响，并且所有子分类会向上移动一级（只针对把状态设置为 无效 的时候）；
func (this *Manager) updateCategoryStatus(id int64, status, updateType int) (err error) {
	var sess = this.DB

	// 锁表
	var lock = dbs.WriteLock(this.Table)
	if _, err = lock.Exec(sess); err != nil {
		return err
	}

	// 解锁
	defer func() {
		var unlock = dbs.UnlockTable()
		unlock.Exec(sess)
	}()

	var tx = dbs.MustTx(sess)
	var category *BaseModel
	if category, err = this._getCategoryWithId(tx, id); err != nil {
		return err
	}

	if category == nil {
		tx.Rollback()
		return ErrCategoryNotExists
	}

	if category.Status == status {
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
			if _, err := tx.ExecUpdateBuilder(ub); err != nil {
				return err
			}

			var ubChild = dbs.NewUpdateBuilder()
			ubChild.Table(this.Table)
			ubChild.SET("left_value", dbs.SQL("left_value + 1"))
			ubChild.SET("right_value", dbs.SQL("right_value + 1"))
			ubChild.SET("depth", dbs.SQL("depth-1"))
			ubChild.SET("updated_on", now)
			ubChild.Where("type = ? AND left_value > ? AND right_value < ?", category.Type, category.LeftValue, category.RightValue)
			if _, err := tx.ExecUpdateBuilder(ubChild); err != nil {
				return err
			}
		}
	case 1:
		var ub = dbs.NewUpdateBuilder()
		ub.Table(this.Table)
		ub.SET("status", status)
		ub.SET("updated_on", now)
		ub.Where("type = ? AND left_value >= ? AND right_value <= ?", category.Type, category.LeftValue, category.RightValue)
		if _, err := tx.ExecUpdateBuilder(ub); err != nil {
			return err
		}
	case 0:
		var ub = dbs.NewUpdateBuilder()
		ub.Table(this.Table)
		ub.SET("status", status)
		ub.SET("updated_on", now)
		ub.Where("id = ?", id)
		ub.Limit(1)
		if _, err := tx.ExecUpdateBuilder(ub); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (this *Manager) moveCategory(position int, id, rid int64) (err error) {
	if id == rid {
		return ErrParentNotAllowed
	}

	var sess = this.DB

	// 锁表
	var lock = dbs.WriteLock(this.Table)
	if _, err = lock.Exec(sess); err != nil {
		return err
	}

	// 解锁
	defer func() {
		var unlock = dbs.UnlockTable()
		unlock.Exec(sess)
	}()

	var tx = dbs.MustTx(sess)

	// 判断被移动的分类是否存在
	var category *BaseModel
	if category, err = this._getCategoryWithId(tx, id); err != nil {
		return err
	}
	if category == nil {
		tx.Rollback()
		return ErrCategoryNotExists
	}

	// 判断参照分类是否存在
	var refer *BaseModel
	if position == K_MOVE_POSITION_ROOT {
		// 如果是添加顶级分类，那么参照分类为 right value 最大的
		if refer, err = this._getCategoryWithMaxRightValue(tx, category.Type); err != nil {
			return err
		}
		if refer != nil && refer.Id == category.Id {
			tx.Rollback()
			return nil
		}
	} else {
		if refer, err = this._getCategoryWithId(tx, rid); err != nil {
			return err
		}
	}
	if refer == nil {
		tx.Rollback()
		return ErrParentCategoryNotExists
	}

	// 判断被移动分类和目标参照分类是否属于同一 type
	if refer.Type != category.Type {
		tx.Rollback()
		return ErrParentNotAllowed
	}

	// 循环连接问题，即 参照分类 是 被移动分类 的子分类
	if refer.LeftValue > category.LeftValue && refer.RightValue < category.RightValue {
		tx.Rollback()
		return ErrParentNotAllowed
	}

	// 判断是否已经是子分类
	//if refer.LeftValue < category.LeftValue && refer.RightValue > category.RightValue && category.Depth - 1 == refer.Depth {
	//	tx.Rollback()
	//	return ErrParentNotAllowed
	//}

	// 查询出被移动分类的所有子分类
	//children, err := this.getCategoryList(category.Id, 0, 0)
	var children []*BaseModel
	if err = this.getCategoryList(category.Id, 0, 0, 0, "", 0, true, &children); err != nil {
		return err
	}

	var updateIdList []int64
	updateIdList = append(updateIdList, category.Id)
	for _, c := range children {
		updateIdList = append(updateIdList, c.Id)
	}

	if err = this.moveCategoryWithPosition(tx, position, category, refer, updateIdList); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (this *Manager) moveCategoryWithPosition(tx *dbs.Tx, position int, category, refer *BaseModel, updateIdList []int64) (err error) {
	var nodeLen = category.RightValue - category.LeftValue + 1
	var now = time.Now()

	// 把要移动的节点及其子节点从原树中删除掉
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value - ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value > ?", category.Type, category.RightValue)
	if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
		return err
	}
	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value - ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value > ?", category.Type, category.RightValue)
	if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
		return err
	}

	if refer.LeftValue > category.RightValue {
		refer.LeftValue -= nodeLen
	}
	if refer.RightValue > category.RightValue {
		refer.RightValue -= nodeLen
	}

	switch position {
	case K_MOVE_POSITION_ROOT:
		return this.moveToRight(tx, category, refer, updateIdList, nodeLen)
	case K_MOVE_POSITION_FIRST:
		return this.moveToFirst(tx, category, refer, updateIdList, nodeLen)
	case K_MOVE_POSITION_LAST:
		return this.moveToLast(tx, category, refer, updateIdList, nodeLen)
	case K_MOVE_POSITION_LEFT:
		return this.moveToLeft(tx, category, refer, updateIdList, nodeLen)
	case K_MOVE_POSITION_RIGHT:
		return this.moveToRight(tx, category, refer, updateIdList, nodeLen)
	}
	tx.Rollback()
	return ErrUnknownPosition
}

func (this *Manager) moveToFirst(tx *dbs.Tx, category, parent *BaseModel, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value > ?", parent.Type, parent.LeftValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value > ?", parent.Type, parent.LeftValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
		return err
	}

	parent.RightValue += nodeLen

	// 更新被移动节点的信息
	var diff = category.LeftValue - parent.LeftValue - 1
	var diffDepth = parent.Depth - category.Depth + 1
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.Table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTree); err != nil {
		return err
	}

	return nil
}

func (this *Manager) moveToLast(tx *dbs.Tx, category, parent *BaseModel, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value > ?", parent.Type, parent.RightValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value >= ?", parent.Type, parent.RightValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
		return err
	}

	parent.RightValue += nodeLen

	// 更新被移动节点的信息
	var diff = category.RightValue - parent.RightValue + 1
	var diffDepth = parent.Depth - category.Depth + 1
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.Table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTree); err != nil {
		return err
	}

	return nil
}

func (this *Manager) moveToLeft(tx *dbs.Tx, category, refer *BaseModel, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value >= ?", refer.Type, refer.LeftValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value >= ?", refer.Type, refer.LeftValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
		return err
	}

	//refer.LeftValue += nodeLen

	// 更新被移动节点的信息
	var diff = category.LeftValue - refer.LeftValue
	var diffDepth = refer.Depth - category.Depth
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.Table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTree); err != nil {
		return err
	}

	return nil
}

func (this *Manager) moveToRight(tx *dbs.Tx, category, refer *BaseModel, updateIdList []int64, nodeLen int) (err error) {
	var now = time.Now()

	// 移出空间用于存放被移动的节点及其子节点
	var ubTreeLeft = dbs.NewUpdateBuilder()
	ubTreeLeft.Table(this.Table)
	ubTreeLeft.SET("left_value", dbs.SQL("left_value + ?", nodeLen))
	ubTreeLeft.SET("updated_on", now)
	ubTreeLeft.Where("type = ? AND left_value > ?", refer.Type, refer.RightValue)
	ubTreeLeft.Where(dbs.NotIn("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTreeLeft); err != nil {
		return err
	}

	var ubTreeRight = dbs.NewUpdateBuilder()
	ubTreeRight.Table(this.Table)
	ubTreeRight.SET("right_value", dbs.SQL("right_value + ?", nodeLen))
	ubTreeRight.SET("updated_on", now)
	ubTreeRight.Where("type = ? AND right_value > ?", refer.Type, refer.RightValue)
	ubTreeRight.Where(dbs.NotIn("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTreeRight); err != nil {
		return err
	}

	// 更新被移动节点的信息
	var diff = category.LeftValue - refer.RightValue - 1
	var diffDepth = refer.Depth - category.Depth
	var ubTree = dbs.NewUpdateBuilder()
	ubTree.Table(this.Table)
	ubTree.SET("left_value", dbs.SQL("left_value - ?", diff))
	ubTree.SET("right_value", dbs.SQL("right_value - ?", diff))
	ubTree.SET("depth", dbs.SQL("depth + ?", diffDepth))
	ubTree.SET("updated_on", now)
	ubTree.Where(dbs.IN("id", updateIdList))
	if _, err = tx.ExecUpdateBuilder(ubTree); err != nil {
		return err
	}

	return nil
}
