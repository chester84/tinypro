package controllers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"

	"tinypro/common/pkg/accesstoken"
	"tinypro/common/pkg/wsh"
	"tinypro/common/pkg/wsprocess"
	"tinypro/common/types"
)

type WebsocketController struct {
	beego.Controller
	AccountID int64
	Ws        *websocket.Conn
}

func (c *WebsocketController) Prepare() {
	// 调用上一级的 Prepare 方法
	//c.Controller.Prepare()

	// Upgrade from http request to WebSocket.
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		c.StopRun()
		return
	} else if err != nil {
		logs.Error("[WebsocketController->Prepare] cannot setup WebSocket connection, err: %v", err)
		c.CustomAbort(401, "Cannot setup WebSocket connection")
		return
	}

	c.Ws = ws
}

func (c *WebsocketController) Get() {
	var running = false
	done := make(chan struct{})

	for {
		var wl sync.Mutex

		_, msgByte, err := c.Ws.ReadMessage()
		if err != nil {
			wl.Lock()
			if !wsh.HadBroken(done) {
				close(done) //! 很重要
			}
			wl.Unlock()

			logs.Error("[WebsocketController] read message has wrong, err: %v", err)
			return
		}

		//logs.Debug("[WebsocketController] msg: %s", string(msgByte))

		// 1. 检查数据包的合法性,非法不容许进入
		clientMsg := types.ClientMsg{}
		errDc := json.Unmarshal(msgByte, &clientMsg)
		// 数据包格式不正确,关闭连接,禁止接入
		if errDc != nil {
			logs.Error("[WebsocketController] client message structure is wrong, would close connection, clientMsg: ", string(msgByte))

			// long lock {{{
			wl.Lock()

			if !wsh.HadBroken(done) {
				close(done) //! 要关闭连接了...由于有时序的问题,需要加锁
			}
			// 服务端主动关闭连接
			errClose := c.Ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if errClose != nil {
				logs.Error("[WebsocketController] close websocket has wrong, errClose; %v", errClose)
			}

			wl.Unlock()
			// end lock }}}

			return
		}

		//logs.Critical("current accountID: %d", c.AccountID)

		// 2. 鉴权和识别用户
		if c.AccountID == 0 || clientMsg.Action == types.ActionAuth {
			logs.Notice("[WebsocketController] auto auth message: %s", string(msgByte))
			ok, accountID := accesstoken.IsValidAccessToken(types.PlatformAndroid, clientMsg.AccessToken)
			if !ok {
				logs.Error("[WebsocketController] token not found, clientMsg.Token: %s", clientMsg.AccessToken)

				// 鉴权失败,关闭连接
				logs.Error("[WebsocketController] auth fail, would close connection, clientMsg: ", string(msgByte))

				reply := types.AuthFailReply()
				message, _ := json.Marshal(reply)

				// long lock {{{
				wl.Lock()

				//errW := c.Ws.WriteMessage(websocket.BinaryMessage, message)
				errW := c.Ws.WriteMessage(websocket.TextMessage, message)
				if errW != nil {
					logs.Error("[WebsocketController] write message has wrong, errW: %v", errW)
				}

				if !wsh.HadBroken(done) {
					close(done)
				}

				errClose := c.Ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if errClose != nil {
					logs.Error("[WebsocketController] close websocket has wrong, errClose; %v", errClose)
				}

				wl.Unlock()
				// end lock }}}

				return
			} else {
				c.AccountID = accountID
				logs.Notice("[WebsocketController] auth pass, accountID: %d", accountID)
			}
		}

		if !running {
			// 3.1 心跳保活,每个活跃连接只需要启一个协程即可
			go wsh.HeartBeat(done, &wl, c.Ws, c.AccountID, "/ws/v1")
			// 3.2 接入单通道广播
			//go wsh.SingleChannelBroadcast(done, &wl, c.Ws, c.AccountID)
			// 3.3 接入公共广播
			go wsh.PublicChannelBroadcast(done, &wl, c.Ws, c.AccountID)

			running = true
		}

		// 4. 接收消息并回复客户端
		logs.Debug("[WebsocketController] device is working, accountID: %d", c.AccountID)
		logs.Debug("[WebsocketController] msg: %s", string(msgByte))

		// 保活/认证
		if clientMsg.Action == types.ActionPong {
			// 客户端回 pong,没有实际逻辑
			logs.Notice("[WebsocketController] device pong, accountID: %d", c.AccountID)
			continue
		}

		var reply types.ServerMsg

		// 调用注册的回调函数,开始工作
		processFuncBox := wsprocess.ProcessFuncBox()
		if funcName, ok := processFuncBox[clientMsg.Action]; ok {
			reply = funcName(c.AccountID, clientMsg)
		} else {
			logs.Notice("[WebsocketController] undefined callback function, clientMsg: %s", string(msgByte))
			reply = types.UndefinedCallback(clientMsg)
		}

		message, _ := json.Marshal(reply)

		wl.Lock()
		//errW := c.Ws.WriteMessage(websocket.BinaryMessage, message)
		errW := c.Ws.WriteMessage(websocket.TextMessage, message)
		wl.Unlock()

		if errW != nil {
			wl.Lock()
			if !wsh.HadBroken(done) {
				close(done)
			}
			wl.Unlock()

			logs.Error("[WebsocketController] write message has wrong, message: %s, errW: %v", string(message), errW)
			return
		}
	}
}
