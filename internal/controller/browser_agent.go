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

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type BrowserAgentController struct {
	service  *service.BrowserAgentService
	hub      *ws.Hub
	upgrader websocket.Upgrader
}

func NewBrowserAgentController(r *gin.Engine,
	mws *middleware.Middlewares, svc *service.BrowserAgentService, hub *ws.Hub) *BrowserAgentController {
	browserAgentCtrl := &BrowserAgentController{
		service: svc,
		hub:     hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  65536,
			WriteBufferSize: 65536,
			// 允许所有源
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}
	agent := r.Group("/api").Group("/browser-agent")
	agent.Use(mws.AuthMiddleware())
	{
		agent.POST("/conversation/create", browserAgentCtrl.CreateConversation)
		agent.GET("/conversation/list", browserAgentCtrl.ListConversations)
		agent.POST("/conversation/rename", browserAgentCtrl.RenameConversation)
		agent.DELETE("/conversation/delete", browserAgentCtrl.DeleteConversation)
		agent.GET("/messages", browserAgentCtrl.ListMessages)
		agent.POST("/message/create", browserAgentCtrl.CreateMessage)
		agent.GET("/actions", browserAgentCtrl.ListActions)
	}

	wsGroup := r.Group("/api").Group("/browser-agent").Group("/ws")
	// ws必须使用ws特别鉴权，否则会导致当作普通http返回
	// 因为ws请求头无法修改，必须放在url里面
	wsGroup.Use(mws.WSAuthMiddleware())
	{
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

	resp, err := ctrl.service.CreateConversation(c, &req)
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
	resp, err := ctrl.service.ListConversations(c, userID, &queryParam)
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
	if err := ctrl.service.RenameConversation(c, &req); err != nil {
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

	if err = ctrl.service.DeleteConversation(c, conversationID); err != nil {
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

	resp, err := ctrl.service.CreateMessage(c, &req)
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

	resp, err := ctrl.service.ListMessages(c, &req)
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

	resp, err := ctrl.service.ListActions(c, &req)
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
		Service:        ctrl.service,
		Ctx:            clientCtx,
		Cancel:         cancel,
	}

	ctrl.hub.Register(client)
	go client.WritePump()
	client.ReadPump()
}
