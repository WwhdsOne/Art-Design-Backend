package utils

import (
	"github.com/sony/sonyflake"
	"time"
)

var sf = sonyflake.NewSonyflake(sonyflake.Settings{
	StartTime: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
})

func GenerateSnowflakeId() (int64, error) {
	id, err := sf.NextID()
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}
