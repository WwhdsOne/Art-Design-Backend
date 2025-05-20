package container

import "github.com/google/wire"

// Container 提供各种库的实例
var Container = wire.NewSet(
	NewSyncMap,
)
