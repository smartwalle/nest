package nest

import "time"

const (
	K_STATUS_ENABLE  = 1000 // 启用
	K_STATUS_DISABLE = 2000 // 禁用
)

type BasicModel struct {
	Id         int64      `json:"id"                        sql:"id"`
	Type       int        `json:"type"                      sql:"type"`
	Name       string     `json:"name"                      sql:"name"`
	LeftValue  int        `json:"left_value"                sql:"left_value"`
	RightValue int        `json:"right_value"               sql:"right_value"`
	Depth      int        `json:"depth"                     sql:"depth"`
	Status     int        `json:"status"                    sql:"status"`
	CreatedOn  *time.Time `json:"created_on,omitempty"      sql:"created_on"`
	UpdatedOn  *time.Time `json:"updated_on,omitempty"      sql:"updated_on"`
}

func (this *BasicModel) IsLeafNode() bool {
	return this.LeftValue+1 == this.RightValue
}
