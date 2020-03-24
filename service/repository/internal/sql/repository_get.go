package sql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/nest"
)

// getLastNode 获取节点列表中的最后一个节点
func (this *Repository) getLastNode(ctx, pId int64) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS c")
	if pId > 0 {
		sb.LeftJoin(this.Table, "AS p ON p.status != ? AND p.left_value < c.left_value AND p.right_value > c.right_value", Delete)
		sb.Where("p.id = ?", pId)
		sb.Where("p.ctx = ?", ctx)
		sb.Where("c.right_value = p.right_value - 1")
	} else {
		sb.OrderBy("c.right_value DESC")
	}
	sb.Where("c.ctx = ?", ctx)
	sb.Where("c.status != ?", Delete)
	sb.Limit(1)
	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// getFirstNode 获取节点列表中的第一个节点
func (this *Repository) getFirstNode(ctx, pId int64) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS c")
	if pId > 0 {
		sb.LeftJoin(this.Table, "AS p ON p.status != ? AND p.left_value < c.left_value AND p.right_value > c.right_value", Delete)
		sb.Where("p.id = ?", pId)
		sb.Where("p.ctx = ?", ctx)
		sb.Where("c.left_value = p.left_value + 1")
	} else {
		sb.OrderBy("c.left_value ASC")
	}
	sb.Where("c.ctx = ?", ctx)
	sb.Where("c.status != ?", Delete)
	sb.Limit(1)
	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Repository) getPreviousNode(ctx, id int64) (result *nest.Node, err error) {
	// p.right_value = c.left_value - 1
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("p.id", "p.ctx", "p.name", "p.left_value", "p.right_value", "p.depth", "p.status", "p.created_on", "p.updated_on")
	sb.From(this.Table, "AS c")
	sb.LeftJoin(this.Table, "AS p ON p.status != ? AND p.right_value = c.left_value - 1", Delete)
	sb.Where("c.ctx = ?", ctx)
	sb.Where("c.id = ?", id)
	sb.Where("c.status != ?", Delete)
	sb.Where("p.ctx = ?", ctx)
	sb.Limit(1)
	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Repository) getNextNode(ctx, id int64) (result *nest.Node, err error) {
	// n.left_value = c.right_value + 1
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("n.id", "n.ctx", "n.name", "n.left_value", "n.right_value", "n.depth", "n.status", "n.created_on", "n.updated_on")
	sb.From(this.Table, "AS c")
	sb.LeftJoin(this.Table, "AS n ON n.status != ? AND n.left_value = c.right_value + 1", Delete)
	sb.Where("c.ctx = ?", ctx)
	sb.Where("c.id = ?", id)
	sb.Where("c.status != ?", Delete)
	sb.Where("n.ctx = ?", ctx)
	sb.Limit(1)
	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Repository) getNodes(ctx, pId int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result []*nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS c")
	if pId > 0 {
		if withParent {
			sb.LeftJoin(this.Table, "AS pc ON pc.status != ? AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value", Delete)
		} else {
			sb.LeftJoin(this.Table, "AS pc ON pc.status != ? AND pc.left_value < c.left_value AND pc.right_value > c.right_value", Delete)
		}
		sb.Where("pc.id = ?", pId)
		sb.Where("pc.ctx = ?", ctx)
	}
	sb.Where("c.ctx = ?", ctx)
	if status != nest.All && status != Delete {
		sb.Where("c.status = ?", status)
	} else {
		sb.Where("c.status != ?", Delete)
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

	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Repository) getNodeIds(ctx, pId int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result []int64, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("c.id")
	sb.From(this.Table, "AS c")
	if pId > 0 {
		if withParent {
			sb.LeftJoin(this.Table, "AS pc ON pc.status != ? AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value", Delete)
		} else {
			sb.LeftJoin(this.Table, "AS pc ON pc.status != ? AND pc.left_value < c.left_value AND pc.right_value > c.right_value", Delete)
		}
		sb.Where("pc.id = ?", pId)
		sb.Where("pc.ctx = ?", ctx)
	}
	sb.Where("c.ctx = ?", ctx)
	if status != nest.All && status != Delete {
		sb.Where("c.status = ?", status)
	} else {
		sb.Where("c.status != ?", Delete)
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

	var nodeList []*nest.Node
	if err = sb.Scan(this.DB, &nodeList); err != nil {
		return nil, err
	}
	if limit > 0 {
		sb.Limit(limit)
	}
	if offset > 0 {
		sb.Offset(offset)
	}

	for _, c := range nodeList {
		result = append(result, c.Id)
	}
	return result, nil
}

func (this *Repository) getNodeWithId(ctx, id int64) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS c")
	sb.Where("c.id = ?", id)
	sb.Where("c.ctx = ?", ctx)
	sb.Where("c.status != ?", Delete)
	sb.Limit(1)
	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Repository) getNodeWithName(ctx int, name string) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS c")
	sb.Where("c.ctx = ? AND c.status != ? AND c.name = ?", ctx, Delete, name)
	sb.Limit(1)
	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

//func (this *Repository) getNodes(ctx, pid int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result *nest.Node, err error) {
//	var sb = dbs.NewSelectBuilder().UseDialect(this.Dialect)
//	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
//	sb.From(this.Table, "AS c")
//	if pid > 0 {
//		if withParent {
//			sb.LeftJoin(this.Table, "AS pc ON pc.ctx = c.ctx AND pc.left_value <= c.left_value AND pc.right_value >= c.right_value")
//		} else {
//			sb.LeftJoin(this.Table, "AS pc ON pc.ctx = c.ctx AND pc.left_value < c.left_value AND pc.right_value > c.right_value")
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
//	if err = sb.Scan(this.DB, &result); err != nil {
//		return nil, err
//	}
//	return result, nil
//}

func (this *Repository) getParentNodes(ctx, id int64, status nest.Status, withCurrentNode bool) (result []*nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS sc")
	if withCurrentNode {
		sb.LeftJoin(this.Table, "AS c ON c.ctx = sc.ctx AND c.left_value <= sc.left_value AND c.right_value >= sc.right_value")
	} else {
		sb.LeftJoin(this.Table, "AS c ON c.ctx = sc.ctx AND c.left_value < sc.left_value AND c.right_value > sc.right_value")
	}
	sb.Where("sc.id = ?", id)
	sb.Where("sc.ctx = ?", ctx)
	sb.Where("sc.status != ?", Delete)
	if status != nest.All && status != Delete {
		sb.Where("c.status = ?", status)
	} else {
		sb.Where("c.status != ?", Delete)
	}
	sb.OrderBy("c.left_value")
	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Repository) getParent(ctx, id int64, status nest.Status) (result *nest.Node, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.Dialect)
	sb.Selects("c.id", "c.ctx", "c.name", "c.left_value", "c.right_value", "c.depth", "c.status", "c.created_on", "c.updated_on")
	sb.From(this.Table, "AS sc")
	sb.LeftJoin(this.Table, "AS c ON c.ctx = sc.ctx AND c.left_value < sc.left_value AND c.right_value > sc.right_value")
	sb.Where("sc.id = ?", id)
	sb.Where("sc.ctx = ?", ctx)
	sb.Where("sc.status != ?", Delete)
	if status != nest.All && status != Delete {
		sb.Where("c.status = ?", status)
	} else {
		sb.Where("c.status != ?", Delete)
	}
	sb.Limit(1)
	sb.OrderBy("c.left_value DESC")

	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, nil
}
