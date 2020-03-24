package sql

import (
	"github.com/smartwalle/nest"
)

func (this *Repository) AddRoot(ctx int64, name string, status nest.Status) (result int64, err error) {
	return this.AddNode(ctx, nest.Root, 0, name, status)
}

func (this *Repository) AddToFirst(ctx, pId int64, name string, status nest.Status) (result int64, err error) {
	return this.AddNode(ctx, nest.First, pId, name, status)
}

func (this *Repository) AddToLast(ctx, pId int64, name string, status nest.Status) (result int64, err error) {
	return this.AddNode(ctx, nest.Last, pId, name, status)
}

func (this *Repository) AddToLeft(ctx, rId int64, name string, status nest.Status) (result int64, err error) {
	return this.AddNode(ctx, nest.Left, rId, name, status)
}

func (this *Repository) AddToRight(ctx, rId int64, name string, status nest.Status) (result int64, err error) {
	return this.AddNode(ctx, nest.Right, rId, name, status)
}

func (this *Repository) AddNode(ctx int64, position nest.Position, rId int64, name string, status nest.Status) (result int64, err error) {
	if position != nest.Root &&
		position != nest.First &&
		position != nest.Last &&
		position != nest.Left &&
		position != nest.Right {
		return 0, nest.ErrUnknownPosition
	}
	if status != nest.Enable && status != nest.Disable {
		return 0, nest.ErrUnknownStatus
	}
	return this.addNode(ctx, position, rId, name, status)
}

func (this *Repository) GetNode(ctx, id int64) (result *nest.Node, err error) {
	return this.getNodeWithId(ctx, id)
}

func (this *Repository) GetNodeWithName(ctx int, name string) (result *nest.Node, err error) {
	return this.getNodeWithName(ctx, name)
}

func (this *Repository) GetParent(ctx, id int64) (result *nest.Node, err error) {
	return this.getParent(ctx, id, nest.All)
}

func (this *Repository) GetLastNode(ctx, pId int64) (result *nest.Node, err error) {
	return this.getLastNode(ctx, pId)
}

func (this *Repository) GetFirstNode(ctx, pId int64) (result *nest.Node, err error) {
	return this.getFirstNode(ctx, pId)
}

func (this *Repository) GetPreviousNode(ctx, id int64) (result *nest.Node, err error) {
	return this.getPreviousNode(ctx, id)
}

func (this *Repository) GetNextNode(ctx, id int64) (result *nest.Node, err error) {
	return this.getNextNode(ctx, id)
}

func (this *Repository) GetChildNodes(ctx, pId int64, status nest.Status, depth int) (result []*nest.Node, err error) {
	return this.getNodes(ctx, pId, status, depth, "", 0, 0, false)
}

func (this *Repository) GetChildNodeIds(ctx, pId int64, status nest.Status, depth int) (result []int64, err error) {
	return this.getNodeIds(ctx, pId, status, depth, "", 0, 0, false)
}

// GetChildNodePaths 获取指定节点的子节点，返回的节点列表包括当前节点
func (this *Repository) GetChildNodePaths(ctx, pId int64, status nest.Status, depth int) (result []*nest.Node, err error) {
	return this.getNodes(ctx, pId, status, depth, "", 0, 0, true)
}

// GetChildNodePathIds 获取指定节点的子节点 id 列表，返回的 id 列表包含当前节点 id
func (this *Repository) GetChildNodePathIds(ctx, pId int64, status nest.Status, depth int) (result []int64, err error) {
	return this.getNodeIds(ctx, pId, status, depth, "", 0, 0, true)
}

func (this *Repository) GetParentNodes(ctx, id int64, status nest.Status) (result []*nest.Node, err error) {
	return this.getParentNodes(ctx, id, status, false)
}

func (this *Repository) GetNodePaths(ctx, id int64, status nest.Status) (result []*nest.Node, err error) {
	return this.getParentNodes(ctx, id, status, true)
}

func (this *Repository) GetNodes(ctx, pId int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result []*nest.Node, err error) {
	return this.getNodes(ctx, pId, status, depth, name, limit, offset, withParent)
}

func (this *Repository) GetNodeIds(ctx, pId int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result []int64, err error) {
	return this.getNodeIds(ctx, pId, status, depth, name, limit, offset, withParent)
}

func (this *Repository) UpdateNodeName(ctx, id int64, name string) (err error) {
	return this.updateNodeName(ctx, id, name)
}

func (this *Repository) UpdateNodeStatus(ctx, id int64, status nest.Status) (err error) {
	if status != nest.Enable && status != nest.Disable {
		return nest.ErrUnknownStatus
	}
	return this.updateNodeStatus(ctx, id, status, 1)
}

func (this *Repository) MoveToRoot(ctx, id int64) (err error) {
	return this.moveNode(nest.Root, ctx, id, 0)
}

func (this *Repository) MoveToFirst(ctx, id, pId int64) (err error) {
	return this.moveNode(nest.First, ctx, id, pId)
}

func (this *Repository) MoveToLast(ctx, id, pId int64) (err error) {
	return this.moveNode(nest.Last, ctx, id, pId)
}

func (this *Repository) MoveToLeft(ctx, id, rId int64) (err error) {
	return this.moveNode(nest.Left, ctx, id, rId)
}

func (this *Repository) MoveToRight(ctx, id, rId int64) (err error) {
	return this.moveNode(nest.Right, ctx, id, rId)
}

func (this *Repository) MoveUp(ctx, id int64) (err error) {
	return this.moveNode(nest.Left, ctx, id, 0)
}

func (this *Repository) MoveDown(ctx, id int64) (err error) {
	return this.moveNode(nest.Right, ctx, id, 0)
}

func (this *Repository) MoveTo(ctx, id, rId int64, position nest.Position) (err error) {
	return this.moveNode(position, ctx, id, rId)
}

func (this *Repository) RemoveNode(ctx, id int64) (err error) {
	return this.removeNode(ctx, id)
}
