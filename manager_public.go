package nest

// --------------------------------------------------------------------------------
// AddRoot 添加顶级分类
func (this *Manager) AddRoot(cType int, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(0, cType, K_ADD_POSITION_ROOT, 0, name, status, ext...)
}

func (this *Manager) AddRootWithId(cId int64, cType int, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(cId, cType, K_ADD_POSITION_ROOT, 0, name, status, ext...)
}

// AddToFirst 添加子分类，新添加的子分类位于子分类列表的前面
func (this *Manager) AddToFirst(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(0, -1, K_ADD_POSITION_FIRST, referTo, name, status, ext...)
}

func (this *Manager) AddToFirstWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(cId, -1, K_ADD_POSITION_FIRST, referTo, name, status, ext...)
}

// AddToLast 添加子分类，新添加的子分类位于子分类列表的后面
func (this *Manager) AddToLast(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(0, -1, K_ADD_POSITION_LAST, referTo, name, status, ext...)
}

func (this *Manager) AddToLastWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(cId, -1, K_ADD_POSITION_LAST, referTo, name, status, ext...)
}

// AddToLeft 添加兄弟分类，新添加的分类位于指定分类的左边(前面)
func (this *Manager) AddToLeft(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(0, -1, K_ADD_POSITION_LEFT, referTo, name, status, ext...)
}

func (this *Manager) AddToLeftWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(cId, -1, K_ADD_POSITION_LEFT, referTo, name, status, ext...)
}

// AddToRight 添加兄弟分类，新添加的分类位于指定分类的右边(后面)
func (this *Manager) AddToRight(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(0, -1, K_ADD_POSITION_RIGHT, referTo, name, status, ext...)
}

func (this *Manager) AddToRightWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.addCategory(cId, -1, K_ADD_POSITION_RIGHT, referTo, name, status, ext...)
}

// --------------------------------------------------------------------------------
// GetCategory 获取分类信息
// id 分类 id
func (this *Manager) GetCategory(id int64, result interface{}) (err error) {
	return this.getCategory(id, result)
}

func (this *Manager) GetCategoryWithName(cType int, name string, result interface{}) (err error) {
	return this.getCategoryWithName(cType, name, result)
}

// GetParent 获取指定分类的父分类
func (this *Manager) GetParent(id int64, result interface{}) (err error) {
	return this.getParent(id, 0, result)
}

// GetNodeList 获取指定分类的子分类
func (this *Manager) GetNodeList(parentId int64, status, depth int, result interface{}) (err error) {
	return this.getCategoryList(parentId, 0, status, depth, "", 0, false, result)
}

// GetNodeIdList 获取指定分类的子分类 id 列表
func (this *Manager) GetNodeIdList(parentId int64, status, depth int) (result []int64, err error) {
	return this.getIdList(parentId, status, depth, false)
}

// GetNodePathList 获取指定分类的子分类，返回的分类列表包括 parentId 对应的分类
func (this *Manager) GetNodePathList(parentId int64, status, depth int, result interface{}) (err error) {
	return this.getCategoryList(parentId, 0, status, depth, "", 0, true, result)
}

// GetNodePathIdList 获取指定分类的子分类 id 列表，返回的 id 列表包含 parentId
func (this *Manager) GetNodePathIdList(parentId int64, status, depth int) (result []int64, err error) {
	return this.getIdList(parentId, status, depth, true)
}

// GetCategoryList 获取分类列表
// parentId: 父分类id，当此参数的值大于 0 的时候，将忽略 cType 参数
// cType: 指定筛选分类的类型
// status: 指定筛选分类的状态
// depth: 指定要获取多少级别内的分类
// name: 模糊匹配 name 字段
// limit: 返回多少条数据
// includeParent: 如果有传递 parentId 参数，将可以通过此参数设定是否需要返回 parentId 对应的分类信息
func (this *Manager) GetCategoryList(parentId int64, cType, status, depth int, name string, limit uint64, includeParent bool, result interface{}) (err error) {
	return this.getCategoryList(parentId, cType, status, depth, name, limit, includeParent, result)
}

// GetParentList 获取指定分类到 root 分类的完整分类列表
func (this *Manager) GetParentList(id int64, status int, result interface{}) (err error) {
	return this.getPathList(id, status, false, result)
}

// GetParentPathList 获取指定分类到 root 分类的完整分类列表，返回的分类列表包括 id 对应的分类
func (this *Manager) GetParentPathList(id int64, status int, result interface{}) (err error) {
	return this.getPathList(id, status, true, result)
}

// --------------------------------------------------------------------------------
// UpdateCategory 更新分类信息
func (this *Manager) UpdateCategory(id int64, updateInfo map[string]interface{}) (err error) {
	return this.updateCategory(id, updateInfo)
}

// UpdateCategoryStatus 更新分类状态
// id: 被更新分类的 id
// status: 新的状态
// updateType:
// 		0、只更新当前分类的状态，子分类的状态不会受到影响，并且不会改变父子关系；
// 		1、子分类的状态会一起更新，不会改变父子关系；
// 		2、子分类的状态不会受到影响，并且所有子分类会向上移动一级（只针对把状态设置为 无效 的时候）；
func (this *Manager) UpdateCategoryStatus(id int64, status, updateType int) (err error) {
	return this.updateCategoryStatus(id, status, updateType)
}

func (this *Manager) MoveToRoot(id int64) (err error) {
	return this.moveCategory(K_MOVE_POSITION_ROOT, id, 0)
}

func (this *Manager) MoveToFirst(id, pid int64) (err error) {
	return this.moveCategory(K_MOVE_POSITION_FIRST, id, pid)
}

func (this *Manager) MoveToLast(id, pid int64) (err error) {
	return this.moveCategory(K_MOVE_POSITION_LAST, id, pid)
}

func (this *Manager) MoveToLeft(id, rid int64) (err error) {
	return this.moveCategory(K_MOVE_POSITION_LEFT, id, rid)
}

func (this *Manager) MoveToRight(id, rid int64) (err error) {
	return this.moveCategory(K_MOVE_POSITION_RIGHT, id, rid)
}
