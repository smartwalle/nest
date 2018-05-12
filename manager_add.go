package nest

import (
	"github.com/smartwalle/dbs"
	"time"
	"fmt"
)

const (
	K_ADD_POSITION_ROOT  = 0 // 顶级分类
	K_ADD_POSITION_FIRST = 1 // 列表头部 (子分类)
	K_ADD_POSITION_LAST  = 2 // 列表尾部 (子分类)
	K_ADD_POSITION_LEFT  = 3 // 左边 (兄弟分类)
	K_ADD_POSITION_RIGHT = 4 // 右边 (兄弟分类)
)

// AddCategory 添加分类
// cType: 分类类型（分类组）
// position:
// 		1、将新的分类添加到参照分类的子分类列表头部；
// 		2、将新的分类添加到参照分类的子分类列表尾部；
// 		3、将新的分类添加到参照分类的左边；
// 		4、将新的分类添加到参照分类的右边；
// referTo: 参照分类 id，如果值等于 0，则表示添加顶级分类
// name: 分类名
// status: 分类状态 1000、有效；2000、无效
// ext: 其它数据
func (this *Manager) AddCategory(cId int64, cType, position int, referTo int64, name string, status int, exts ...map[string]interface{}) (result int64, err error) {
	var sess = this.DB

	// 锁表
	var lock = dbs.WriteLock(this.Table)
	if _, err = lock.Exec(sess); err != nil {
		return 0, err
	}

	// 解锁
	defer func() {
		var unlock = dbs.UnlockTable()
		unlock.Exec(sess)
	}()

	var tx = dbs.MustTx(sess)

	// 查询出参照分类的信息
	var referCategory *BasicModel

	if position == K_ADD_POSITION_ROOT {
		// 如果是添加顶级分类，那么参照分类为 right value 最大的
		if err = this._getCategoryWithMaxRightValue(tx, cType, &referCategory); err != nil {
			return 0, err
		}

		// 如果参照分类为 nil，则创建一个虚拟的
		if referCategory == nil {
			referCategory = &BasicModel{}
			referCategory.Id = -1
			referCategory.Type = cType
			referCategory.LeftValue = 0
			referCategory.RightValue = 0
			referCategory.Depth = 1
		}
	} else {
		if err = this._getCategoryWithId(tx, referTo, &referCategory); err != nil {
			return 0, err
		}
		if referCategory == nil {
			tx.Rollback()
			return 0, ErrCategoryNotExists
		}
	}

	var ext map[string]interface{}
	if len(exts) > 0 {
		ext = exts[0]
	}

	if result, err = this.addCategoryWithPosition(tx, referCategory, cId, position, name, status, ext); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *Manager) addCategoryWithPosition(tx *dbs.Tx, refer *BasicModel, cId int64, position int, name string, status int, ext map[string]interface{}) (id int64, err error) {
	switch position {
	case K_ADD_POSITION_ROOT:
		return this.insertCategoryToRoot(tx, refer, cId, name, status, ext)
	case K_ADD_POSITION_FIRST:
		return this.insertCategoryToFirst(tx, refer, cId, name, status, ext)
	case K_ADD_POSITION_LAST:
		return this.insertCategoryToLast(tx, refer, cId, name, status, ext)
	case K_ADD_POSITION_LEFT:
		return this.insertCategoryToLeft(tx, refer, cId, name, status, ext)
	case K_ADD_POSITION_RIGHT:
		return this.insertCategoryToRight(tx, refer, cId, name, status, ext)
	}
	tx.Rollback()
	return 0, ErrUnknownPosition
}

func (this *Manager) insertCategoryToRoot(tx *dbs.Tx, refer *BasicModel, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
	var cType = refer.Type
	var leftValue = refer.RightValue + 1
	var rightValue = refer.RightValue + 2
	var depth = refer.Depth
	if id, err = this.insertCategory(tx, cId, cType, name, leftValue, rightValue, depth, status, ext); err != nil {
		return 0, err
	}
	return id, nil
}

func (this *Manager) insertCategoryToFirst(tx *dbs.Tx, refer *BasicModel, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
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

	if id, err = this.insertCategory(tx, cId, refer.Type, name, refer.LeftValue+1, refer.LeftValue+2, refer.Depth+1, status, ext); err != nil {
		return 0, err
	}
	return id, nil
}

func (this *Manager) insertCategoryToLast(tx *dbs.Tx, refer *BasicModel, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
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

	if id, err = this.insertCategory(tx, cId, refer.Type, name, refer.RightValue, refer.RightValue+1, refer.Depth+1, status, ext); err != nil {
		return 0, err
	}

	return id, nil
}

func (this *Manager) insertCategoryToLeft(tx *dbs.Tx, refer *BasicModel, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
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

	if id, err = this.insertCategory(tx, cId, refer.Type, name, refer.LeftValue, refer.LeftValue+1, refer.Depth, status, ext); err != nil {
		return 0, err
	}
	return id, nil
}

func (this *Manager) insertCategoryToRight(tx *dbs.Tx, refer *BasicModel, cId int64, name string, status int, ext map[string]interface{}) (id int64, err error) {
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

	if id, err = this.insertCategory(tx, cId, refer.Type, name, refer.RightValue+1, refer.RightValue+2, refer.Depth, status, ext); err != nil {
		return 0, err
	}
	return id, nil
}

func (this *Manager) insertCategory(tx *dbs.Tx, cId int64, cType int, name string, leftValue, rightValue, depth, status int, ext map[string]interface{}) (id int64, err error) {
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

	fmt.Println(ib.ToSQL())

	if result, err := tx.ExecInsertBuilder(ib); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}
