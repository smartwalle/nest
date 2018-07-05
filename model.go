package nest

import "time"

const (
	K_STATUS_ENABLE  = 1000 // 启用
	K_STATUS_DISABLE = 2000 // 禁用
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
