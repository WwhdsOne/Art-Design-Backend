package common

import (
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
)

// LongStringID 定义自定义类型，用于表示前端传递的长整形的字符串
type LongStringID int64

// UnmarshalJSON 实现 UnmarshalJSON 方法，用于将字符串反序列化为 int64
func (id *LongStringID) UnmarshalJSON(data []byte) error {
	// 尝试将 JSON 数据解析为字符串
	var strValue string
	if err := sonic.Unmarshal(data, &strValue); err != nil {
		return err
	}

	// 将字符串转换为 int64
	intValue, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return err
	}

	// 赋值给 LongStringID
	*id = LongStringID(intValue)
	return nil
}

// LongStringIDs 自定义类型，用于实现 IDs 的解析逻辑
type LongStringIDs []int64

// UnmarshalJSON 实现 UnmarshalJSON 方法
func (c *LongStringIDs) UnmarshalJSON(data []byte) error {
	// 定义一个临时变量，用于解析 JSON 中的字符串数组
	var temp []string
	if err := sonic.Unmarshal(data, &temp); err != nil {
		return err
	}

	// 将字符串数组转换为 int64 数组
	*c = make([]int64, len(temp))
	for i, idStr := range temp {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return fmt.Errorf("无法将 %s 转换为 int64: %v", idStr, err)
		}
		(*c)[i] = id
	}
	return nil
}

// MarshalJSON 实现 LongStringID 序列化为 JSON 字符串
func (id *LongStringID) MarshalJSON() ([]byte, error) {
	str := strconv.FormatInt(int64(*id), 10)
	return sonic.Marshal(str)
}

// MarshalJSON 实现 LongStringIDs 序列化为 JSON 字符串数组
func (c *LongStringIDs) MarshalJSON() ([]byte, error) {
	strIDs := make([]string, len(*c))
	for i, id := range *c {
		strIDs[i] = strconv.FormatInt(id, 10)
	}
	return sonic.Marshal(strIDs)
}
