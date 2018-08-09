package mysql

import (
	"github.com/smartwalle/nest"
)

// --------------------------------------------------------------------------------
// AddRoot 添加顶级节点
func (this *nestRepository) AddRoot(ctx int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, ctx, nest.K_ADD_POSITION_ROOT, 0, name, status, ext...)
}

func (this *nestRepository) AddRootWithId(ctx, id int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(id, ctx, nest.K_ADD_POSITION_ROOT, 0, name, status, ext...)
}

// AddToFirst 添加子节点，新添加的子节点位于子节点列表的前面
func (this *nestRepository) AddToFirst(pid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, -1, nest.K_ADD_POSITION_FIRST, pid, name, status, ext...)
}

func (this *nestRepository) AddToFirstWithId(id, pid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(id, -1, nest.K_ADD_POSITION_FIRST, pid, name, status, ext...)
}

// AddToLast 添加子节点，新添加的子节点位于子节点列表的后面
func (this *nestRepository) AddToLast(pid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, -1, nest.K_ADD_POSITION_LAST, pid, name, status, ext...)
}

func (this *nestRepository) AddToLastWithId(id, pid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(id, -1, nest.K_ADD_POSITION_LAST, pid, name, status, ext...)
}

// AddToLeft 添加兄弟节点，新添加的节点位于指定节点的左边(前面)
func (this *nestRepository) AddToLeft(rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, -1, nest.K_ADD_POSITION_LEFT, rid, name, status, ext...)
}

func (this *nestRepository) AddToLeftWithId(id, rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(id, -1, nest.K_ADD_POSITION_LEFT, rid, name, status, ext...)
}

// AddToRight 添加兄弟节点，新添加的节点位于指定节点的右边(后面)
func (this *nestRepository) AddToRight(rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, -1, nest.K_ADD_POSITION_RIGHT, rid, name, status, ext...)
}

func (this *nestRepository) AddToRightWithId(id, rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(id, -1, nest.K_ADD_POSITION_RIGHT, rid, name, status, ext...)
}

func (this *nestRepository) AddNode(ctx, id int64, position int, rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	if position != nest.K_ADD_POSITION_ROOT &&
		position != nest.K_ADD_POSITION_FIRST &&
		position != nest.K_ADD_POSITION_LAST &&
		position != nest.K_ADD_POSITION_LEFT &&
		position != nest.K_ADD_POSITION_RIGHT {
		return 0, nest.ErrUnknownPosition
	}
	return this.addNode(id, ctx, position, rid, name, status, ext...)
}

// --------------------------------------------------------------------------------
// GetNode 获取节点信息
// id: 节点 id
func (this *nestRepository) GetNode(ctx, id int64, result interface{}) (err error) {
	return this.getNode(ctx, id, result)
}

func (this *nestRepository) GetNodeWithName(ctx int, name string, result interface{}) (err error) {
	return this.getNodeWithName(ctx, name, result)
}

// GetParent 获取指定节点的父节点
// id: 节点 id
func (this *nestRepository) GetParent(ctx, id int64, result interface{}) (err error) {
	return this.getParent(ctx, id, 0, result)
}

// GetSubNodeList 获取指定节点的子节点
func (this *nestRepository) GetSubNodeList(ctx, id int64, status, depth int, result interface{}) (err error) {
	return this.getNodeList(ctx, id, status, depth, "", 0, 0, false, result)
}

// GetSubNodeIdList 获取指定节点的子节点 id 列表
func (this *nestRepository) GetSubNodeIdList(ctx, id int64, status, depth int) (result []int64, err error) {
	return this.getIdList(ctx, id, status, depth, false)
}

// GetSubNodePathList 获取指定节点的子节点，返回的节点列表包括 id 对应的节点
func (this *nestRepository) GetSubNodePathList(ctx, id int64, status, depth int, result interface{}) (err error) {
	return this.getNodeList(ctx, id, status, depth, "", 0, 0, true, result)
}

// GetSubNodePathIdList 获取指定节点的子节点 id 列表，返回的 id 列表包含 id
func (this *nestRepository) GetSubNodePathIdList(ctx, id int64, status, depth int) (result []int64, err error) {
	return this.getIdList(ctx, id, status, depth, true)
}

// GetParentList 获取指定节点到 root 节点的完整节点列表
func (this *nestRepository) GetParentList(ctx, id int64, status int, result interface{}) (err error) {
	return this.getPathList(ctx, id, status, false, result)
}

// GetParentPathList 获取指定节点到 root 节点的完整节点列表，返回的节点列表包括 id 对应的节点
func (this *nestRepository) GetParentPathList(ctx, id int64, status int, result interface{}) (err error) {
	return this.getPathList(ctx, id, status, true, result)
}

// GetNodeList 获取节点列表
// ctx: 指定筛选节点的类型
// pid: 父节点id，当此参数的值大于 0 的时候，将忽略 ctx 参数
// status: 指定筛选节点的状态
// depth: 指定要获取多少级别内的节点
// name: 模糊匹配 name 字段
// limit: 返回多少条数据
// includeParent: 如果有传递 parentId 参数，将可以通过此参数设定是否需要返回 parentId 对应的节点信息
func (this *nestRepository) GetNodeList(ctx, pid int64, status, depth int, name string, limit, offset int64, includeParent bool, result interface{}) (err error) {
	return this.getNodeList(ctx, pid, status, depth, name, limit, offset, includeParent, result)
}

// --------------------------------------------------------------------------------
// UpdateNode 更新节点信息
func (this *nestRepository) UpdateNode(id int64, updateInfo map[string]interface{}) (err error) {
	return this.updateNode(id, updateInfo)
}

// UpdateNodeStatus 更新节点状态
// id: 被更新节点的 id
// status: 新的状态
// updateType:
// 		0、只更新当前节点的状态，子节点的状态不会受到影响，并且不会改变父子关系；
// 		1、子节点的状态会一起更新，不会改变父子关系；
// 		2、子节点的状态不会受到影响，并且所有子节点会向上移动一级（只针对把状态设置为 无效 的时候）；
func (this *nestRepository) UpdateNodeStatus(id int64, status, updateType int) (err error) {
	return this.updateNodeStatus(id, status, updateType)
}

func (this *nestRepository) MoveToRoot(id int64) (err error) {
	return this.moveNode(nest.K_MOVE_POSITION_ROOT, id, 0)
}

func (this *nestRepository) MoveToFirst(id, pid int64) (err error) {
	return this.moveNode(nest.K_MOVE_POSITION_FIRST, id, pid)
}

func (this *nestRepository) MoveToLast(id, pid int64) (err error) {
	return this.moveNode(nest.K_MOVE_POSITION_LAST, id, pid)
}

func (this *nestRepository) MoveToLeft(id, rid int64) (err error) {
	return this.moveNode(nest.K_MOVE_POSITION_LEFT, id, rid)
}

func (this *nestRepository) MoveToRight(id, rid int64) (err error) {
	return this.moveNode(nest.K_MOVE_POSITION_RIGHT, id, rid)
}

func (this *nestRepository) MoveTo(id, rid int64, position int) (err error) {
	return this.moveNode(position, id, rid)
}
