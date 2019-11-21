package nest

import "github.com/smartwalle/dbs"

type Repository interface {
	BeginTx() (dbs.TX, Repository)

	WithTx(tx dbs.TX) Repository

	AddRoot(ctx int64, name string, status Status) (result int64, err error)

	AddToFirst(ctx, pId int64, name string, status Status) (result int64, err error)

	AddToLast(ctx, pId int64, name string, status Status) (result int64, err error)

	AddToLeft(ctx, rId int64, name string, status Status) (result int64, err error)

	AddToRight(ctx, rId int64, name string, status Status) (result int64, err error)

	AddNode(ctx int64, position Position, rId int64, name string, status Status) (result int64, err error)

	GetNode(ctx, id int64) (result *Node, err error)

	GetNodeWithName(ctx int, name string, ) (result *Node, err error)

	GetParent(ctx, id int64) (result *Node, err error)

	GetChildren(ctx, pId int64, status Status, depth int) (result []*Node, err error)

	GetChildrenIdList(ctx, pId int64, status Status, depth int) (result []int64, err error)

	GetChildrenPathList(ctx, pId int64, status Status, depth int) (result []*Node, err error)

	GetChildrenPathIdList(ctx, pId int64, status Status, depth int) (result []int64, err error)

	GetParentList(ctx, id int64, status Status) (result []*Node, err error)

	GetParentPathList(ctx, id int64, status Status) (result []*Node, err error)

	GetNodeList(ctx, pId int64, status Status, depth int, name string, limit, offset int64, withParent bool) (result []*Node, err error)

	UpdateNodeStatus(ctx, id int64, status Status, updateType int) (err error)

	MoveToRoot(ctx, id int64) (err error)

	MoveToFirst(ctx, id, pId int64) (err error)

	MoveToLast(ctx, id, pId int64) (err error)

	MoveToLeft(ctx, id, rId int64) (err error)

	MoveToRight(ctx, id, rId int64) (err error)

	MoveTo(ctx, id, rId int64, position Position) (err error)
}
