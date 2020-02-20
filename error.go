package nest

import "errors"

var (
	ErrNodeNotExist     = errors.New("节点不存在")
	ErrParentNotExist   = errors.New("父节点不存在")
	ErrParentNotAllowed = errors.New("不满足条件的父节点")
	ErrUnknownPosition  = errors.New("未知位置")
	ErrNodeExists       = errors.New("节点已经存在")
	ErrAlreadyFirstNode = errors.New("已经是第一个节点")
	ErrAlreadyLastNode  = errors.New("已经是最后一个节点")
)
