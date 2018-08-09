package nest

import "time"

const (
	K_STATUS_ENABLE  = 1000 // 启用
	K_STATUS_DISABLE = 2000 // 禁用
)

const (
	K_ADD_POSITION_ROOT  = 0 // 顶级节点
	K_ADD_POSITION_FIRST = 1 // 列表头部 (子节点)
	K_ADD_POSITION_LAST  = 2 // 列表尾部 (子节点)
	K_ADD_POSITION_LEFT  = 3 // 左边 (兄弟节点)
	K_ADD_POSITION_RIGHT = 4 // 右边 (兄弟节点)
)

const (
	K_MOVE_POSITION_ROOT  = 0 // 顶级节点
	K_MOVE_POSITION_FIRST = 1 // 列表头部 (子节点)
	K_MOVE_POSITION_LAST  = 2 // 列表尾部 (子节点)
	K_MOVE_POSITION_LEFT  = 3 // 左边 (兄弟节点)
	K_MOVE_POSITION_RIGHT = 4 // 右边 (兄弟节点)
)

type Node struct {
	Id         int64      `json:"id"                        sql:"id"`
	Ctx        int64      `json:"ctx"                       sql:"ctx"`
	Name       string     `json:"name"                      sql:"name"`
	LeftValue  int        `json:"left_value"                sql:"left_value"`
	RightValue int        `json:"right_value"               sql:"right_value"`
	Depth      int        `json:"depth"                     sql:"depth"`
	Status     int        `json:"status"                    sql:"status"`
	CreatedOn  *time.Time `json:"created_on,omitempty"      sql:"created_on"`
	UpdatedOn  *time.Time `json:"updated_on,omitempty"      sql:"updated_on"`
}

func (this *Node) IsLeaf() bool {
	return this.LeftValue+1 == this.RightValue
}

func (this *Node) IsValid() bool {
	return this.Status == K_STATUS_ENABLE
}
