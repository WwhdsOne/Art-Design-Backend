package ws

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 65536
)

type BrowserAgentService interface {
	HandleTask(ctx context.Context, messageID int64, pageState *PageState) (*Action, error)
	HandleResult(ctx context.Context, conversationID int64, msg *ClientMessage) (*Action, bool, error)
}

type Client struct {
	Hub            *Hub
	Conn           *websocket.Conn
	ConversationID int64
	UserID         int64
	Send           chan []byte
	Service        BrowserAgentService
	Ctx            context.Context
	Cancel         context.CancelFunc
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var clientMsg ClientMessage
		if err := sonic.Unmarshal(message, &clientMsg); err != nil {
			c.sendError("消息格式错误")
			continue
		}

		switch clientMsg.Type {
		case "task":
			c.handleTask(&clientMsg)
		case "result":
			c.handleResult(&clientMsg)
		default:
			c.sendError("未知的消息类型")
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.Ctx.Done():
			return
		}
	}
}

func (c *Client) Close() {
	c.Cancel()
}

func (c *Client) handleTask(msg *ClientMessage) {
	action, err := c.Service.HandleTask(c.Ctx, msg.MessageID, msg.PageState)
	if err != nil {
		c.sendError(err.Error())
		return
	}
	c.sendAction(action)
}

func (c *Client) handleResult(msg *ClientMessage) {
	action, finished, err := c.Service.HandleResult(c.Ctx, c.ConversationID, msg)
	if err != nil {
		c.sendError(err.Error())
		return
	}
	if finished {
		c.sendFinish("任务已完成")
		return
	}
	c.sendAction(action)
}

func (c *Client) sendAction(action *Action) {
	msg := ServerMessage{Type: "action", Action: action}
	data, _ := sonic.Marshal(msg)
	c.Send <- data
}

func (c *Client) sendFinish(message string) {
	msg := ServerMessage{Type: "finish", Message: message}
	data, _ := sonic.Marshal(msg)
	c.Send <- data
}

func (c *Client) sendError(message string) {
	msg := ServerMessage{Type: "error", Message: message}
	data, _ := sonic.Marshal(msg)
	c.Send <- data
}
