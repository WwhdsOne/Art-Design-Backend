package ws

import "sync"

// Hub 是 WebSocket 连接的中心管理器（Connection Hub）
//
// 设计职责：
//  1. 统一管理所有 WebSocket Client 的生命周期
//  2. 保证同一个 ConversationID 同一时间只存在一个连接
//  3. 处理 Client 的注册与注销（线程安全）
//
// 并发模型说明：
//   - register / unregister 通过 channel 串行化处理
//   - clients map 通过 RWMutex 保证并发安全
//   - Run 方法应在单独 goroutine 中长期运行
type Hub struct {

	// clients 保存当前在线的 WebSocket 客户端
	// key   : ConversationID
	// value : 对应的 Client 连接
	clients map[int64]*Client

	// clientsMux 用于保护 clients map 的并发读写
	clientsMux sync.RWMutex

	// register 用于接收新连接的注册请求
	// Client 建立成功后，通过该 channel 发送给 Hub
	register chan *Client

	// unregister 用于接收连接断开的注销请求
	// Client 主动关闭或异常退出时发送
	unregister chan *Client
}

// NewHub 创建一个新的 Hub 实例
//
// 注意：
//   - 创建后必须调用 go hub.Run() 启动事件循环
//   - Hub 本身不启动 goroutine，交由调用方控制
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
	}
}

// Run 启动 Hub 的事件循环
//
// 该方法会阻塞执行，通常应在单独的 goroutine 中启动：
//
//	go hub.Run()
//
// 功能说明：
//   - 监听 register / unregister channel
//   - 保证对 clients map 的修改是串行且安全的
func (h *Hub) Run() {
	for {
		select {

		// 处理新客户端注册
		case client := <-h.register:
			h.clientsMux.Lock()

			// 如果同一个 ConversationID 已存在连接
			// 则主动关闭旧连接，保证唯一性
			if old, ok := h.clients[client.ConversationID]; ok {
				old.Close()
			}

			// 注册新连接
			h.clients[client.ConversationID] = client
			h.clientsMux.Unlock()

		// 处理客户端注销
		case client := <-h.unregister:
			h.clientsMux.Lock()

			// 仅当当前连接与传入的 client 一致时才删除
			// 防止误删已被新连接替换的 Client
			if c, ok := h.clients[client.ConversationID]; ok && c == client {
				delete(h.clients, client.ConversationID)
			}

			h.clientsMux.Unlock()
		}
	}
}

// Register 向 Hub 注册一个新的 WebSocket Client
//
// 该方法是并发安全的：
//   - 实际注册逻辑在 Run() 的事件循环中执行
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister 向 Hub 注销一个 WebSocket Client
//
// 通常在以下场景调用：
//   - 客户端主动关闭连接
//   - 读写发生错误
//   - 服务端主动断开连接
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
