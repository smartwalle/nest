package nest

import "github.com/smartwalle/dbs"

type Repository interface {
	BeginTx() (dbs.TX, Repository)

	WithTx(tx dbs.TX) Repository

	// AddRoot 添加顶级节点
	AddRoot(ctx int64, name string, status Status) (result int64, err error)

	// AddToFirst 添加子节点，新添加的子节点位于子节点列表的最前面
	AddToFirst(ctx, pId int64, name string, status Status) (result int64, err error)

	// AddToLast 添加子节点，新添加的子节点位于子节点列表的最后面
	AddToLast(ctx, pId int64, name string, status Status) (result int64, err error)

	// AddToLeft 添加兄弟节点，新添加的节点位于指定节点的左边(前面)
	AddToLeft(ctx, rId int64, name string, status Status) (result int64, err error)

	// AddToRight 添加兄弟节点，新添加的节点位于指定节点的右边(后面)
	AddToRight(ctx, rId int64, name string, status Status) (result int64, err error)

	// AddNode 添加节点
	// ctx: 节点的类型，相当于对节点进行分组，不同组的节点相互之间不影响；
	// position: 新添加的节点相对于 rId 的位置；
	// rId: 参照节点的 id，根据 position 参数决定新的节点是其子节点还是兄弟节点；
	// name: 节点的名字；
	// status: 节点的状态；
	AddNode(ctx int64, position Position, rId int64, name string, status Status) (result int64, err error)

	// GetNode 获取节点信息
	GetNode(ctx, id int64) (result *Node, err error)

	// GetNodeWithName 获取节点信息
	GetNodeWithName(ctx int, name string) (result *Node, err error)

	// GetParent 获取指定节点的父节点
	GetParent(ctx, id int64) (result *Node, err error)

	// GetPreviousNode 获取相邻的上一节点(前面的节点)
	GetPreviousNode(ctx, id int64) (result *Node, err error)

	// GetNextNode 获取相邻的下一节点(后面的节点)
	GetNextNode(ctx, id int64) (result *Node, err error)

	// GetChildren 获取指定节点的子节点
	GetChildren(ctx, pId int64, status Status, depth int) (result []*Node, err error)

	// GetChildrenIdList 获取指定节点的子节点 id 列表
	GetChildrenIdList(ctx, pId int64, status Status, depth int) (result []int64, err error)

	// GetChildrenPathList 获取指定节点的子节点，返回的节点列表包括当前节点
	GetChildrenPathList(ctx, pId int64, status Status, depth int) (result []*Node, err error)

	// GetChildrenPathIdList 获取指定节点的子节点 id 列表，返回的 id 列表包含当前节点 id
	GetChildrenPathIdList(ctx, pId int64, status Status, depth int) (result []int64, err error)

	// GetParentList 获取指定节点到 root 节点的完整节点列表，返回的节点列表不包括当前节点
	GetParentList(ctx, id int64, status Status) (result []*Node, err error)

	// GetParentPathList 获取指定节点到 root 节点的完整节点列表，返回的节点列表包括当前节点
	GetParentPathList(ctx, id int64, status Status) (result []*Node, err error)

	// GetNodeList 获取节点列表
	// ctx: 指定筛选节点的类型；
	// pId: 父节点id，当此参数的值大于 0 的时候，将忽略 ctx 参数；
	// status: 指定筛选节点的状态；
	// depth: 指定要获取多少级别内的节点；
	// name: 模糊匹配 name 字段；
	// limit: 返回多少条数据；
	// withParent: 如果有传递 pId 参数，将可以通过此参数设定是否需要返回 pId 对应的节点信息；
	GetNodeList(ctx, pId int64, status Status, depth int, name string, limit, offset int64, withParent bool) (result []*Node, err error)

	// UpdateNodeName 更新节点名称
	UpdateNodeName(ctx, id int64, name string) (err error)

	// UpdateNodeStatus 更新节点状态
	// id: 被更新节点的 id；
	// status: 新的状态；
	// updateType:
	// 		0、只更新当前节点的状态，子节点的状态不会受到影响，并且不会改变父子关系；
	// 		1、子节点的状态会一起更新，不会改变父子关系；
	// 		2、子节点的状态不会受到影响，并且所有子节点会向上移动一级（只针对把状态设置为 无效 的时候）；
	UpdateNodeStatus(ctx, id int64, status Status, updateType int) (err error)

	// MoveToRoot 将指定节点移动到最顶层
	MoveToRoot(ctx, id int64) (err error)

	// MoveToFirst 将节点调整为指定节点的子节点，并将该节点作为指定节点列表的第一个节点
	MoveToFirst(ctx, id, pId int64) (err error)

	// MoveToLast 将节点调整为指定节点的子节点，并将该节点作为指定节点列表的最后一个节点
	MoveToLast(ctx, id, pId int64) (err error)

	// MoveToLeft 将节点调整为指定节点的兄弟节点，并将该节点位于指定节点的左边(前面)
	// 如果 rId 参数小于等于 0，则本方法的作用与方法 MoveUp 一致
	MoveToLeft(ctx, id, rId int64) (err error)

	// MoveToRight 将节点调整为指定节点的兄弟节点，并将该节点位于指定节点的右边(后面)
	// 如果 rId 参数小于等于 0，则本方法的作用与方法 MoveDown 一致
	MoveToRight(ctx, id, rId int64) (err error)

	// MoveUp 将节点向前移动一位，即向左移动一位
	MoveUp(ctx, id int64) (err error)

	// MoveDown 将节点向后移动一位，即向右移动一位
	MoveDown(ctx, id int64) (err error)

	// MoveTo 移动节点
	MoveTo(ctx, id, rId int64, position Position) (err error)
}
