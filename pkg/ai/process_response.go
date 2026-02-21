package ai

import (
	"errors"
	"strings"

	"github.com/bytedance/sonic"
)

func ExtractJSONFromLLMOutput(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	// 情况 1：```json ... ```
	if after, ok := strings.CutPrefix(raw, "```"); ok {
		raw = after
		raw = strings.TrimPrefix(raw, "json")
		raw = strings.TrimSuffix(raw, "```")
		raw = strings.TrimSpace(raw)
	}

	// 情况 2：本身就是 JSON
	if sonic.Valid([]byte(raw)) {
		return raw, nil
	}

	// 情况 3：混合文本，尝试截取第一个 JSON 对象
	start := strings.IndexAny(raw, "{[")
	if start == -1 {
		return "", errors.New("未找到 JSON 起始符号")
	}

	// 尝试从 start 开始逐步缩短尾部，直到合法 JSON
	for i := len(raw); i > start; i-- {
		candidate := raw[start:i]
		if sonic.Valid([]byte(candidate)) {
			return candidate, nil
		}
	}

	return "", errors.New("无法从 LLM 输出中提取合法 JSON")
}
