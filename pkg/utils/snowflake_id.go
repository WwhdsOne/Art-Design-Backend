package utils

import (
	"github.com/sony/sonyflake"
	"time"
)

var sf = sonyflake.NewSonyflake(sonyflake.Settings{
	StartTime: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
})

func GenerateSnowflakeId() int64 {
	id, _ := sf.NextID()
	return int64(id)
}
