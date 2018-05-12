package nest

import (
	"github.com/smartwalle/dbs"
)

func (this *Manager) GetCategory(id int64, result interface{}) (err error) {
	var tx = dbs.MustTx(this.db)

	if err = this.getCategoryWithId(tx, id, result); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (this *Manager) getCategoryWithId(tx *dbs.Tx, id int64, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.table, "AS c")
	sb.Where("c.id = ?", id)
	sb.Limit(1)
	if err = tx.ExecSelectBuilder(sb, result); err != nil {
		return err
	}
	return nil
}

func (this *Manager) GetCategoryWithName(cType int, name string, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.table, "AS c")
	sb.Where("c.type = ? AND c.name = ?", cType, name)
	sb.Limit(1)
	if err = sb.Scan(this.db, result); err != nil {
		return err
	}
	return nil
}

func (this *Manager) getCategoryWithMaxRightValue(tx *dbs.Tx, cType int, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.table, "AS c")
	if cType > 0 {
		sb.Where("c.type = ?", cType)
	}
	sb.OrderBy("c.right_value DESC")
	sb.Limit(1)
	if err = tx.ExecSelectBuilder(sb, result); err != nil {
		return err
	}
	return nil
}

func (this *Manager) GetCategoryList(parentId int64, cType, status, depth int, name string, limit uint64, includeParent bool, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.table, "AS c")
	if parentId > 0 {
		if includeParent {
			sb.LeftJoin(this.table, "AS pc ON pc.type = c.type AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
		} else {
			sb.LeftJoin(this.table, "AS pc ON pc.type = c.type AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
		}
		sb.Where("pc.id = ?", parentId)
	} else {
		if cType > 0 {
			sb.Where("c.type = ?", cType)
		}
	}
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	if depth > 0 {
		if parentId > 0 {
			sb.Where("c.depth - pc.depth <= ?", depth)
		} else {
			sb.Where("c.depth <= ?", depth)
		}
	}
	if name != "" {
		var keyword = "%" + name + "%"
		sb.Where("c.name LIKE ?", keyword)
	}
	sb.OrderBy("c.type", "c.left_value")
	if limit > 0 {
		sb.Limit(limit)
	}

	if err = sb.Scan(this.db, result); err != nil {
		return err
	}
	return nil
}

func (this *Manager) GetIdList(parentId int64, status, depth int, includeParent bool) (result []int64, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id")
	sb.From(this.table, "AS c")
	if parentId > 0 {
		if includeParent {
			sb.LeftJoin(this.table, "AS pc ON pc.type = c.type AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
		} else {
			sb.LeftJoin(this.table, "AS pc ON pc.type = c.type AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
		}
		sb.Where("pc.id = ?", parentId)
	}
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	if depth > 0 {
		if parentId > 0 {
			sb.Where("c.depth - pc.depth <= ?", depth)
		} else {
			sb.Where("c.depth <= ?", depth)
		}
	}
	sb.OrderBy("c.type", "c.left_value")

	var categoryList []*BasicModel
	if err = sb.Scan(this.db, &categoryList); err != nil {
		return nil, err
	}

	for _, c := range categoryList {
		result = append(result, c.Id)
	}
	return result, nil
}

func (this *Manager) GetPathList(id int64, status int, includeLastNode bool, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.table, "AS sc")
	if includeLastNode {
		sb.LeftJoin(this.table, "AS c ON c.type = sc.type AND c.left_value <= sc.left_value AND c.right_value >= sc.right_value")
	} else {
		sb.LeftJoin(this.table, "AS c ON c.type = sc.type AND c.left_value < sc.left_value AND c.right_value > sc.right_value")
	}
	sb.Where("sc.id = ?", id)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	sb.OrderBy("c.left_value")
	if err = sb.Scan(this.db, result); err != nil {
		return err
	}
	return nil
}
