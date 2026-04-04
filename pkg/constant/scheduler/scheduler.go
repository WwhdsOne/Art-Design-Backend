package scheduler

import "time"

const (
	// BrowserAgentStaleActionCron 更新长时间未处理的任务状态
	BrowserAgentStaleActionCron = "0 0 * * * *" // 秒 分 时 日 月 周 年(可选，部分 cron 库支持秒)

	// BrowserAgentActionMaxDuration 任务最大允许执行时间
	BrowserAgentActionMaxDuration = 10 * time.Minute

	// BrowserAgentMessageMaxDuration 消息最大允许执行时间
	BrowserAgentMessageMaxDuration = 60 * time.Minute

	// BrowserAgentLLMCacheTTL LLM 配置（Provider/Model）内存缓存刷新间隔
	BrowserAgentLLMCacheTTL = 10 * time.Minute
)
