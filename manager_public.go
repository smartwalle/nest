package nest

// --------------------------------------------------------------------------------
// AddRoot 添加顶级节点
func (this *Manager) AddRoot(ctx int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, ctx, K_ADD_POSITION_ROOT, 0, name, status, ext...)
}

func (this *Manager) AddRootWithId(cId, ctx int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(cId, ctx, K_ADD_POSITION_ROOT, 0, name, status, ext...)
}

// AddToFirst 添加子节点，新添加的子节点位于子节点列表的前面
func (this *Manager) AddToFirst(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, -1, K_ADD_POSITION_FIRST, referTo, name, status, ext...)
}

func (this *Manager) AddToFirstWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(cId, -1, K_ADD_POSITION_FIRST, referTo, name, status, ext...)
}

// AddToLast 添加子节点，新添加的子节点位于子节点列表的后面
func (this *Manager) AddToLast(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, -1, K_ADD_POSITION_LAST, referTo, name, status, ext...)
}

func (this *Manager) AddToLastWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(cId, -1, K_ADD_POSITION_LAST, referTo, name, status, ext...)
}

// AddToLeft 添加兄弟节点，新添加的节点位于指定节点的左边(前面)
func (this *Manager) AddToLeft(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, -1, K_ADD_POSITION_LEFT, referTo, name, status, ext...)
}

func (this *Manager) AddToLeftWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(cId, -1, K_ADD_POSITION_LEFT, referTo, name, status, ext...)
}

// AddToRight 添加兄弟节点，新添加的节点位于指定节点的右边(后面)
func (this *Manager) AddToRight(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(0, -1, K_ADD_POSITION_RIGHT, referTo, name, status, ext...)
}

func (this *Manager) AddToRightWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addNode(cId, -1, K_ADD_POSITION_RIGHT, referTo, name, status, ext...)
}

func (this *Manager) AddNode(cId, ctx int64, position int, referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	if position != K_ADD_POSITION_ROOT &&
		position != K_ADD_POSITION_FIRST &&
		position != K_ADD_POSITION_LAST &&
		position != K_ADD_POSITION_LEFT &&
		position != K_ADD_POSITION_RIGHT {
		return 0, ErrUnknownPosition
	}
	return this.addNode(cId, ctx, position, referTo, name, status, ext...)
}

// --------------------------------------------------------------------------------
// GetNode 获取节点信息
// id: 节点 id
func (this *Manager) GetNode(id int64, result interface{}) (err error) {
	return this.getNode(id, result)
}

func (this *Manager) GetNodeWithName(ctx int, name string, result interface{}) (err error) {
	return this.getNodeWithName(ctx, name, result)
}

// GetParent 获取指定节点的父节点
// id: 节点 id
func (this *Manager) GetParent(id int64, result interface{}) (err error) {
	return this.getParent(id, 0, result)
}

// GetSubNodeList 获取指定节点的子节点
func (this *Manager) GetSubNodeList(id int64, status, depth int, result interface{}) (err error) {
	return this.getNodeList(id, 0, status, depth, "", 0, false, result)
}

// GetSubNodeIdList 获取指定节点的子节点 id 列表
func (this *Manager) GetSubNodeIdList(id int64, status, depth int) (result []int64, err error) {
	return this.getIdList(id, status, depth, false)
}

// GetSubNodePathList 获取指定节点的子节点，返回的节点列表包括 id 对应的节点
func (this *Manager) GetSubNodePathList(id int64, status, depth int, result interface{}) (err error) {
	return this.getNodeList(id, 0, status, depth, "", 0, true, result)
}

// GetSubNodePathIdList 获取指定节点的子节点 id 列表，返回的 id 列表包含 id
func (this *Manager) GetSubNodePathIdList(id int64, status, depth int) (result []int64, err error) {
	return this.getIdList(id, status, depth, true)
}

// GetParentList 获取指定节点到 root 节点的完整节点列表
func (this *Manager) GetParentList(id int64, status int, result interface{}) (err error) {
	return this.getPathList(id, status, false, result)
}

// GetParentPathList 获取指定节点到 root 节点的完整节点列表，返回的节点列表包括 id 对应的节点
func (this *Manager) GetParentPathList(id int64, status int, result interface{}) (err error) {
	return this.getPathList(id, status, true, result)
}

// GetNodeList 获取节点列表
// parentId: 父节点id，当此参数的值大于 0 的时候，将忽略 ctx 参数
// ctx: 指定筛选节点的类型
// status: 指定筛选节点的状态
// depth: 指定要获取多少级别内的节点
// name: 模糊匹配 name 字段
// limit: 返回多少条数据
// includeParent: 如果有传递 parentId 参数，将可以通过此参数设定是否需要返回 parentId 对应的节点信息
func (this *Manager) GetNodeList(parentId int64, ctx, status, depth int, name string, limit uint64, includeParent bool, result interface{}) (err error) {
	return this.getNodeList(parentId, ctx, status, depth, name, limit, includeParent, result)
}

// --------------------------------------------------------------------------------
// UpdateNode 更新节点信息
func (this *Manager) UpdateNode(id int64, updateInfo map[string]interface{}) (err error) {
	return this.updateNode(id, updateInfo)
}

// UpdateNodeStatus 更新节点状态
// id: 被更新节点的 id
// status: 新的状态
// updateType:
// 		0、只更新当前节点的状态，子节点的状态不会受到影响，并且不会改变父子关系；
// 		1、子节点的状态会一起更新，不会改变父子关系；
// 		2、子节点的状态不会受到影响，并且所有子节点会向上移动一级（只针对把状态设置为 无效 的时候）；
func (this *Manager) UpdateNodeStatus(id int64, status, updateType int) (err error) {
	return this.updateNodeStatus(id, status, updateType)
}

func (this *Manager) MoveToRoot(id int64) (err error) {
	return this.moveNode(K_MOVE_POSITION_ROOT, id, 0)
}

func (this *Manager) MoveToFirst(id, pid int64) (err error) {
	return this.moveNode(K_MOVE_POSITION_FIRST, id, pid)
}

func (this *Manager) MoveToLast(id, pid int64) (err error) {
	return this.moveNode(K_MOVE_POSITION_LAST, id, pid)
}

func (this *Manager) MoveToLeft(id, rid int64) (err error) {
	return this.moveNode(K_MOVE_POSITION_LEFT, id, rid)
}

func (this *Manager) MoveToRight(id, rid int64) (err error) {
	return this.moveNode(K_MOVE_POSITION_RIGHT, id, rid)
}
