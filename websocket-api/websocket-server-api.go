package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrade = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 向客户端写数据
func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		conn *websocket.Conn
		err  error
		data []byte
	)
	// 升级为WebSocket连接，并进行握手
	if conn, err = upgrade.Upgrade(w, r, nil); err != nil {
		return
	}
	// 死循环
	for {
		// 使用conn.ReadMessage()读出接收到的data
		if _, data, err = conn.ReadMessage(); err != nil {
			_ = conn.Close()
		}
		// 将data通过conn.WriteMessage()原路返回给客户端
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			_ = conn.Close()
		}
	}
}

func main() {
	// 监听地址及端口
	_ = http.ListenAndServe("0.0.0.0:7777", nil)
	// 通过wsHandler中的upgrade升级为WebSocket长连接
	http.HandleFunc("/ws", wsHandler)
}
