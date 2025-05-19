package provider

import "sync"

// ProvideSyncMap 提供一个可注入的 *sync.Map 实例
func ProvideSyncMap() *sync.Map {
	return &sync.Map{}
}
