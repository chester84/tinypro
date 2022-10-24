// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/websocket"

	"github.com/chester84/libtools"
	"tinypro/common/types"
)

var (
	wsPath = "/ws/v1"

	env   = flag.String("env", "dev", `set env, [dev|prod]`)
	token = flag.String("token", "", `有效token`)
	help  = flag.Bool("h", false, `show help`)

	action string
)

func init() {
	flag.StringVar(&action, `action`, ``, `action
idiom-watch-img-guess
unlimited-guess
idiom-guess
single-challenge
feedback
digg
favorite
comment
`)
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	accessToken := *token

	var wLock = sync.Mutex{}

	var apiUrl url.URL
	switch *env {
	case `dev`:
		apiUrl = url.URL{Scheme: "ws", Host: "127.0.0.1:8565", Path: wsPath}
		//accessToken = `a0a3918a8b881eaa6dcb6034a1e7d662`

	case `prod`:
		apiUrl = url.URL{Scheme: "wss", Host: "api.indianpandit.in", Path: wsPath}
		//accessToken = `41e99520011e5a8e5160076e4b4aeb84`

	default:
		flag.PrintDefaults()
		os.Exit(0)
	}

	logs.Info("connecting to %s", apiUrl.String())

	var clientMsg types.ClientMsg
	switch action {
	case types.ActionIdiomWatchImgGuess:
		clientMsg = types.ClientMsg{
			ClientMsgBase: types.ClientMsgBase{
				Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
				AccessToken: accessToken,
				Action:      types.ActionIdiomWatchImgGuess,
			},
		}

	case types.ActionUnlimitedGuess:
		clientMsg = types.ClientMsg{
			ClientMsgBase: types.ClientMsgBase{
				Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
				AccessToken: accessToken,
				Action:      types.ActionUnlimitedGuess,
			},
			Message: types.ClientMessageBody{
				LevelNum: 1,
			},
		}

	case types.ActionIdiomGuess:
		clientMsg = types.ClientMsg{
			ClientMsgBase: types.ClientMsgBase{
				Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
				AccessToken: accessToken,
				Action:      types.ActionIdiomGuess,
			},
			Message: types.ClientMessageBody{
				LevelNum: 1,
			},
		}

	case types.ActionFeedback:
		clientMsg = types.ClientMsg{
			ClientMsgBase: types.ClientMsgBase{
				Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
				AccessToken: accessToken,
				Action:      types.ActionFeedback,
			},
			Message: types.ClientMessageBody{
				Type:    1,
				Content: `测试反馈003`,
				Contact: `QQ: 209876788`,
				//Img1:    `iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAYAAADED76LAAAAE0lEQVR4nGJiIABGiAJAAAAA`,
				//Img2:    `iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAYAAADED76LAAAAE0lEQVR4nGJiIABGiAJAAAAA`,
				//Img3:    `iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAYAAADED76LAAAAE0lEQVR4nGJiIABGiAJAAAAA`,
			},
		}

	case types.ActionDigg:
		clientMsg = types.ClientMsg{
			ClientMsgBase: types.ClientMsgBase{
				Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
				AccessToken: accessToken,
				Action:      types.ActionDigg,
			},
			Message: types.ClientMessageBody{
				ObjSN:  `200229860000000279`,
				OpCode: 1,
			},
		}

	case types.ActionFavorite:
		clientMsg = types.ClientMsg{
			ClientMsgBase: types.ClientMsgBase{
				Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
				AccessToken: accessToken,
				Action:      types.ActionFavorite,
			},
			Message: types.ClientMessageBody{
				ObjSN:  `200229860000000279`,
				OpCode: -1,
			},
		}

	case types.ActionComment:
		clientMsg = types.ClientMsg{
			ClientMsgBase: types.ClientMsgBase{
				Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
				AccessToken: accessToken,
				Action:      types.ActionComment,
			},
			Message: types.ClientMessageBody{
				ObjSN:   `200229860000000279`,
				Content: `测试评论003`,
			},
		}

	default:
		flag.PrintDefaults()
		logs.Error("不支持的action")
		os.Exit(0)
	}

	clientMsgJson, _ := json.MarshalIndent(clientMsg, ``, `  `)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(apiUrl.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read get exception, err:", err)
				return
			}
			logs.Info("recv: %s, messageType: %d", message, messageType)
		}
	}()

	go func() {
		defer close(done)
		for {
			pong := types.ClientMsg{
				ClientMsgBase: types.ClientMsgBase{
					Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
					AccessToken: accessToken,
					Action:      types.ActionPong,
				},
			}
			pongJson, _ := json.MarshalIndent(&pong, "", "  ")
			logs.Info("pong: %s", string(pongJson))

			wLock.Lock()
			err := c.WriteMessage(websocket.PingMessage, pongJson)
			wLock.Unlock()

			if err != nil {
				logs.Error("ping get error: %v", err)
				return
			}

			time.Sleep(10 * time.Second)
		}
	}()

	auth := types.ClientMsg{
		ClientMsgBase: types.ClientMsgBase{
			Mid:         libtools.Int642Str(libtools.GetUnixMillis()),
			AccessToken: accessToken,
			Action:      types.ActionAuth,
		},
	}
	msgJson, _ := json.MarshalIndent(auth, ``, `  `)
	logs.Info("auth msg: %s", string(msgJson))

	wLock.Lock()
	err = c.WriteMessage(websocket.TextMessage, msgJson)
	wLock.Unlock()

	if err != nil {
		logs.Error("auth fail. err:", err)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			logs.Info("time.NewTicker at:", libtools.UnixMsec2Date(t.Unix()*1000, `Y-m-d H:i:s`))

			logs.Info("action: %s, msg: %s", action, string(clientMsgJson))

			wLock.Lock()
			err := c.WriteMessage(websocket.BinaryMessage, clientMsgJson)
			wLock.Unlock()

			if err != nil {
				logs.Error("write wrong, err:", err)
				return
			}

		case <-interrupt:
			logs.Notice("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.

			wLock.Lock()
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			wLock.Unlock()

			if err != nil {
				logs.Error("write close:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}

			return
		}
	}
}
