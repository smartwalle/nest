package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
)

func (this *nestRepository) getMaxRightNode(ctx int64) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.table, "AS c")
	sb.Where("c.ctx = ?", ctx)
	sb.OrderBy("c.right_value DESC")
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *nestRepository) getPreviousNode(ctx, id int64) (result *nest.Node, err error) {
	// p.right_value = c.left_value - 1
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.name", "p.left_value", "p.right_value", "p.depth", "p.status", "p.created_on", "p.updated_on")
	sb.From(this.table, "AS c")
	sb.LeftJoin(this.table, "AS p ON p.right_value = c.left_value - 1")
	sb.Where("c.ctx = ?", ctx)
	sb.Where("c.id = ?", id)
	sb.Where("p.ctx = ?", ctx)
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *nestRepository) getNextNode(ctx, id int64) (result *nest.Node, err error) {
	// n.left_value = c.right_value + 1
	var sb = dbs.NewSelectBuilder()
	sb.Selects("n.id", "n.ctx", "n.name", "n.left_value", "n.right_value", "n.depth", "n.status", "n.created_on", "n.updated_on")
	sb.From(this.table, "AS c")
	sb.LeftJoin(this.table, "AS n ON n.left_value = c.right_value + 1")
	sb.Where("c.ctx = ?", ctx)
	sb.Where("c.id = ?", id)
	sb.Where("n.ctx = ?", ctx)
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *nestRepository) getNodeList(ctx, pId int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result []*nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.table, "AS c")
	if pId > 0 {
		if withParent {
			sb.LeftJoin(this.table, "AS pc ON pc.ctx = c.ctx AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
		} else {
			sb.LeftJoin(this.table, "AS pc ON pc.ctx = c.ctx AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
		}
		sb.Where("pc.id = ?", pId)
	}
	sb.Where("c.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	if depth > 0 {
		if pId > 0 {
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

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *nestRepository) getNodeWithId(ctx, id int64) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.table, "AS c")
	sb.Where("c.id = ?", id)
	sb.Where("c.ctx = ?", ctx)
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *nestRepository) getNodeWithName(ctx int, name string) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.table, "AS c")
	sb.Where("c.ctx = ? AND c.name = ?", ctx, name)
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

//func (this *nestRepository) getNodeList(ctx, pid int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result *nest.Node, err error) {
//	var sb = dbs.NewSelectBuilder()
//	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
//	sb.From(this.table, "AS c")
//	if pid > 0 {
//		if withParent {
//			sb.LeftJoin(this.table, "AS pc ON pc.ctx = c.ctx AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
//		} else {
//			sb.LeftJoin(this.table, "AS pc ON pc.ctx = c.ctx AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
//		}
//		sb.Where("pc.id = ?", pid)
//	}
//	sb.Where("c.ctx = ?", ctx)
//	if status > 0 {
//		sb.Where("c.status = ?", status)
//	}
//	if depth > 0 {
//		if pid > 0 {
//			sb.Where("c.depth - pc.depth <= ?", depth)
//		} else {
//			sb.Where("c.depth <= ?", depth)
//		}
//	}
//	if name != "" {
//		var keyword = "%" + name + "%"
//		sb.Where("c.name LIKE ?", keyword)
//	}
//	sb.OrderBy("c.ctx", "c.left_value")
//	if limit > 0 {
//		sb.Limit(limit)
//	}
//	if offset > 0 {
//		sb.Offset(offset)
//	}
//
//	if err = sb.Scan(this.db, &result); err != nil {
//		return nil, err
//	}
//	return result, nil
//}

func (this *nestRepository) getChildrenIdList(ctx, pId int64, status nest.Status, depth int, withParent bool) (result []int64, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id")
	sb.From(this.table, "AS c")
	if pId > 0 {
		if withParent {
			sb.LeftJoin(this.table, "AS pc ON pc.ctx = c.ctx AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
		} else {
			sb.LeftJoin(this.table, "AS pc ON pc.ctx = c.ctx AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
		}
		sb.Where("pc.id = ?", pId)
	}
	sb.Where("c.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	if depth > 0 {
		if pId > 0 {
			sb.Where("c.depth - pc.depth <= ?", depth)
		} else {
			sb.Where("c.depth <= ?", depth)
		}
	}
	sb.OrderBy("c.ctx", "c.left_value")

	var nodeList []*nest.Node
	if err = sb.Scan(this.db, &nodeList); err != nil {
		return nil, err
	}

	for _, c := range nodeList {
		result = append(result, c.Id)
	}
	return result, nil
}

func (this *nestRepository) getPathList(ctx, id int64, status nest.Status, withCurrentNode bool) (result []*nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.table, "AS sc")
	if withCurrentNode {
		sb.LeftJoin(this.table, "AS c ON c.ctx = sc.ctx AND c.left_value <= sc.left_value AND c.right_value >= sc.right_value")
	} else {
		sb.LeftJoin(this.table, "AS c ON c.ctx = sc.ctx AND c.left_value < sc.left_value AND c.right_value > sc.right_value")
	}
	sb.Where("sc.id = ?", id)
	sb.Where("sc.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	sb.OrderBy("c.left_value")
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *nestRepository) getParent(ctx, id int64, status nest.Status) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.table, "AS sc")
	sb.LeftJoin(this.table, "AS c ON c.ctx = sc.ctx AND c.left_value < sc.left_value AND c.right_value > sc.right_value")
	sb.Where("sc.id = ?", id)
	sb.Where("c.ctx = ?", ctx)
	if status > 0 {
		sb.Where("c.status = ?", status)
	}
	sb.Limit(1)
	sb.OrderBy("c.left_value DESC")

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}
