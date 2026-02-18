package controller

import (
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"Art-Design-Backend/pkg/ws"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type BrowserAgentController struct {
	browserAgentService          *service.BrowserAgentService
	browserAgentDashboardService *service.BrowserAgentDashboardService
	hub                          *ws.Hub
	upgrader                     websocket.Upgrader
}

func NewBrowserAgentController(
	r *gin.Engine,
	mws *middleware.Middlewares,
	bas *service.BrowserAgentService, bads *service.BrowserAgentDashboardService,
	hub *ws.Hub) *BrowserAgentController {
	browserAgentCtrl := &BrowserAgentController{
		browserAgentService:          bas,
		browserAgentDashboardService: bads,
		hub:                          hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  65536,
			WriteBufferSize: 65536,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}
	{
		agent := r.Group("/api").
			Group("/browser-agent")
		agent.Use(mws.AuthMiddleware())
		agent.POST("/conversation/create", browserAgentCtrl.CreateConversation)
		agent.GET("/conversation/list", browserAgentCtrl.ListConversations)
		agent.POST("/conversation/rename", browserAgentCtrl.RenameConversation)
		agent.DELETE("/conversation/delete", browserAgentCtrl.DeleteConversation)
		agent.GET("/messages", browserAgentCtrl.ListMessages)
		agent.POST("/message/create", browserAgentCtrl.CreateMessage)
		agent.GET("/actions", browserAgentCtrl.ListActions)
	}

	{
		adminDashboard := r.
			Group("/api").
			Group("/browser-agent").
			Group("/dashboard").
			Group("/admin")
		adminDashboard.Use(mws.AuthMiddleware())
		adminDashboard.GET("/summary", browserAgentCtrl.GetAdminSummary)
		adminDashboard.GET("/weekly-task-volume", browserAgentCtrl.GetAdminWeeklyTaskVolume)
		adminDashboard.GET("/weekly-task-success-rate", browserAgentCtrl.GetAdminWeeklyTaskSuccessRate)
		adminDashboard.GET("/total-task-volume", browserAgentCtrl.GetAdminTotalTaskVolume)
		adminDashboard.GET("/task-classification", browserAgentCtrl.GetAdminTaskClassification)
		adminDashboard.GET("/weekly-operation-volume", browserAgentCtrl.GetAdminWeeklyOperationVolume)
		adminDashboard.GET("/weekly-operation-success-rate", browserAgentCtrl.GetAdminWeeklyOperationSuccessRate)
		adminDashboard.GET("/active-sessions", browserAgentCtrl.GetAdminActiveSessions)
		adminDashboard.GET("/annual-task-stats", browserAgentCtrl.GetAdminAnnualTaskStats)
		adminDashboard.GET("/hot-task-list", browserAgentCtrl.GetAdminHotTaskList)
		adminDashboard.POST("/messages", browserAgentCtrl.GetMessagePage)
		adminDashboard.GET("/actions", browserAgentCtrl.GetActionsByMessageID)
	}

	{
		userDashboard := r.
			Group("/api").
			Group("/browser-agent").
			Group("/dashboard").
			Group("/user")
		userDashboard.Use(mws.AuthMiddleware())
		userDashboard.GET("/summary", browserAgentCtrl.GetUserSummary)
		userDashboard.GET("/weekly-task-volume", browserAgentCtrl.GetUserWeeklyTaskVolume)
		userDashboard.GET("/weekly-task-success-rate", browserAgentCtrl.GetUserWeeklyTaskSuccessRate)
		userDashboard.GET("/task-overview", browserAgentCtrl.GetUserTaskOverview)
		userDashboard.GET("/task-trend", browserAgentCtrl.GetUserTaskTrend)
	}

	{
		wsGroup := r.Group("/api").Group("/browser-agent").Group("/ws")
		wsGroup.Use(mws.WSAuthMiddleware())
		wsGroup.GET("/:conversationId", browserAgentCtrl.HandleWebSocket)
	}
	return browserAgentCtrl
}

func (ctrl *BrowserAgentController) CreateConversation(c *gin.Context) {
	var req request.CreateConversationRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := ctrl.browserAgentService.CreateConversation(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) ListConversations(c *gin.Context) {
	var queryParam query.BrowserAgentConversation
	if err := c.ShouldBindQuery(&queryParam); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	userID := authutils.GetUserID(c)
	resp, err := ctrl.browserAgentService.ListConversations(c, userID, &queryParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) RenameConversation(c *gin.Context) {
	var req request.RenameConversationRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}
	if err := ctrl.browserAgentService.RenameConversation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	result.OkWithMessage("重命名成功", c)
}

func (ctrl *BrowserAgentController) DeleteConversation(c *gin.Context) {
	idStr := c.Query("id")
	conversationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.FailWithMessage("无效的ID", c)
		return
	}

	if err = ctrl.browserAgentService.DeleteConversation(c, conversationID); err != nil {
		_ = c.Error(err)
		return
	}

	result.OkWithMessage("删除成功", c)
}

func (ctrl *BrowserAgentController) CreateMessage(c *gin.Context) {
	var req request.CreateMessageRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := ctrl.browserAgentService.CreateMessage(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) ListMessages(c *gin.Context) {
	var req request.GetMessagesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := ctrl.browserAgentService.ListMessages(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) ListActions(c *gin.Context) {
	var req request.GetActionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := ctrl.browserAgentService.ListActions(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) HandleWebSocket(c *gin.Context) {
	conversationIDStr := c.Param("conversationId")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		zap.L().Error("无效的会话ID", zap.Error(err))
		result.FailWithMessage("无效的会话ID", c)
		return
	}

	conn, err := ctrl.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("升级为WebSocket失败", zap.Error(err))
		return
	}

	clientCtx, cancel := context.WithCancel(context.Background())
	client := &ws.Client{
		Hub:            ctrl.hub,
		Conn:           conn,
		ConversationID: conversationID,
		UserID:         authutils.GetUserID(c),
		Send:           make(chan []byte, 256),
		Service:        ctrl.browserAgentService,
		Ctx:            clientCtx,
		Cancel:         cancel,
	}

	ctrl.hub.Register(client)
	go client.WritePump()
	client.ReadPump()
}

func (ctrl *BrowserAgentController) GetAdminSummary(c *gin.Context) {
	resp, err := ctrl.browserAgentDashboardService.GetAdminSummary(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminWeeklyTaskVolume(c *gin.Context) {
	resp, err := ctrl.browserAgentDashboardService.GetAdminWeeklyTaskVolume(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminWeeklyTaskSuccessRate(c *gin.Context) {
	resp, err := ctrl.browserAgentDashboardService.GetAdminWeeklyTaskSuccessRate(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminTotalTaskVolume(c *gin.Context) {
	resp, err := ctrl.browserAgentDashboardService.GetAdminTotalTaskVolume(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminTaskClassification(c *gin.Context) {
	resp, err := ctrl.browserAgentDashboardService.GetAdminTaskClassification(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminWeeklyOperationVolume(c *gin.Context) {
	resp, err := ctrl.browserAgentDashboardService.GetAdminWeeklyOperationVolume(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminWeeklyOperationSuccessRate(c *gin.Context) {
	resp, err := ctrl.browserAgentDashboardService.GetAdminWeeklyOperationSuccessRate(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminActiveSessions(c *gin.Context) {
	resp, err := ctrl.browserAgentDashboardService.GetAdminActiveSessions(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminAnnualTaskStats(c *gin.Context) {
	var req request.YearRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Year = time.Now().Year()
	}
	if req.Year == 0 {
		req.Year = time.Now().Year()
	}

	resp, err := ctrl.browserAgentDashboardService.GetAdminAnnualTaskStats(c, req.Year)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetAdminHotTaskList(c *gin.Context) {
	var req request.LimitRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Limit = 6
	}

	resp, err := ctrl.browserAgentDashboardService.GetAdminHotTaskList(c, req.Limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetMessagePage(c *gin.Context) {
	var queryParam query.BrowserAgentMessage
	if err := c.ShouldBindBodyWithJSON(&queryParam); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := ctrl.browserAgentDashboardService.GetMessagePage(c, &queryParam)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetUserSummary(c *gin.Context) {
	userID := authutils.GetUserID(c)
	resp, err := ctrl.browserAgentDashboardService.GetUserSummary(c, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetUserWeeklyTaskVolume(c *gin.Context) {
	userID := authutils.GetUserID(c)
	resp, err := ctrl.browserAgentDashboardService.GetUserWeeklyTaskVolume(c, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetUserWeeklyTaskSuccessRate(c *gin.Context) {
	userID := authutils.GetUserID(c)
	resp, err := ctrl.browserAgentDashboardService.GetUserWeeklyTaskSuccessRate(c, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetUserTaskOverview(c *gin.Context) {
	userID := authutils.GetUserID(c)
	resp, err := ctrl.browserAgentDashboardService.GetUserTaskOverview(c, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetUserTaskTrend(c *gin.Context) {
	var req request.YearRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Year = time.Now().Year()
	}
	if req.Year == 0 {
		req.Year = time.Now().Year()
	}

	userID := authutils.GetUserID(c)
	resp, err := ctrl.browserAgentDashboardService.GetUserTaskTrend(c, userID, req.Year)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(resp, c)
}

func (ctrl *BrowserAgentController) GetActionsByMessageID(c *gin.Context) {
	var req request.GetActionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := ctrl.browserAgentService.ListActions(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result.OkWithData(resp, c)
}
