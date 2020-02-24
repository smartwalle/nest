package nest

import "errors"

var (
	ErrNodeNotExist     = errors.New("节点不存在")
	ErrParentNotExist   = errors.New("父节点或者兄弟节点不存在")
	ErrParentNotAllowed = errors.New("不满足条件的父节点或者兄弟节点")
	ErrUnknownPosition  = errors.New("未知位置")
	ErrNodeExists       = errors.New("节点已经存在")
	ErrNotLeafNode      = errors.New("不是叶子节点")
	ErrUnknownStatus    = errors.New("未知状态")
)
