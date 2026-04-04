package service

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/ai"
	"Art-Design-Backend/pkg/constant/llmid"
	"Art-Design-Backend/pkg/constant/prompt"
	"Art-Design-Backend/pkg/constant/scheduler"
	"Art-Design-Backend/pkg/ws"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

type BrowserAgentService struct {
	BrowserAgentRepo *repository.BrowserAgentRepo
	AIModelRepo      *repository.AIModelRepo
	AIProviderRepo   *repository.AIProviderRepo
	AIModelClient    *ai.AIModelClient
	GormTX           *db.GormTransactionManager

	// LLM 配置内存缓存（每 10 分钟刷新一次）
	llmCacheMu     sync.Mutex
	llmProvider    *entity.AIProvider
	llmModel       *entity.AIModel
	llmCacheExpire time.Time
}

func NewBrowserAgentService(
	browserAgentRepo *repository.BrowserAgentRepo,
	aiModelRepo *repository.AIModelRepo,
	aiProviderRepo *repository.AIProviderRepo,
	aiModelClient *ai.AIModelClient,
	gormTX *db.GormTransactionManager,
) *BrowserAgentService {
	c := cron.New()
	b := &BrowserAgentService{
		BrowserAgentRepo: browserAgentRepo,
		AIModelRepo:      aiModelRepo,
		AIProviderRepo:   aiProviderRepo,
		AIModelClient:    aiModelClient,
		GormTX:           gormTX,
	}
	_ = c.AddFunc(scheduler.BrowserAgentStaleActionCron, func() {
		if err := b.BrowserAgentRepo.
			MarkStaleActionsFailed(context.Background(), scheduler.BrowserAgentActionMaxDuration); err != nil {
			zap.L().Error("更新旧动作状态失败", zap.Error(err))
		}
		if err := b.BrowserAgentRepo.
			MarkStaleAndFailedMessages(context.Background(), scheduler.BrowserAgentMessageMaxDuration); err != nil {
			zap.L().Error("更新旧任务状态为失败", zap.Error(err))
		}
		zap.L().Info("未完成旧任务状态已更新为失败")
	})
	c.Start()
	return b
}

// =========================
// 1. Conversation CRUD
// =========================

func (s *BrowserAgentService) CreateConversation(c *gin.Context, req *request.CreateConversationRequest) (*response.ConversationResponse, error) {
	conv := &entity.BrowserAgentConversation{
		Title:       req.Title,
		BrowserType: req.BrowserType,
	}
	if err := s.BrowserAgentRepo.CreateConversation(c, conv); err != nil {
		return nil, err
	}
	var convResp response.ConversationResponse
	_ = copier.Copy(&convResp, conv)
	return &convResp, nil
}

func (s *BrowserAgentService) GetConversationByID(c *gin.Context, id int64) (*response.ConversationResponse, error) {
	conv, err := s.BrowserAgentRepo.GetConversationByID(c, id)
	if err != nil {
		return nil, err
	}
	var convResp response.ConversationResponse
	_ = copier.Copy(&convResp, conv)
	return &convResp, nil
}

func (s *BrowserAgentService) ListConversations(c *gin.Context, userID int64, queryParam *query.BrowserAgentConversation) (*common.PaginationResp[response.ConversationResponse], error) {
	conversations, total, err := s.BrowserAgentRepo.ListConversationsByUserID(c, userID, queryParam)
	if err != nil {
		zap.L().Error("查询浏览器代理会话列表失败", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}

	responses := make([]response.ConversationResponse, len(conversations))
	for i := range conversations {
		_ = copier.Copy(&responses[i], &conversations[i])
	}

	return common.BuildPageResp[response.ConversationResponse](responses, total, queryParam.PaginationReq), nil
}

func (s *BrowserAgentService) RenameConversation(c *gin.Context, req *request.RenameConversationRequest) error {
	conv, err := s.BrowserAgentRepo.GetConversationByID(c, int64(req.ID))
	if err != nil {
		return err
	}
	conv.Title = req.Title
	return s.BrowserAgentRepo.UpdateConversation(c, conv)
}

func (s *BrowserAgentService) DeleteConversation(c *gin.Context, conversationID int64) error {
	err := s.GormTX.Transaction(c, func(ctx context.Context) (err error) {
		messageIDList, err := s.BrowserAgentRepo.ListMessagesIDListByConversationID(ctx, conversationID)
		if err != nil {
			return
		}
		if err = s.BrowserAgentRepo.DeleteActionsByMessageIDList(ctx, messageIDList); err != nil {
			return
		}
		if err = s.BrowserAgentRepo.DeleteMessagesByConversationID(ctx, conversationID); err != nil {
			return
		}
		if err = s.BrowserAgentRepo.DeleteConversation(ctx, conversationID); err != nil {
			return
		}
		return
	})
	return err
}

// =========================
// 2. Message CRUD
// =========================

func (s *BrowserAgentService) CreateMessage(c *gin.Context, req *request.CreateMessageRequest) (*response.MessageResponse, error) {
	msg := &entity.BrowserAgentMessage{
		ConversationID: req.ConversationID,
		Content:        req.Content,
	}
	if err := s.BrowserAgentRepo.CreateMessage(c, msg); err != nil {
		return nil, err
	}
	var msgResp response.MessageResponse
	_ = copier.Copy(&msgResp, msg)
	return &msgResp, nil
}

func (s *BrowserAgentService) ListMessages(c *gin.Context, req *request.GetMessagesRequest) ([]response.MessageResponse, error) {
	messages, err := s.BrowserAgentRepo.ListMessagesByConversationID(c, req.ConversationID)
	if err != nil {
		return nil, err
	}

	responses := make([]response.MessageResponse, len(messages))
	for i := range messages {
		_ = copier.Copy(&responses[i], &messages[i])
	}

	return responses, nil
}

// =========================
// 3. Action CRUD
// =========================

func (s *BrowserAgentService) ListActions(c *gin.Context, req *request.GetActionsRequest) ([]response.ActionResponse, error) {
	actions, err := s.BrowserAgentRepo.ListActionsByMessageID(c, req.MessageID)
	if err != nil {
		return nil, err
	}

	responses := make([]response.ActionResponse, len(actions))

	for i := range actions {
		_ = copier.Copy(&responses[i], &actions[i])
	}

	return responses, nil
}

// =========================
// 4. 任务处理
// =========================

func (s *BrowserAgentService) HandleTask(c context.Context, messageID int64, pageState *ws.PageState) ([]*ws.Action, error) {
	msg, err := s.BrowserAgentRepo.GetMessageByID(c, messageID)
	if err != nil {
		return nil, err
	}

	zap.L().Info("========== 收到新任务 ==========",
		zap.Int64("messageID", messageID),
		zap.Int64("conversationID", msg.ConversationID),
		zap.String("task", msg.Content),
	)

	if pageState == nil {
		return nil, errors.New("页面状态为空")
	}

	elementsCount := len(pageState.Elements)

	zap.L().Info("页面状态",
		zap.String("url", pageState.URL),
		zap.String("title", pageState.Title),
		zap.Int("elementsCount", elementsCount),
	)

	elements := make([]string, elementsCount)

	for i, elem := range pageState.Elements {
		elements[i] = fmt.Sprintf("%d.[%s]%s", i+1, elem.Tag, elem.Text)
	}

	zap.L().Info("可交互元素", zap.Strings("elements", elements))

	actions, err := s.decideAction(c, msg.Content, pageState)
	if err != nil {
		return nil, err
	}

	// 所有动作存入数据库
	for _, action := range actions {
		dbAction := s.wsActionToEntity(messageID, action)
		if err = s.BrowserAgentRepo.CreateAction(c, dbAction); err != nil {
			return nil, err
		}
		action.ActionID = dbAction.ID
	}

	s.logAction("首次决策", actions[0])
	if len(actions) > 1 {
		zap.L().Info("多动作模式", zap.Int("count", len(actions)))
	}

	return actions, nil
}

func (s *BrowserAgentService) HandleResult(c context.Context, msg *ws.ClientMessage) ([]*ws.Action, bool, error) {
	zap.L().Info("========== 收到执行结果 ==========",
		zap.Int64("actionID", msg.ActionID),
		zap.Int64("messageID", msg.MessageID),
		zap.Bool("success", msg.Success),
		zap.Int("executionTime(ms)", msg.ExecutionTime),
	)

	var errPtr *string
	if msg.Error != "" {
		errPtr = new(msg.Error)
	}
	execTimePtr := new(msg.ExecutionTime)

	if !msg.Success {
		zap.L().Error("操作执行失败",
			zap.Int64("actionID", msg.ActionID),
			zap.String("error", msg.Error),
		)
		if err := s.GormTX.Transaction(c, func(ctx context.Context) (err error) {
			if err = s.BrowserAgentRepo.UpdateActionStatus(ctx, msg.ActionID, entity.ActionStatusFailed, errPtr, execTimePtr); err != nil {
				return
			}
			if err = s.BrowserAgentRepo.UpdateMessageState(ctx, msg.MessageID, entity.MessageStateError); err != nil {
				return
			}
			return
		}); err != nil {
			return nil, false, err
		}

		return nil, false, errors.New(msg.Error)
	}

	if err := s.BrowserAgentRepo.UpdateActionStatus(c, msg.ActionID, entity.ActionStatusSuccess, nil, execTimePtr); err != nil {
		return nil, false, err
	}

	action, err := s.BrowserAgentRepo.GetActionByID(c, msg.ActionID)
	if err != nil {
		return nil, false, err
	}

	pageState := msg.PageState
	if pageState != nil {
		zap.L().Info("当前页面状态",
			zap.String("url", pageState.URL),
			zap.String("title", pageState.Title),
			zap.Int("elementsCount", len(pageState.Elements)),
		)
	}

	nextActions, finished, err := s.decideNextAction(c, pageState, msg.Task, action.MessageID)
	if err != nil {
		return nil, false, err
	}

	if finished {
		zap.L().Info("任务完成", zap.Int64("messageID", msg.MessageID))
		if err = s.BrowserAgentRepo.UpdateMessageState(c, msg.MessageID, entity.MessageStateFinished); err != nil {
			return nil, false, err
		}
		return nil, true, nil
	}

	// 所有动作存入数据库
	for _, nextAction := range nextActions {
		dbAction := s.wsActionToEntity(action.MessageID, nextAction)
		if err = s.BrowserAgentRepo.CreateAction(c, dbAction); err != nil {
			return nil, false, err
		}
		nextAction.ActionID = dbAction.ID
	}

	s.logAction("下一步决策", nextActions[0])
	if len(nextActions) > 1 {
		zap.L().Info("多动作模式", zap.Int("count", len(nextActions)))
	}

	return nextActions, false, nil
}

func (s *BrowserAgentService) wsActionToEntity(messageID int64, action *ws.Action) *entity.BrowserAgentAction {
	return &entity.BrowserAgentAction{
		MessageID:    messageID,
		ActionType:   action.Action,
		ElementIndex: action.Index,
		Status:       entity.ActionStatusPending,
		URL:          action.URL,
		Selector:     action.Selector,
		Value:        action.Value,
		Distance:     action.Distance,
		Timeout:      action.Timeout,
	}
}

func (s *BrowserAgentService) logAction(stage string, action *ws.Action) {
	fields := []zap.Field{
		zap.String("stage", stage),
		zap.Int64("actionID", action.ActionID),
		zap.String("action", action.Action),
	}

	if action.Index != nil {
		fields = append(fields, zap.Int("index", *action.Index))
	}

	switch action.Action {
	case "goto":
		if action.URL != nil {
			fields = append(fields, zap.String("url", *action.URL))
		}
	case "click":
		if action.Selector != nil {
			fields = append(fields, zap.String("selector", *action.Selector))
		}
	case "input", "select":
		if action.Selector != nil {
			fields = append(fields, zap.String("selector", *action.Selector))
		}
		if action.Value != nil {
			fields = append(fields, zap.String("value", *action.Value))
		}
	case "scroll":
		if action.Distance != nil {
			fields = append(fields, zap.Int("distance", *action.Distance))
		}
	case "wait":
		if action.Timeout != nil {
			fields = append(fields, zap.Int("timeout(ms)", *action.Timeout))
		}
	case "finish_task":
		fields = append(fields, zap.String("message", "任务完成"))
	}

	zap.L().Info("========== 发送操作指令 ==========", fields...)
}

// =========================
// 5. 大模型相关
// =========================

func (s *BrowserAgentService) callLLM(
	c context.Context,
	systemPrompt,
	promptText string,
) (string, error) {

	// 获取 LLM 配置（内存缓存 10 分钟，过期后重新查 Redis）
	provider, modelInfo, err := s.getLLMConfig(c)
	if err != nil {
		return "", err
	}

	const maxRetries = 3
	var respJSON []byte
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		respJSON, lastErr = s.AIModelClient.ChatRequest(
			c,
			provider.BaseURL+modelInfo.APIPath,
			provider.APIKey,
			ai.DefaultChatRequest(
				modelInfo.Model,
				[]ai.ChatMessage{
					{Role: "system", Content: systemPrompt},
					{Role: "user", Content: promptText},
				},
			),
		)
		if lastErr == nil {
			break
		}
		zap.L().Warn("调用LLM失败，准备重试",
			zap.Int("attempt", attempt),
			zap.Int("maxRetries", maxRetries),
			zap.Error(lastErr),
		)
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}
	if lastErr != nil {
		zap.L().Error("调用LLM失败（已重试"+strconv.Itoa(maxRetries)+"次）",
			zap.String("promptText", promptText), zap.Error(lastErr))
		return "", fmt.Errorf("调用LLM失败（已重试%d次）: %w", maxRetries, lastErr)
	}

	var browserResp ai.ChatCompletionResponse
	if err := sonic.Unmarshal(respJSON, &browserResp); err != nil {
		zap.L().Error("解析 LLM 原始响应失败", zap.Error(err))
		return "", fmt.Errorf("解析 LLM 原始响应失败: %w", err)
	}

	rawContent := strings.TrimSpace(browserResp.FirstText())
	if rawContent == "" {
		return "", errors.New("LLM 返回内容为空")
	}

	cleanJSON, err := ai.ExtractJSONFromLLMOutput(rawContent)
	if err != nil {
		zap.L().Error(
			"无法从 LLM 输出中提取 JSON",
			zap.String("raw", rawContent),
			zap.Error(err),
		)
		return "", err
	}

	zap.L().Debug(
		"LLM JSON 提取成功",
		zap.String("promptText", promptText),
		zap.String("json", cleanJSON),
	)

	return cleanJSON, nil
}

// getLLMConfig 获取 LLM 供应商和模型配置，内存缓存 10 分钟
func (s *BrowserAgentService) getLLMConfig(c context.Context) (*entity.AIProvider, *entity.AIModel, error) {
	s.llmCacheMu.Lock()
	defer s.llmCacheMu.Unlock()

	now := time.Now()
	if s.llmProvider != nil && s.llmModel != nil && now.Before(s.llmCacheExpire) {
		return s.llmProvider, s.llmModel, nil
	}

	// 缓存过期或首次调用，重新查询
	provider, err := s.AIProviderRepo.GetAIProviderByIDWithCache(c, llmid.BrowserProviderID)
	if err != nil {
		return nil, nil, fmt.Errorf("获取浏览器模型供应商失败: %w", err)
	}

	model, err := s.AIModelRepo.GetAIModelByIDWithCache(c, llmid.BrowserModelID)
	if err != nil {
		return nil, nil, fmt.Errorf("获取浏览器模型失败: %w", err)
	}

	s.llmProvider = provider
	s.llmModel = model
	s.llmCacheExpire = now.Add(scheduler.BrowserAgentLLMCacheTTL)

	zap.L().Debug("LLM 配置缓存已刷新",
		zap.String("provider", provider.Name),
		zap.String("model", model.Model),
		zap.Time("nextRefresh", s.llmCacheExpire),
	)

	return provider, model, nil
}

func (s *BrowserAgentService) decideAction(
	c context.Context,
	task string,
	pageState *ws.PageState,
) ([]*ws.Action, error) {

	zap.L().Info(
		"开始智能任务处理",
		zap.String("task", task),
		zap.Any("pageState", pageState),
	)

	decidePrompt := s.buildPrompt(task, pageState)

	resp, err := s.callLLM(c, prompt.BrowserSystemPrompt, decidePrompt)
	if err != nil {
		return nil, err
	}

	actions, err := s.parseAgentOutput(resp)
	if err != nil {
		return nil, err
	}

	for _, action := range actions {
		if err = s.validateAction(action); err != nil {
			return nil, err
		}
	}

	return actions, nil
}

func (s *BrowserAgentService) decideNextAction(
	c context.Context,
	pageState *ws.PageState,
	task string,
	messageID int64,
) ([]*ws.Action, bool, error) {

	// 获取已完成的操作历史，帮助 AI 判断任务进度
	completedActions := s.buildCompletedActionsSummary(c, messageID)
	nextActionPrompt := s.buildNextPrompt(pageState, task, completedActions)

	resp, err := s.callLLM(c, prompt.BrowserSystemPrompt, nextActionPrompt)
	if err != nil {
		return nil, false, err
	}

	if strings.Contains(resp, `"action":"finish_task"`) ||
		strings.Contains(resp, `"action": "finish_task"`) {
		return nil, true, nil
	}

	actions, err := s.parseAgentOutput(resp)
	if err != nil {
		return nil, false, err
	}

	for _, action := range actions {
		if err = s.validateAction(action); err != nil {
			return nil, false, err
		}
	}

	return actions, false, nil
}

// parseAgentOutput 解析 LLM 返回的结构化输出
// 支持单动作和多动作格式，返回动作列表
func (s *BrowserAgentService) parseAgentOutput(resp string) ([]*ws.Action, error) {
	var output ws.AgentOutput
	zap.L().Debug("LLM返回结果", zap.String("response", resp))
	if err := sonic.Unmarshal([]byte(resp), &output); err != nil {
		// 降级：尝试直接作为 Action 解析（兼容旧格式）
		var action ws.Action
		if err2 := sonic.Unmarshal([]byte(resp), &action); err2 == nil && action.Action != "" {
			return []*ws.Action{&action}, nil
		}
		return nil, fmt.Errorf("解析 AgentOutput 失败: %w", err)
	}

	// 记录结构化思维
	if output.Thinking != "" {
		zap.L().Debug("思维链", zap.String("thinking", output.Thinking))
	}
	if output.Evaluation != "" {
		zap.L().Debug("上步评估", zap.String("evaluation", output.Evaluation))
	}
	if output.Memory != "" {
		zap.L().Debug("任务记忆", zap.String("memory", output.Memory))
	}
	if output.NextGoal != "" {
		zap.L().Debug("下一步目标", zap.String("next_goal", output.NextGoal))
	}

	// 优先处理多动作
	if len(output.Actions) > 0 {
		actions := make([]*ws.Action, 0, len(output.Actions))
		for _, a := range output.Actions {
			action := &ws.Action{
				Action: a.Action,
				Index:  a.Index,
			}
			if a.Selector != nil {
				action.Selector = a.Selector
			}
			if a.Value != nil {
				action.Value = a.Value
			}
			if a.URL != nil {
				action.URL = a.URL
			}
			if a.Distance != nil {
				action.Distance = a.Distance
			}
			if a.Timeout != nil {
				action.Timeout = a.Timeout
			}
			actions = append(actions, action)
		}
		return actions, nil
	}

	// 单动作模式
	if output.Action != "" {
		action := &ws.Action{
			Action: output.Action,
		}

		// 从原始 JSON 提取 index、selector、value 等字段
		var rawMap map[string]any
		if err := sonic.Unmarshal([]byte(resp), &rawMap); err == nil {
			if idx, ok := rawMap["index"]; ok {
				if idxFloat, ok := idx.(float64); ok {
					action.Index = new(int(idxFloat))
				}
			}
			if sel, ok := rawMap["selector"]; ok {
				if selStr, ok := sel.(string); ok {
					action.Selector = new(selStr)
				}
			}
			if val, ok := rawMap["value"]; ok {
				if valStr, ok := val.(string); ok {
					action.Value = new(valStr)
				}
			}
			if url, ok := rawMap["url"]; ok {
				if urlStr, ok := url.(string); ok {
					action.URL = new(urlStr)
				}
			}
			if dist, ok := rawMap["distance"]; ok {
				if distFloat, ok := dist.(float64); ok {
					action.Distance = new(int(distFloat))
				}
			}
			if timeout, ok := rawMap["timeout"]; ok {
				if timeoutFloat, ok := timeout.(float64); ok {
					action.Timeout = new(int(timeoutFloat))
				}
			}
		}

		return []*ws.Action{action}, nil
	}

	return nil, errors.New("LLM 输出中未找到有效动作")
}

// =========================
// 6. Prompt 构建
// =========================

func (s *BrowserAgentService) buildPrompt(
	task string,
	pageState *ws.PageState,
) string {

	return "【用户目标】\n" +
		task + "\n\n" +
		s.buildPageStateSection(pageState)
}

// buildCompletedActionsSummary 生成已执行操作的摘要，供 AI 评估任务进度
func (s *BrowserAgentService) buildCompletedActionsSummary(c context.Context, messageID int64) []string {
	actions, err := s.BrowserAgentRepo.ListActionsByMessageID(c, messageID)
	if err != nil {
		return nil
	}

	var summary []string
	step := 0
	for _, a := range actions {
		// 跳过未执行的操作（多动作批次中排在后面的）
		if a.Status == entity.ActionStatusPending || a.Status == entity.ActionStatusRunning {
			continue
		}
		step++
		desc := fmt.Sprintf("%d. [%s]", step, a.ActionType)
		switch a.ActionType {
		case "goto":
			if a.URL != nil {
				desc += fmt.Sprintf(" %s", *a.URL)
			}
		case "click":
			if a.ElementIndex != nil {
				desc += fmt.Sprintf(" 元素[%d]", *a.ElementIndex)
			} else if a.Selector != nil {
				desc += fmt.Sprintf(" %s", *a.Selector)
			}
		case "input":
			if a.ElementIndex != nil {
				desc += fmt.Sprintf(" 元素[%d]", *a.ElementIndex)
			} else if a.Selector != nil {
				desc += fmt.Sprintf(" %s", *a.Selector)
			}
			if a.Value != nil {
				desc += fmt.Sprintf(": \"%s\"", *a.Value)
			}
		case "select":
			if a.ElementIndex != nil {
				desc += fmt.Sprintf(" 元素[%d]", *a.ElementIndex)
			} else if a.Selector != nil {
				desc += fmt.Sprintf(" %s", *a.Selector)
			}
			if a.Value != nil {
				desc += fmt.Sprintf(" → \"%s\"", *a.Value)
			}
		case "scroll":
			if a.Distance != nil {
				desc += fmt.Sprintf(" 距离%d", *a.Distance)
			}
		case "wait":
			if a.Timeout != nil {
				desc += fmt.Sprintf(" %dms", *a.Timeout)
			}
		}
		// 标记失败的操作
		if a.Status == entity.ActionStatusFailed {
			desc += " (失败)"
		}
		summary = append(summary, desc)
	}
	return summary
}

func (s *BrowserAgentService) buildNextPrompt(
	pageState *ws.PageState,
	task string,
	completedActions []string,
) string {
	var sb strings.Builder
	sb.WriteString("【任务进度评估】\n\n")
	sb.WriteString("原始任务: " + task + "\n\n")

	if len(completedActions) > 0 {
		sb.WriteString(fmt.Sprintf("已完成操作（共 %d 步）:\n", len(completedActions)))
		for _, a := range completedActions {
			sb.WriteString(a + "\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString(s.buildPageStateSection(pageState))

	sb.WriteString("\n请先判断原始任务是否已完成：\n")
	sb.WriteString("- 对比「原始任务」和「已完成操作」，确认所有步骤是否已执行\n")
	sb.WriteString("- 如果任务目标已达成，返回 finish_task\n")
	sb.WriteString("- 如果还有未完成步骤，返回下一步操作\n")

	return sb.String()
}

// maxElementsInPrompt 限制发送给 LLM 的最大元素数量，超出时截断以减少 token 消耗
const maxElementsInPrompt = 60

// truncateElementText 截断客户端传来的元素索引文本，保留前 max 行
func truncateElementText(text string, maxLines int) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= maxLines {
		return text
	}
	truncated := strings.Join(lines[:maxLines], "\n")
	truncated += fmt.Sprintf("\n\n... 已省略 %d 个元素（共 %d 个）", len(lines)-maxLines, len(lines))
	return truncated
}

func (s *BrowserAgentService) buildPageStateSection(pageState *ws.PageState) string {
	if pageState == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("【页面状态】\n")
	sb.WriteString("URL: " + pageState.URL + "\n")
	if pageState.Title != "" {
		sb.WriteString("标题: " + pageState.Title + "\n")
	}

	if pageState.ScrollInfo != nil {
		si := pageState.ScrollInfo
		sb.WriteString(fmt.Sprintf(
			"滚动: scrollTop=%.0f / clientHeight=%.0f / scrollHeight=%.0f | hasMoreAbove=%t, hasMoreBelow=%t\n",
			si.ScrollTop,
			si.ClientHeight,
			si.ScrollHeight,
			si.HasMoreAbove,
			si.HasMoreBelow,
		))
	}

	// 优先使用客户端生成的编号索引文本
	if pageState.ElementText != "" {
		sb.WriteString("\n【可交互元素】\n")
		sb.WriteString(truncateElementText(pageState.ElementText, maxElementsInPrompt))
	} else {
		// 降级：使用旧格式
		sb.WriteString("\n【可交互元素】\n")
		elements := pageState.Elements
		if len(elements) > maxElementsInPrompt {
			zap.L().Warn("元素数量超出上限，已截断",
				zap.Int("total", len(elements)),
				zap.Int("kept", maxElementsInPrompt),
			)
			elements = elements[:maxElementsInPrompt]
		}
		for i, elem := range elements {
			sb.WriteString(fmt.Sprintf("%d. ", i))
			sb.WriteString(s.formatElementText(&elem))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (s *BrowserAgentService) formatElementText(elem *ws.PageElement) string {
	tag := elem.Tag
	attrs := []string{}

	if elem.Type != nil && *elem.Type != tag {
		attrs = append(attrs, fmt.Sprintf(`type="%s"`, *elem.Type))
	}
	if elem.Role != nil {
		attrs = append(attrs, fmt.Sprintf(`role="%s"`, *elem.Role))
	}
	if elem.AriaLabel != nil {
		attrs = append(attrs, fmt.Sprintf(`aria-label="%s"`, *elem.AriaLabel))
	}
	if elem.AriaExpanded != nil {
		attrs = append(attrs, fmt.Sprintf(`aria-expanded="%s"`, *elem.AriaExpanded))
	}
	if elem.AriaChecked != nil {
		attrs = append(attrs, fmt.Sprintf(`aria-checked="%s"`, *elem.AriaChecked))
	}
	if elem.Required != nil && *elem.Required {
		attrs = append(attrs, "required")
	}
	if elem.Disabled != nil && *elem.Disabled {
		attrs = append(attrs, "disabled")
	}

	attrStr := ""
	if len(attrs) > 0 {
		attrStr = " " + strings.Join(attrs, " ")
	}

	text := elem.Text
	if elem.Label != nil && *elem.Label != "" && *elem.Label != text {
		text = *elem.Label
	}
	if elem.Value != nil && *elem.Value != "" && *elem.Type != "password" {
		val := *elem.Value
		if len(val) > 50 {
			val = val[:50] + "..."
		}
		return fmt.Sprintf(`<%s%s> %s (当前值: "%s")`, tag, attrStr, text, val)
	}

	return fmt.Sprintf(`<%s%s> %s`, tag, attrStr, text)
}

func (s *BrowserAgentService) buildHistorySection(history []*entity.BrowserAgentMessage) string {
	if len(history) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("【历史任务】\n")

	for i, h := range history {
		sb.WriteString(fmt.Sprintf(
			"%d. %s\n",
			i+1,
			h.Content,
		))
	}

	return sb.String()
}

// =========================
// 7. 解析与校验
// =========================

func (s *BrowserAgentService) parseAction(resp string) (*ws.Action, error) {
	var action ws.Action
	zap.L().Debug("LLM返回结果", zap.String("response", resp))
	if err := sonic.Unmarshal([]byte(resp), &action); err != nil {
		return nil, fmt.Errorf("解析 Action 失败: %w", err)
	}
	return &action, nil
}

func (s *BrowserAgentService) validateAction(action *ws.Action) error {
	validActions := map[string]bool{
		"goto":        true,
		"click":       true,
		"input":       true,
		"select":      true,
		"scroll":      true,
		"wait":        true,
		"finish_task": true,
	}

	if !validActions[action.Action] {
		return fmt.Errorf("非法 Action: %s", action.Action)
	}

	switch action.Action {
	case "goto":
		if action.URL == nil || *action.URL == "" {
			return errors.New("goto 缺少 url")
		}
	case "click":
		// 支持 index 或 selector 两种定位方式
		if action.Index == nil && (action.Selector == nil || *action.Selector == "") {
			return errors.New("click 缺少 index 或 selector")
		}
	case "input", "select":
		if action.Index == nil && (action.Selector == nil || *action.Selector == "") {
			return fmt.Errorf("%s 缺少 index 或 selector", action.Action)
		}
		if action.Value == nil {
			return fmt.Errorf("%s 缺少 value", action.Action)
		}
	case "scroll":
		if action.Distance == nil {
			return errors.New("scroll 缺少 distance")
		}
	case "wait":
		if action.Timeout == nil {
			return errors.New("wait 缺少 timeout")
		}
	}

	return nil
}
