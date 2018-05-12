package nest

type Service struct {
	m IManager
}

func NewService(m IManager) *Service {
	var s = &Service{}
	s.m = m
	return s
}

// --------------------------------------------------------------------------------
// AddRoot 添加顶级分类
func (this *Service) AddRoot(cType int, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(0, cType, K_ADD_POSITION_ROOT, 0, name, status, ext...)
}

func (this *Service) AddRootWithId(cId int64, cType int, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(cId, cType, K_ADD_POSITION_ROOT, 0, name, status, ext...)
}

// AddToFirst 添加子分类，新添加的子分类位于子分类列表的前面
func (this *Service) AddToFirst(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(0, -1, K_ADD_POSITION_FIRST, referTo, name, status, ext...)
}

func (this *Service) AddToFirstWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(cId, -1, K_ADD_POSITION_FIRST, referTo, name, status, ext...)
}

// AddToLast 添加子分类，新添加的子分类位于子分类列表的后面
func (this *Service) AddToLast(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(0, -1, K_ADD_POSITION_LAST, referTo, name, status, ext...)
}

func (this *Service) AddToLastWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(cId, -1, K_ADD_POSITION_LAST, referTo, name, status, ext...)
}

// AddToLeft 添加兄弟分类，新添加的分类位于指定分类的左边(前面)
func (this *Service) AddToLeft(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(0, -1, K_ADD_POSITION_LEFT, referTo, name, status, ext...)
}

func (this *Service) AddToLeftWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(cId, -1, K_ADD_POSITION_LEFT, referTo, name, status, ext...)
}

// AddToRight 添加兄弟分类，新添加的分类位于指定分类的右边(后面)
func (this *Service) AddToRight(referTo int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(0, -1, K_ADD_POSITION_RIGHT, referTo, name, status, ext...)
}

func (this *Service) AddToRightWithId(referTo, cId int64, name string, status int, ext ...map[string]interface{}) (result int64, err error) {
	return this.m.AddCategory(cId, -1, K_ADD_POSITION_RIGHT, referTo, name, status, ext...)
}

// --------------------------------------------------------------------------------
// GetCategory 获取分类信息
// id 分类 id
func (this *Service) GetCategory(id int64, result interface{}) (err error) {
	return this.m.GetCategory(id, result)
}

func (this *Service) GetCategoryWithName(cType int, name string, result interface{}) (err error) {
	return this.m.GetCategoryWithName(cType, name, result)
}

// GetCategoryAdvList 获取分类列表
// parentId: 父分类id，当此参数的值大于 0 的时候，将忽略 cType 参数
// cType: 指定筛选分类的类型
// status: 指定筛选分类的状态
// depth: 指定要获取多少级别内的分类
func (this *Service) GetCategoryAdvList(parentId int64, cType, status, depth int, name string, limit uint64, result interface{}) (err error) {
	return this.m.GetCategoryList(parentId, cType, status, depth, name, limit, false, result)
}

// GetNodeList 获取指定分类的子分类
func (this *Service) GetNodeList(parentId int64, status, depth int, result interface{}) (err error) {
	return this.m.GetCategoryList(parentId, 0, status, depth, "", 0, false, result)
}

// GetNodeIdList 获取指定分类的子分类 id 列表
func (this *Service) GetNodeIdList(parentId int64, status, depth int) (result []int64, err error) {
	return this.m.GetIdList(parentId, status, depth, false)
}

// GetCategoryList 获取指定分类的子分类列表，返回的列表包含指定的分类
func (this *Service) GetCategoryList(parentId int64, status, depth int, result interface{}) (err error) {
	return this.m.GetCategoryList(parentId, 0, status, depth, "", 0, true, result)
}

// GetIdList 获取指定分类的子分类 id 列表，返回的 id 列表包含指定的分类
func (this *Service) GetIdList(parentId int64, status, depth int) (result []int64, err error) {
	return this.m.GetIdList(parentId, status, depth, true)
}

// GetParentList 获取指定分类的父分类列表
func (this *Service) GetParentList(id int64, status int, result interface{}) (err error) {
	return this.m.GetPathList(id, status, false, result)
}

// GetPathList 获取指定分类到 root 分类的完整分类列表，包括自身
func (this *Service) GetPathList(id int64, status int, result interface{}) (err error) {
	return this.m.GetPathList(id, status, true, result)
}

// --------------------------------------------------------------------------------
// UpdateCategory 更新分类信息
func (this *Service) UpdateCategory(id int64, updateInfo map[string]interface{}) (err error) {
	return this.m.UpdateCategory(id, updateInfo)
}

// UpdateCategoryStatus 更新分类状态
// id: 被更新分类的 id
// status: 新的状态
// updateType:
// 		0、只更新当前分类的状态，子分类的状态不会受到影响，并且不会改变父子关系；
// 		1、子分类的状态会一起更新，不会改变父子关系；
// 		2、子分类的状态不会受到影响，并且所有子分类会向上移动一级（只针对把状态设置为 无效 的时候）；
func (this *Service) UpdateCategoryStatus(id int64, status, updateType int) (err error) {
	return this.m.UpdateCategoryStatus(id, status, updateType)
}

func (this *Service) MoveToRoot(id int64) (err error) {
	return this.m.MoveCategory(K_MOVE_POSITION_ROOT, id, 0)
}

func (this *Service) MoveToFirst(id, pid int64) (err error) {
	return this.m.MoveCategory(K_MOVE_POSITION_FIRST, id, pid)
}

func (this *Service) MoveToLast(id, pid int64) (err error) {
	return this.m.MoveCategory(K_MOVE_POSITION_LAST, id, pid)
}

func (this *Service) MoveToLeft(id, rid int64) (err error) {
	return this.m.MoveCategory(K_MOVE_POSITION_LEFT, id, rid)
}

func (this *Service) MoveToRight(id, rid int64) (err error) {
	return this.m.MoveCategory(K_MOVE_POSITION_RIGHT, id, rid)
}
