package utils

import (
	"strconv"
	"strings"
	"time"
)

// ParseDuration 将字符串解析为 time.Duration 类型，支持标准格式和自定义扩展。
// 支持的格式包括：
//   - 标准 time.ParseDuration 能识别的所有格式（如 "1h", "30m", "5s" 等）
//   - 扩展了对天数 "d" 的支持（例如 "3d", "1d2h", "7d12h30m" 等）
//
// 如果输入无效或无法解析，则返回最接近的时间值，或 0。
func ParseDuration(d string) time.Duration {
	// 去除前后空格
	d = strings.TrimSpace(d)

	// 首先尝试使用标准库解析
	dr, _ := time.ParseDuration(d)
	if dr != 0 {
		return dr
	}

	// 如果包含 "d"，表示有天数部分
	if strings.Contains(d, "d") {
		before, after, _ := strings.Cut(d, "d")
		// 解析 "d" 前面的天数
		hour, _ := strconv.Atoi(before)
		dr = time.Hour * 24 * time.Duration(hour)

		// 解析 "d" 后面的部分作为剩余时间
		ndr, _ := time.ParseDuration(after)
		return dr + ndr
	}

	// 最后尝试将整个字符串视为整数秒（默认单位为秒）
	dv, _ := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv) * time.Second
}
