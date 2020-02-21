package mysql

import (
	"github.com/smartwalle/nest"
)

func (this *nestRepository) AddRoot(ctx int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.Root, 0, name, status)
}

func (this *nestRepository) AddToFirst(ctx, pId int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.First, pId, name, status)
}

func (this *nestRepository) AddToLast(ctx, pId int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.Last, pId, name, status)
}

func (this *nestRepository) AddToLeft(ctx, rId int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.Left, rId, name, status)
}

func (this *nestRepository) AddToRight(ctx, rId int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.Right, rId, name, status)
}

func (this *nestRepository) AddNode(ctx int64, position nest.Position, rId int64, name string, status nest.Status) (result int64, err error) {
	if position != nest.Root &&
		position != nest.First &&
		position != nest.Last &&
		position != nest.Left &&
		position != nest.Right {
		return 0, nest.ErrUnknownPosition
	}
	return this.addNode(ctx, position, rId, name, status)
}

func (this *nestRepository) GetNode(ctx, id int64) (result *nest.Node, err error) {
	return this.getNodeWithId(ctx, id)
}

func (this *nestRepository) GetNodeWithName(ctx int, name string) (result *nest.Node, err error) {
	return this.getNodeWithName(ctx, name)
}

func (this *nestRepository) GetParent(ctx, id int64) (result *nest.Node, err error) {
	return this.getParent(ctx, id, 0)
}

func (this *nestRepository) GetLastNode(ctx, pId int64) (result *nest.Node, err error) {
	return this.getLastNode(ctx, pId)
}

func (this *nestRepository) GetFirstNode(ctx, pId int64) (result *nest.Node, err error) {
	return this.getFirstNode(ctx, pId)
}

func (this *nestRepository) GetPreviousNode(ctx, id int64) (result *nest.Node, err error) {
	return this.getPreviousNode(ctx, id)
}

func (this *nestRepository) GetNextNode(ctx, id int64) (result *nest.Node, err error) {
	return this.getNextNode(ctx, id)
}

func (this *nestRepository) GetChildren(ctx, pId int64, status nest.Status, depth int) (result []*nest.Node, err error) {
	return this.getNodeList(ctx, pId, status, depth, "", 0, 0, false)
}

func (this *nestRepository) GetChildrenIdList(ctx, pId int64, status nest.Status, depth int) (result []int64, err error) {
	return this.getChildrenIdList(ctx, pId, status, depth, false)
}

func (this *nestRepository) GetChildrenPathList(ctx, pId int64, status nest.Status, depth int) (result []*nest.Node, err error) {
	return this.getNodeList(ctx, pId, status, depth, "", 0, 0, true)
}

func (this *nestRepository) GetChildrenPathIdList(ctx, pId int64, status nest.Status, depth int) (result []int64, err error) {
	return this.getChildrenIdList(ctx, pId, status, depth, true)
}

func (this *nestRepository) GetParentList(ctx, id int64, status nest.Status) (result []*nest.Node, err error) {
	return this.getPathList(ctx, id, status, false)
}

func (this *nestRepository) GetParentPathList(ctx, id int64, status nest.Status) (result []*nest.Node, err error) {
	return this.getPathList(ctx, id, status, true)
}

func (this *nestRepository) GetNodeList(ctx, pId int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result []*nest.Node, err error) {
	return this.getNodeList(ctx, pId, status, depth, name, limit, offset, withParent)
}

func (this *nestRepository) UpdateNodeName(ctx, id int64, name string) (err error) {
	return this.updateNodeName(ctx, id, name)
}

func (this *nestRepository) UpdateNodeStatus(ctx, id int64, status nest.Status, updateType int) (err error) {
	return this.updateNodeStatus(ctx, id, status, updateType)
}

func (this *nestRepository) MoveToRoot(ctx, id int64) (err error) {
	return this.moveNode(nest.Root, ctx, id, 0)
}

func (this *nestRepository) MoveToFirst(ctx, id, pId int64) (err error) {
	return this.moveNode(nest.First, ctx, id, pId)
}

func (this *nestRepository) MoveToLast(ctx, id, pId int64) (err error) {
	return this.moveNode(nest.Last, ctx, id, pId)
}

func (this *nestRepository) MoveToLeft(ctx, id, rId int64) (err error) {
	return this.moveNode(nest.Left, ctx, id, rId)
}

func (this *nestRepository) MoveToRight(ctx, id, rId int64) (err error) {
	return this.moveNode(nest.Right, ctx, id, rId)
}

func (this *nestRepository) MoveUp(ctx, id int64) (err error) {
	return this.moveNode(nest.Left, ctx, id, 0)
}

func (this *nestRepository) MoveDown(ctx, id int64) (err error) {
	return this.moveNode(nest.Right, ctx, id, 0)
}

func (this *nestRepository) MoveTo(ctx, id, rId int64, position nest.Position) (err error) {
	return this.moveNode(position, ctx, id, rId)
}
