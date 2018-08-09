package nest

import (
	"github.com/smartwalle/dbs"
)

func (this *Manager) _getNodeWithId(tx dbs.TX, id int64) (result *Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS c")
	sb.Where("c.id = ?", id)
	sb.Limit(1)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Manager) _getNodeWithMaxRightValue(tx dbs.TX, ctx int64) (result *Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS c")
	sb.Where("c.ctx = ?", ctx)
	sb.OrderBy("c.right_value DESC")
	sb.Limit(1)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Manager) _getNodeList(ctx, pid int64, status, depth int, name string, limit, offset int64, includeParent bool) (result []*Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS c")
	if pid > 0 {
		if includeParent {
			sb.LeftJoin(this.Table, "AS pc ON pc.ctx = c.ctx AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
		} else {
			sb.LeftJoin(this.Table, "AS pc ON pc.ctx = c.ctx AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
		}
		sb.Where("pc.id = ?", pid)
	}
	sb.Where("c.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	if depth > 0 {
		if pid > 0 {
			sb.Where("c.depth - pc.depth <= ?", depth)
		} else {
			sb.Where("c.depth <= ?", depth)
		}
	}
	if name != "" {
		var keyword = "%" + name + "%"
		sb.Where("c.name LIKE ?", keyword)
	}
	sb.OrderBy("c.ctx", "c.left_value")
	if limit > 0 {
		sb.Limit(limit)
	}
	if offset > 0 {
		sb.Offset(offset)
	}

	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Manager) getNode(ctx, id int64, result interface{}) (err error) {
	var tx = dbs.MustTx(this.DB)

	if err = this.getNodeWithId(tx, ctx, id, result); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (this *Manager) getNodeWithId(tx dbs.TX, ctx, id int64, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.Table, "AS c")
	sb.Where("c.id = ?", id)
	sb.Where("c.ctx = ?", ctx)
	sb.Limit(1)
	if err = sb.Scan(tx, result); err != nil {
		return err
	}
	return nil
}

func (this *Manager) getNodeWithName(ctx int, name string, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.Table, "AS c")
	sb.Where("c.ctx = ? AND c.name = ?", ctx, name)
	sb.Limit(1)
	if err = sb.Scan(this.DB, result); err != nil {
		return err
	}
	return nil
}

func (this *Manager) getNodeList(ctx, pid int64, status, depth int, name string, limit, offset int64, includeParent bool, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.Table, "AS c")
	if pid > 0 {
		if includeParent {
			sb.LeftJoin(this.Table, "AS pc ON pc.ctx = c.ctx AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
		} else {
			sb.LeftJoin(this.Table, "AS pc ON pc.ctx = c.ctx AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
		}
		sb.Where("pc.id = ?", pid)
	}
	sb.Where("c.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	if depth > 0 {
		if pid > 0 {
			sb.Where("c.depth - pc.depth <= ?", depth)
		} else {
			sb.Where("c.depth <= ?", depth)
		}
	}
	if name != "" {
		var keyword = "%" + name + "%"
		sb.Where("c.name LIKE ?", keyword)
	}
	sb.OrderBy("c.ctx", "c.left_value")
	if limit > 0 {
		sb.Limit(limit)
	}
	if offset > 0 {
		sb.Offset(offset)
	}

	if err = sb.Scan(this.DB, result); err != nil {
		return err
	}
	return nil
}

func (this *Manager) getIdList(ctx, pid int64, status, depth int, includeParent bool) (result []int64, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id")
	sb.From(this.Table, "AS c")
	if pid > 0 {
		if includeParent {
			sb.LeftJoin(this.Table, "AS pc ON pc.ctx = c.ctx AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
		} else {
			sb.LeftJoin(this.Table, "AS pc ON pc.ctx = c.ctx AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
		}
		sb.Where("pc.id = ?", pid)
	}
	sb.Where("c.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	if depth > 0 {
		if pid > 0 {
			sb.Where("c.depth - pc.depth <= ?", depth)
		} else {
			sb.Where("c.depth <= ?", depth)
		}
	}
	sb.OrderBy("c.ctx", "c.left_value")

	var nodeList []*Node
	if err = sb.Scan(this.DB, &nodeList); err != nil {
		return nil, err
	}

	for _, c := range nodeList {
		result = append(result, c.Id)
	}
	return result, nil
}

func (this *Manager) getPathList(ctx, id int64, status int, includeLastNode bool, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.Table, "AS sc")
	if includeLastNode {
		sb.LeftJoin(this.Table, "AS c ON c.ctx = sc.ctx AND c.left_value <= sc.left_value AND c.right_value >= sc.right_value")
	} else {
		sb.LeftJoin(this.Table, "AS c ON c.ctx = sc.ctx AND c.left_value < sc.left_value AND c.right_value > sc.right_value")
	}
	sb.Where("sc.id = ?", id)
	sb.Where("sc.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	sb.OrderBy("c.left_value")
	if err = sb.Scan(this.DB, result); err != nil {
		return err
	}
	return nil
}

func (this *Manager) getParent(ctx, id int64, status int, result interface{}) (err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects(this.SelectFields...)
	sb.From(this.Table, "AS sc")
	sb.LeftJoin(this.Table, "AS c ON c.ctx = sc.ctx AND c.left_value < sc.left_value AND c.right_value > sc.right_value")
	sb.Where("sc.id = ?", id)
	sb.Where("c.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	sb.Limit(1)
	sb.OrderBy("c.left_value DESC")

	if err = sb.Scan(this.DB, result); err != nil {
		return err
	}
	return nil
}
