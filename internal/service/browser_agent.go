package service

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
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
}

// =========================
// 1. Conversation CRUD
// =========================

func (s *BrowserAgentService) CreateConversation(c *gin.Context, req *request.CreateConversationRequest) (*response.ConversationResponse, error) {
	conv := &entity.BrowserAgentConversation{
		Title: req.Title,
		State: entity.ConversationStateRunning,
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

	responses := make([]response.ConversationResponse, 0, len(conversations))
	for _, conv := range conversations {
		var convResp response.ConversationResponse
		if err = copier.Copy(&convResp, conv); err != nil {
			zap.L().Error("拷贝浏览器代理会话属性失败", zap.Int64("conversationID", conv.ID), zap.Error(err))
			continue
		}
		responses = append(responses, convResp)
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

func (s *BrowserAgentService) DeleteConversation(c *gin.Context, id int64) error {
	return s.BrowserAgentRepo.DeleteConversation(c, id)
}

// =========================
// 2. Message CRUD
// =========================

func (s *BrowserAgentService) CreateMessage(c *gin.Context, req *request.CreateMessageRequest) (*response.MessageResponse, error) {

	msg := &entity.BrowserAgentMessage{
		ConversationID: req.ConversationID,
		Role:           req.Role,
		Content:        req.Content,
	}
	if err := s.BrowserAgentRepo.CreateMessage(c, msg); err != nil {
		return nil, err
	}
	var msgResp response.MessageResponse
	_ = copier.Copy(&msgResp, msg)
	return &msgResp, nil
}

func (s *BrowserAgentService) ListMessages(c *gin.Context, req *request.GetMessagesRequest) ([]*response.MessageResponse, error) {

	messages, err := s.BrowserAgentRepo.ListMessagesByConversationID(c, req.ConversationID)
	if err != nil {
		return nil, err
	}

	responses := make([]*response.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		var msgResp response.MessageResponse
		_ = copier.Copy(&msgResp, msg)
		responses = append(responses, &msgResp)
	}

	return responses, nil
}

// =========================
// 3. Action CRUD
// =========================

func (s *BrowserAgentService) CreateAction(c *gin.Context, req *request.CreateActionRequest) (*response.ActionResponse, error) {
	_, err := s.BrowserAgentRepo.GetMessageByID(c, req.MessageID)
	if err != nil {
		return nil, err
	}

	action := &entity.BrowserAgentAction{
		MessageID:  req.MessageID,
		ActionType: req.ActionType,
		Sequence:   req.Sequence,
		Status:     entity.ActionStatusPending,
		URL:        req.URL,
		Selector:   req.Selector,
		Value:      req.Value,
		Distance:   req.Distance,
		Timeout:    req.Timeout,
	}
	if err := s.BrowserAgentRepo.CreateAction(c, action); err != nil {
		return nil, err
	}
	var actionResp response.ActionResponse
	_ = copier.Copy(&actionResp, action)
	return &actionResp, nil
}

func (s *BrowserAgentService) ListActions(c *gin.Context, req *request.GetActionsRequest) ([]*response.ActionResponse, error) {
	actions, err := s.BrowserAgentRepo.ListActionsByMessageID(c, req.MessageID)
	if err != nil {
		return nil, err
	}

	responses := make([]*response.ActionResponse, 0, len(actions))
	for _, action := range actions {
		var actionResp response.ActionResponse
		_ = copier.Copy(&actionResp, action)
		responses = append(responses, &actionResp)
	}

	return responses, nil
}

// =========================
// 4. 任务处理
// =========================

func (s *BrowserAgentService) HandleTask(c context.Context, conversationID int64, task string, pageState *ws.PageState) (*ws.Action, error) {
	userMsg := &entity.BrowserAgentMessage{
		ConversationID: conversationID,
		Role:           entity.MessageRoleUser,
		Content:        task,
	}
	if err := s.BrowserAgentRepo.CreateMessage(c, userMsg); err != nil {
		return nil, err
	}

	history, _ := s.BrowserAgentRepo.GetRecentMessages(c, conversationID, 10)

	action, err := s.decideAction(c, task, pageState, history)
	if err != nil {
		return nil, err
	}

	assistantMsg := &entity.BrowserAgentMessage{
		ConversationID: conversationID,
		Role:           entity.MessageRoleAssistant,
		Content:        action.Action,
	}
	if err = s.BrowserAgentRepo.CreateMessage(c, assistantMsg); err != nil {
		return nil, err
	}

	dbAction := s.wsActionToEntity(assistantMsg.ID, action, 1)
	if err = s.BrowserAgentRepo.CreateAction(c, dbAction); err != nil {
		return nil, err
	}

	return action, nil
}

func (s *BrowserAgentService) HandleResult(c context.Context, conversationID int64, success bool, errMsg string, pageState *ws.PageState) (*ws.Action, bool, error) {
	if !success {
		_ = s.BrowserAgentRepo.UpdateConversationState(c, conversationID, entity.ConversationStateError)
		return nil, false, errors.New(errMsg)
	}

	history, _ := s.BrowserAgentRepo.GetRecentMessages(c, conversationID, 20)

	action, finished, err := s.decideNextAction(c, pageState, history)
	if err != nil {
		return nil, false, err
	}

	if finished {
		_ = s.BrowserAgentRepo.UpdateConversationState(c, conversationID, entity.ConversationStateFinished)
		return nil, true, nil
	}

	assistantMsg := &entity.BrowserAgentMessage{
		ConversationID: conversationID,
		Role:           entity.MessageRoleAssistant,
		Content:        action.Action,
	}
	if err := s.BrowserAgentRepo.CreateMessage(c, assistantMsg); err != nil {
		return nil, false, err
	}

	dbAction := s.wsActionToEntity(assistantMsg.ID, action, 1)
	if err := s.BrowserAgentRepo.CreateAction(c, dbAction); err != nil {
		return nil, false, err
	}

	return action, false, nil
}

func (s *BrowserAgentService) HandleResume(c context.Context, conversationID int64, pageState *ws.PageState) (*ws.Action, bool, error) {
	_ = s.BrowserAgentRepo.UpdateConversationState(c, conversationID, entity.ConversationStateRunning)

	history, _ := s.BrowserAgentRepo.GetRecentMessages(c, conversationID, 20)

	action, finished, err := s.decideNextAction(c, pageState, history)
	if err != nil {
		return nil, false, err
	}

	if finished {
		_ = s.BrowserAgentRepo.UpdateConversationState(c, conversationID, entity.ConversationStateFinished)
		return nil, true, nil
	}

	assistantMsg := &entity.BrowserAgentMessage{
		ConversationID: conversationID,
		Role:           entity.MessageRoleAssistant,
		Content:        action.Action,
	}
	if err := s.BrowserAgentRepo.CreateMessage(c, assistantMsg); err != nil {
		return nil, false, err
	}

	dbAction := s.wsActionToEntity(assistantMsg.ID, action, 1)
	if err := s.BrowserAgentRepo.CreateAction(c, dbAction); err != nil {
		return nil, false, err
	}

	return action, false, nil
}

func (s *BrowserAgentService) wsActionToEntity(messageID int64, action *ws.Action, sequence int) *entity.BrowserAgentAction {
	return &entity.BrowserAgentAction{
		MessageID:  messageID,
		ActionType: action.Action,
		Sequence:   sequence,
		Status:     entity.ActionStatusPending,
		URL:        action.URL,
		Selector:   action.Selector,
		Value:      action.Value,
		Distance:   action.Distance,
		Timeout:    action.Timeout,
	}
}

// =========================
// 5. 大模型相关
// =========================

func (s *BrowserAgentService) callLLM(
	c context.Context,
	systemPrompt,
	promptText string,
) (string, error) {

	provider, err := s.AIProviderRepo.GetAIProviderByIDWithCache(c, llmid.BrowserProviderZhipuID)
	if err != nil {
		zap.L().Error("获取浏览器智谱模型供应商失败", zap.Error(err))
		return "", fmt.Errorf("获取浏览器智谱模型供应商失败: %w", err)
	}

	modelInfo, err := s.AIModelRepo.GetAIModelByIDWithCache(c, llmid.BrowserModelZhipuID)
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
		zap.L().Error("调用智谱LLM失败", zap.String("promptText", promptText), zap.Error(err))
		return "", fmt.Errorf("调用智谱LLM失败: %w", err)
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

	// ⭐ 核心：抽取 JSON
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
	history []*entity.BrowserAgentMessage,
) (*ws.Action, error) {

	decidePrompt := s.buildPrompt(task, pageState, history)

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
	history []*entity.BrowserAgentMessage,
) (*ws.Action, bool, error) {

	nextActionPrompt := s.buildNextPrompt(pageState, history)

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
	history []*entity.BrowserAgentMessage,
) string {

	return "【用户目标】\n" +
		task + "\n\n" +
		s.buildPageStateSection(pageState) +
		s.buildHistorySection(history)
}

func (s *BrowserAgentService) buildNextPrompt(
	pageState *ws.PageState,
	history []*entity.BrowserAgentMessage,
) string {

	return "【继续执行当前任务】\n\n" +
		s.buildPageStateSection(pageState) +
		s.buildHistorySection(history)
}

func (s *BrowserAgentService) buildPageStateSection(pageState *ws.PageState) string {
	if pageState == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("【页面状态】\n")
	sb.WriteString("URL: " + pageState.URL + "\n")
	sb.WriteString("可交互元素:\n")

	for i, elem := range pageState.Elements {
		if elem.Visible && !elem.Disabled {
			sb.WriteString(fmt.Sprintf(
				"%d. <%s> %s | selector=%s\n",
				i+1,
				elem.Tag,
				elem.Text,
				elem.Selector,
			))
		}
	}

	return sb.String()
}

func (s *BrowserAgentService) buildHistorySection(history []*entity.BrowserAgentMessage) string {
	if len(history) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("【历史记录】\n")

	for _, h := range history {
		sb.WriteString(fmt.Sprintf(
			"- %s: %s\n",
			h.Role,
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
