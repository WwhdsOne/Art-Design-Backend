package llmid

const (
	// 内部记录，不对外开放
	embedProviderQwenID       = 51088793876300041
	multiModelQwenID          = 61331207874412809
	browserProviderZhipuID    = 81681722898382850
	browserModelZhipuID       = 81682220275728386
	browserProviderDeepSeekID = 50361636266975107
	browserModelDeepSeekID    = 50372307280994179

	// 视觉模型（Qwen 3.5 Flash via 通义千问）
	browserVisionProviderID = embedProviderQwenID // 复用通义千问供应商
	browserVisionModelID    = 91702166376415492                  // TODO: 需在数据库中创建 Qwen 3.5 Flash 模型记录后填入

	// 开发使用

	EmbedProviderID       = embedProviderQwenID
	MultiModelID          = multiModelQwenID
	BrowserProviderID     = browserProviderDeepSeekID
	BrowserModelID        = browserModelDeepSeekID
	BrowserVisionProvider = browserVisionProviderID
	BrowserVisionModel    = browserVisionModelID
)
