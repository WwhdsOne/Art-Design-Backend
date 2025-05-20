package container

import "sync"

// NewSyncMap 提供一个可注入的 *sync.Map 实例
func NewSyncMap() *sync.Map {
	return &sync.Map{}
}
