package nest

import "time"

type Status int

const (
	All     Status = 0 // 所有，用于查询
	Enable  Status = 1 // 启用
	Disable Status = 2 // 禁用
)

type Position int

const (
	Root  Position = 0 // 顶级节点
	First Position = 1 // 列表头部 (子节点)
	Last  Position = 2 // 列表尾部 (子节点)
	Left  Position = 3 // 左边 (兄弟节点)
	Right Position = 4 // 右边 (兄弟节点)
)

type Node struct {
	Id         int64      `json:"id"                        sql:"id"`
	Ctx        int64      `json:"ctx"                       sql:"ctx"`
	Name       string     `json:"name"                      sql:"name"`
	LeftValue  int64      `json:"left_value"                sql:"left_value"`
	RightValue int64      `json:"right_value"               sql:"right_value"`
	Depth      int        `json:"depth"                     sql:"depth"`
	Status     Status     `json:"status"                    sql:"status"`
	CreatedOn  *time.Time `json:"created_on,omitempty"      sql:"created_on"`
	UpdatedOn  *time.Time `json:"updated_on,omitempty"      sql:"updated_on"`
}

func (this *Node) HasChildNodes() bool {
	return this.RightValue-this.LeftValue > 1
}

func (this *Node) IsValid() bool {
	return this.Status == Enable
}
