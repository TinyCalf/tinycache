package tinycache

import (
	"./lru"
)

type tinyString string

//Len 获取tinyString的长度
func (s tinyString) Len() int {
	return len(s)
}

var _ lru.Value = (*tinyString)(nil)