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
	"Art-Design-Backend/pkg/ws"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type BrowserAgentService struct {
	BrowserAgentRepo *repository.BrowserAgentRepo
	AIModelRepo      *repository.AIModelRepo
	AIProviderRepo   *repository.AIProviderRepo
	AIModelClient    *ai.AIModelClient
	GormTX           *db.GormTransactionManager
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

func (s *BrowserAgentService) HandleTask(c context.Context, messageID int64, pageState *ws.PageState) (*ws.Action, error) {
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

	if elementsCount == 0 {
		return nil, errors.New("页面元素为空")
	}

	elements := make([]string, elementsCount)

	for i, elem := range pageState.Elements {
		elements[i] = fmt.Sprintf("%d.[%s]%s", i+1, elem.Tag, elem.Text)
	}

	zap.L().Info("可交互元素", zap.Strings("elements", elements))

	action, err := s.decideAction(c, msg.Content, pageState)
	if err != nil {
		return nil, err
	}

	dbAction := s.wsActionToEntity(messageID, action)
	if err = s.BrowserAgentRepo.CreateAction(c, dbAction); err != nil {
		return nil, err
	}

	action.ActionID = dbAction.ID

	s.logAction("首次决策", action)

	return action, nil
}

func (s *BrowserAgentService) HandleResult(c context.Context, msg *ws.ClientMessage) (*ws.Action, bool, error) {
	zap.L().Info("========== 收到执行结果 ==========",
		zap.Int64("actionID", msg.ActionID),
		zap.Int64("messageID", msg.MessageID),
		zap.Bool("success", msg.Success),
		zap.Int("executionTime(ms)", msg.ExecutionTime),
	)

	var errPtr *string
	if msg.Error != "" {
		errPtr = &msg.Error
	}
	execTimePtr := &msg.ExecutionTime

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

	nextAction, finished, err := s.decideNextAction(c, pageState, msg.Task)
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

	dbAction := s.wsActionToEntity(action.MessageID, nextAction)
	if err = s.BrowserAgentRepo.CreateAction(c, dbAction); err != nil {
		return nil, false, err
	}

	nextAction.ActionID = dbAction.ID

	s.logAction("下一步决策", nextAction)

	return nextAction, false, nil
}

func (s *BrowserAgentService) wsActionToEntity(messageID int64, action *ws.Action) *entity.BrowserAgentAction {
	return &entity.BrowserAgentAction{
		MessageID:  messageID,
		ActionType: action.Action,
		Status:     entity.ActionStatusPending,
		URL:        action.URL,
		Selector:   action.Selector,
		Value:      action.Value,
		Distance:   action.Distance,
		Timeout:    action.Timeout,
	}
}

func (s *BrowserAgentService) logAction(stage string, action *ws.Action) {
	fields := []zap.Field{
		zap.String("stage", stage),
		zap.Int64("actionID", action.ActionID),
		zap.String("action", action.Action),
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
	case "close_browser":
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

	provider, err := s.AIProviderRepo.GetAIProviderByIDWithCache(c, llmid.BrowserProviderID)
	if err != nil {
		zap.L().Error("获取浏览器智谱模型供应商失败", zap.Error(err))
		return "", fmt.Errorf("获取浏览器智谱模型供应商失败: %w", err)
	}

	modelInfo, err := s.AIModelRepo.GetAIModelByIDWithCache(c, llmid.BrowserModelID)
	if err != nil {
		zap.L().Error("获取浏览器智谱模型失败", zap.Error(err))
		return "", fmt.Errorf("获取浏览器智谱模型失败: %w", err)
	}

	respJSON, err := s.AIModelClient.ChatRequest(
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
	if err != nil {
		zap.L().Error("调用LLM失败", zap.String("promptText", promptText), zap.Error(err))
		return "", fmt.Errorf("调用LLM失败: %w", err)
	}

	var browserResp ai.ChatCompletionResponse
	if err = sonic.Unmarshal(respJSON, &browserResp); err != nil {
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

func (s *BrowserAgentService) decideAction(
	c context.Context,
	task string,
	pageState *ws.PageState,
) (*ws.Action, error) {

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

	action, err := s.parseAction(resp)
	if err != nil {
		return nil, err
	}

	return action, s.validateAction(action)
}

func (s *BrowserAgentService) decideNextAction(
	c context.Context,
	pageState *ws.PageState,
	task string,
) (*ws.Action, bool, error) {

	nextActionPrompt := s.buildNextPrompt(pageState, task)

	resp, err := s.callLLM(c, prompt.BrowserSystemPrompt, nextActionPrompt)
	if err != nil {
		return nil, false, err
	}

	if strings.Contains(resp, `"action":"close_browser"`) ||
		strings.Contains(resp, `"action": "close_browser"`) {
		return nil, true, nil
	}

	action, err := s.parseAction(resp)
	if err != nil {
		return nil, false, err
	}

	return action, false, s.validateAction(action)
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

func (s *BrowserAgentService) buildNextPrompt(
	pageState *ws.PageState,
	task string,
) string {

	return "【继续执行当前任务】\n\n" +
		"原始任务:" + task + "\n\n" +
		s.buildPageStateSection(pageState)
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
	sb.WriteString("可交互元素:\n")

	for i, elem := range pageState.Elements {

		var valueInfo string
		if elem.Value != nil && *elem.Value != "" {
			valueInfo = fmt.Sprintf(" | value=%s", *elem.Value)
		}

		sb.WriteString(fmt.Sprintf(
			"%d. <%s> %s | selector=%s%s\n",
			i+1,
			elem.Tag,
			elem.Text,
			elem.Selector,
			valueInfo,
		))
	}

	return sb.String()
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
		"goto":          true,
		"click":         true,
		"input":         true,
		"select":        true,
		"scroll":        true,
		"wait":          true,
		"close_browser": true,
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
		if action.Selector == nil || *action.Selector == "" {
			return errors.New("click 缺少 selector")
		}
	case "input", "select":
		if action.Selector == nil || action.Value == nil {
			return errors.New("input/select 缺少 selector 或 value")
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
