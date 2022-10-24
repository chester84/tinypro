package wsh

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"

	"github.com/chester84/libtools"
	"tinypro/common/types"
)

const healthCheckPeriod = time.Minute

var redisServerAddr string

func init() {
	redisHost, _ := config.String("cache_redis_host")
	redisPort, _ := config.Int("cache_redis_port")
	redisServerAddr = fmt.Sprintf("%s:%d", redisHost, redisPort)
}

func HadBroken(bc <-chan struct{}) bool {
	select {
	case <-bc:
		return true
	default:
		return false
	}
}

// HeartBeat 长连接服务端心跳检测机制,如果写失败,要关闭长连接,让客户端重新连接
func HeartBeat(bc chan struct{}, wl *sync.Mutex, ws *websocket.Conn, accountID int64, router string) {
	logs.Notice("[HeartBeat] for accountID: %d, router: %s", accountID, router)

	for {
		if HadBroken(bc) {
			logs.Informational("[HeartBeat] websocket connection has broken normally, accountID: %d, router: %s", accountID, router)
			// 用户掉线/下线
			// TODO

			logs.Informational("[HeartBeat] [HadBroken] get yes, device off-line, accountID: %d", accountID)
			return
		}

		msg := types.PingCmd(libtools.Int642Str(libtools.GetUnixMillis()))
		data, err := json.Marshal(msg)
		if err != nil {
			logs.Error("[HeartBeat] fail to marshal msg, err: %v, accountID: %d, router: %s", err, accountID, router)
			return
		}

		wl.Lock()
		err = ws.WriteMessage(websocket.PingMessage, data)
		wl.Unlock()

		if err != nil {
			logs.Error("[HeartBeat] websocket connection has broken unexpected, accountID: %d, router: %s, %#v", accountID, router, err)

			wl.Lock()
			if !HadBroken(bc) {
				close(bc)
			}
			wl.Unlock()

			// 用户掉线掉线
			// TODO

			logs.Informational("[HeartBeat] device off-line, accountID: %d", accountID)

			return
		}

		//logs.Critical("send ping: %s", string(data))

		// 用户在线
		// TODO

		//time.Sleep(time.Minute)
		time.Sleep(43 * time.Second)
	}
}

func SingleChannelBroadcast(bc chan struct{}, wl *sync.Mutex, ws *websocket.Conn, accountID int64) {
	channelName := types.SingleChannelName(accountID)
	logs.Notice("[SingleChannelBroadcast] start single channel broadcast, channelName: %s", channelName)

	err := PubSubChannels(
		func() error {
			// The start callback is a good place to backfill missed
			logs.Notice("[SingleChannelBroadcast] onStart for channelName: %s", channelName)
			return nil
		},
		func(channel string, message []byte) error {
			wl.Lock()
			errW := ws.WriteMessage(websocket.TextMessage, message)
			if errW != nil {
				// 任何写失败,需要关闭所有资源
				logs.Error("[SingleChannelBroadcast] write message has wrong, channelName: %s, message: %s, errW: %v", channelName, string(message), errW)

				if !HadBroken(bc) {
					close(bc) // 广播事件!
				}
			}
			wl.Unlock()

			return nil
		},
		HadBroken,
		bc,
		channelName)

	if err != nil {
		logs.Error("[SingleChannelBroadcast] the program quits abnormally, channelName: %s, err: %v", channelName, err)
		return
	}
}

func PublicChannelBroadcast(bc chan struct{}, wl *sync.Mutex, ws *websocket.Conn, accountID int64) {
	channelName := types.PublicChannelName()
	logs.Notice("[PublicChannelBroadcast] start single channel broadcast, channelName: %s", channelName)

	err := PubSubChannels(
		func() error {
			// The start callback is a good place to backfill missed
			logs.Notice("[PublicChannelBroadcast] onStart for channelName: %s", channelName)
			return nil
		},
		func(channel string, message []byte) error {
			// 解析协议,如果和当前用户无法,跳过
			var broadcast types.ServerBroadcastMsg
			errJ := json.Unmarshal(message, &broadcast)
			if errJ != nil {
				logs.Error("[PublicChannelBroadcast] broadcast message json decode exception, msg: %s, err: %v", string(message), errJ)
				return nil
			}
			if broadcast.Receiver != types.BroadcastAll && broadcast.Receiver != accountID {
				logs.Info("[PublicChannelBroadcast] broadcast message has nothing to do with the current user. msg: %s, accountID: %d", string(message), accountID)
				return nil
			}

			wl.Lock()
			errW := ws.WriteMessage(websocket.TextMessage, message)
			if errW != nil {
				// 任何写失败,需要关闭所有资源
				logs.Error("[PublicChannelBroadcast] write message has wrong, channelName: %s, message: %s, errW: %v", channelName, string(message), errW)

				if !HadBroken(bc) {
					close(bc) // 广播事件!
				}
			}
			wl.Unlock()

			return nil
		},
		HadBroken,
		bc,
		channelName)

	if err != nil {
		logs.Error("[PublicChannelBroadcast] the program quits abnormally, channelName: %s, err: %v", channelName, err)
		return
	}
}

// PubSubChannels websocket 的辅助方法,用于服务端向客户端推送消息
func PubSubChannels(
	onStart func() error,
	onMessage func(channel string, data []byte) error,
	hadBroken func(bc <-chan struct{}) bool,
	bc <-chan struct{},
	channels ...string) error {
	// A ping is set to the server with this period to test for the health of
	// the connection and server.
	c, err := redis.Dial("tcp", redisServerAddr,
		// Read timeout on server should be greater than ping period.
		redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
		redis.DialWriteTimeout(10*time.Second))
	if err != nil {
		logs.Error("[PubSubChannels] redis server unavailable, redisServerAddr: %s, channels: %v", redisServerAddr, channels)
		return err
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		logs.Error("[PubSubChannels] redis Subscribe has wrong, redisServerAddr: %s, channels: %v, err: %v",
			redisServerAddr, channels, err)
		return err
	}

	done := make(chan error, 1)
	wsDone := make(chan struct{})

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			// websocket 连接已断开,需要正确的 Unsubscribe
			if hadBroken(bc) {
				close(wsDone) //! 需要正确的处理 Unsubscribe

				logs.Notice("[PubSubChannels] has broken, channels: %v", channels)
				return
			}

			switch n := psc.Receive().(type) {
			case error:
				done <- n
				logs.Error("[PubSubChannels] psc.Receive().(type) has wrong, channels: %v, n: %v", channels, n)
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					logs.Error("[PubSubChannels] onMessage has wrong and then need exit, channels: %v, err: %v", channels, err)
					return
				}
			case redis.Subscription:
				switch n.Count {
				case len(channels):
					// Notify application when all channels are subscribed.
					if err := onStart(); err != nil {
						done <- err
						logs.Error("[PubSubChannels] onStart has wrong, channels: %v, err: %v", channels, err)
						return
					}
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					logs.Informational("[PubSubChannels] all channels are unsubscribed,  channels: %v", channels)
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				logs.Error("[PubSubChannels] ping pong is not received, channels: %v, err: %v", channels, err)
				break loop
			}
		case <-wsDone:
			logs.Informational("[PubSubChannels] websocket link had broken, need Unsubscribe, channels: %v", channels)
			break loop
		case err := <-done:
			// Return error from the receive goroutine.
			logs.Error("[PubSubChannels] return error from the receive goroutine, channels: %v, err: %v", channels, err)
			return err
		}
	}

	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	errU := psc.Unsubscribe(redis.Args{}.AddFlat(channels)...)
	if errU != nil {
		logs.Error("[PubSubChannels] redis> Unsubscribe %v , err: %v", channels, errU)
	}
	logs.Notice("[PubSubChannels] unsubscribing from all channels: %v", channels)

	// Wait for goroutine to complete.
	return <-done
}
