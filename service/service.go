package service

type NestRepository interface {
	AddRoot(ctx int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddRootWithId(ctx, id int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddToFirst(pid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddToFirstWithId(id, pid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddToLast(pid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddToLastWithId(id, pid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddToLeft(rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddToLeftWithId(id, rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddToRight(rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddToRightWithId(id, rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	AddNode(ctx, id int64, position int, rid int64, name string, status int, ext ...map[string]interface{}) (result int64, err error)

	GetNode(ctx, id int64, result interface{}) (err error)

	GetNodeWithName(ctx int, name string, result interface{}) (err error)

	GetParent(ctx, id int64, result interface{}) (err error)

	GetSubNodeList(ctx, id int64, status, depth int, result interface{}) (err error)

	GetSubNodeIdList(ctx, id int64, status, depth int) (result []int64, err error)

	GetSubNodePathList(ctx, id int64, status, depth int, result interface{}) (err error)

	GetSubNodePathIdList(ctx, id int64, status, depth int) (result []int64, err error)

	GetParentList(ctx, id int64, status int, result interface{}) (err error)

	GetParentPathList(ctx, id int64, status int, result interface{}) (err error)

	GetNodeList(ctx, pid int64, status, depth int, name string, limit, offset int64, includeParent bool, result interface{}) (err error)

	UpdateNode(id int64, updateInfo map[string]interface{}) (err error)

	UpdateNodeStatus(id int64, status, updateType int) (err error)

	MoveToRoot(id int64) (err error)

	MoveToFirst(id, pid int64) (err error)

	MoveToLast(id, pid int64) (err error)

	MoveToLeft(id, rid int64) (err error)

	MoveToRight(id, rid int64) (err error)

	MoveTo(id, rid int64, position int) (err error)
}
