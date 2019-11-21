package mysql

import (
	"github.com/smartwalle/nest"
)

// --------------------------------------------------------------------------------
// AddRoot 添加顶级节点
func (this *nestRepository) AddRoot(ctx int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.Root, 0, name, status)
}

// AddToFirst 添加子节点，新添加的子节点位于子节点列表的前面
func (this *nestRepository) AddToFirst(ctx, pId int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.First, pId, name, status)
}

// AddToLast 添加子节点，新添加的子节点位于子节点列表的后面
func (this *nestRepository) AddToLast(ctx, pId int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.Last, pId, name, status)
}

// AddToLeft 添加兄弟节点，新添加的节点位于指定节点的左边(前面)
func (this *nestRepository) AddToLeft(ctx, rId int64, name string, status nest.Status) (result int64, err error) {
	return this.addNode(ctx, nest.Left, rId, name, status)
}

// AddToRight 添加兄弟节点，新添加的节点位于指定节点的右边(后面)
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

// --------------------------------------------------------------------------------
// GetNode 获取节点信息
// id: 节点 id
func (this *nestRepository) GetNode(ctx, id int64) (result *nest.Node, err error) {
	return this.getNodeWithId(ctx, id)
}

func (this *nestRepository) GetNodeWithName(ctx int, name string) (result *nest.Node, err error) {
	return this.getNodeWithName(ctx, name)
}

// GetParent 获取指定节点的父节点
// id: 节点 id
func (this *nestRepository) GetParent(ctx, id int64) (result *nest.Node, err error) {
	return this.getParent(ctx, id, 0)
}

// GetChildren 获取指定节点的子节点
func (this *nestRepository) GetChildren(ctx, pId int64, status nest.Status, depth int) (result []*nest.Node, err error) {
	return this.getNodeList(ctx, pId, status, depth, "", 0, 0, false)
}

// GetChildrenIdList 获取指定节点的子节点 id 列表
func (this *nestRepository) GetChildrenIdList(ctx, pId int64, status nest.Status, depth int) (result []int64, err error) {
	return this.getChildrenIdList(ctx, pId, status, depth, false)
}

// GetChildrenPathList 获取指定节点的子节点，返回的节点列表包括 id 对应的节点
func (this *nestRepository) GetChildrenPathList(ctx, pId int64, status nest.Status, depth int) (result []*nest.Node, err error) {
	return this.getNodeList(ctx, pId, status, depth, "", 0, 0, true)
}

// GetChildrenPathIdList 获取指定节点的子节点 id 列表，返回的 id 列表包含 id
func (this *nestRepository) GetChildrenPathIdList(ctx, pId int64, status nest.Status, depth int) (result []int64, err error) {
	return this.getChildrenIdList(ctx, pId, status, depth, true)
}

// GetParentList 获取指定节点到 root 节点的完整节点列表
func (this *nestRepository) GetParentList(ctx, id int64, status nest.Status) (result []*nest.Node, err error) {
	return this.getPathList(ctx, id, status, false)
}

// GetParentPathList 获取指定节点到 root 节点的完整节点列表，返回的节点列表包括 id 对应的节点
func (this *nestRepository) GetParentPathList(ctx, id int64, status nest.Status) (result []*nest.Node, err error) {
	return this.getPathList(ctx, id, status, true)
}

// GetNodeList 获取节点列表
// ctx: 指定筛选节点的类型
// pId: 父节点id，当此参数的值大于 0 的时候，将忽略 ctx 参数
// status: 指定筛选节点的状态
// depth: 指定要获取多少级别内的节点
// name: 模糊匹配 name 字段
// limit: 返回多少条数据
// withParent: 如果有传递 parentId 参数，将可以通过此参数设定是否需要返回 parentId 对应的节点信息
func (this *nestRepository) GetNodeList(ctx, pId int64, status nest.Status, depth int, name string, limit, offset int64, withParent bool) (result []*nest.Node, err error) {
	return this.getNodeList(ctx, pId, status, depth, name, limit, offset, withParent)
}

// --------------------------------------------------------------------------------
// UpdateNodeName 更新节点名称信息
func (this *nestRepository) UpdateNodeName(ctx, id int64, name string) (err error) {
	return this.updateNodeName(ctx, id, name)
}

// UpdateNodeStatus 更新节点状态
// id: 被更新节点的 id
// status: 新的状态
// updateType:
// 		0、只更新当前节点的状态，子节点的状态不会受到影响，并且不会改变父子关系；
// 		1、子节点的状态会一起更新，不会改变父子关系；
// 		2、子节点的状态不会受到影响，并且所有子节点会向上移动一级（只针对把状态设置为 无效 的时候）；
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

func (this *nestRepository) MoveTo(ctx, id, rId int64, position nest.Position) (err error) {
	return this.moveNode(position, ctx, id, rId)
}
