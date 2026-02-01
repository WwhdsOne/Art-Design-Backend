package utils

import (
	"time"

	"github.com/sony/sonyflake"
)

var sf = sonyflake.NewSonyflake(sonyflake.Settings{
	StartTime: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
})

func GenerateSnowflakeID() int64 {
	id, _ := sf.NextID()
	return int64(id)
}
