package types

import (
	"fmt"
)

const (
	CmdNormal    = ``     // 不做任何动作
	CmdHeartBeat = `1314` // 心跳
)

const (
	CodeNormal             int = 200 // 正常
	CodeAuthFail           int = 222 // 鉴权失败
	CodeNoMoreData         int = 300 // 没有更多数据
	CodeUndefinedCallback  int = 404 // 回调方法未注册
	CodeInvalidData        int = 414 // 无效提交数据
	CodeOutOfQuotas        int = 432 // 超出配额
	CodeOperateFailed      int = 444 // 后端操作失败,客户端需要处理重试机制
	CodeDataStructWrong    int = 110 // 数据结构有误
	CodeServiceUnavailable int = 500 // 后端服务不可用
)

var code2MessageMap = map[int]string{
	CodeNormal:             "ok",
	CodeAuthFail:           "auth fail",
	CodeUndefinedCallback:  "callback function undefined",
	CodeOperateFailed:      "operate failed, need retry",
	CodeDataStructWrong:    "data struct wrong",
	CodeServiceUnavailable: "back-end service is not available",
}

const (
	TypeReply = "reply"
	TypeCmd   = "cmd"
)

const (
	ActionAuth  = "auth"
	ActionReply = "reply"
	ActionPong  = "pong"

	ActionEventsReport = "events-report"

	ActionIdiomWatchImgGuess = `idiom-watch-img-guess`
	ActionIdiomGuess         = `idiom-guess`
	ActionUnlimitedGuess     = `unlimited-guess`
	ActionSingleChallenge    = `single-challenge`

	ActionFeedback = `feedback`
	ActionDigg     = `digg`
	ActionFavorite = `favorite`
	ActionComment  = `comment`
)

const (
	BroadcastAll int64 = -1
)

type ClientMsgBase struct {
	Mid         string `json:"mid"`
	AccessToken string `json:"access_token"`
	Action      string `json:"action"`
}

type EventReportItem struct {
	EventID         int64    `json:"event_id"`          // 事件ID
	RelatedSN       string   `json:"related_sn"`        // 关联数据的唯一编号,可以是学习资料编号,题库编号
	StartAt         int64    `json:"start_at"`          // 事件开始时间
	EndAt           int64    `json:"end_at"`            // 事件结束时间
	Duration        int64    `json:"duration"`          // 客户计算的用户用时
	AppUserSolution []string `json:"app_ans,omitempty"` // 用户提交的解题答案
	JudgeResult     int      `json:"judge,omitempty"`   // 客户端对用户答案的判定结果. 0: 答错了; 1: 答对了
}

type ClientMessageBody struct {
	LevelNum   int               `json:"level_num,omitempty"` // 关卡数, TODO: 一期客户端本地记录
	EventsList []EventReportItem `json:"events_list,omitempty"`

	SingleEvent *EventReportItem `json:"event,omitempty"` // 单事件

	// 评论/点赞/收藏
	ObjSN   string `json:"obj_sn,omitempty"`  // 对象编号,可以是任何分布式编号
	Content string `json:"content,omitempty"` // 反馈/评论的具体内容

	// 点赞/收藏
	OpCode int `json:"op_code,omitempty"` // 1: 赞/收藏; -1: 踩/取消收藏

	// 反馈相关
	Type    int    `json:"type,omitempty"`  // 反馈的类型
	Img1    string `json:"img_1,omitempty"` // 反馈时上传的图片
	Img2    string `json:"img_2,omitempty"`
	Img3    string `json:"img_3,omitempty"`
	Contact string `json:"contact,omitempty"` // 可选的联系方式
}

type ClientMsg struct {
	ClientMsgBase
	Message ClientMessageBody `json:"message"`
}

type ServerMsg struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Mid     string      `json:"mid"`
	Type    string      `json:"type"`
	Action  string      `json:"action"`
	Cmd     string      `json:"cmd"`
	Data    interface{} `json:"data"`
}

type ServerBroadcastMsg struct {
	Sender   int64 `json:"sender"`
	Receiver int64 `json:"receiver"`
	ServerMsg
}

func SingleChannelName(accountID int64) string {
	return fmt.Sprintf(`fm:ch:sgl:%d`, accountID)
}

func PublicChannelName() string {
	return fmt.Sprintf(`fm:ch:public-broadcast`)
}

// 服务端应答包构造助手函数,有点重复,先这样吧...

func PingCmd(mid string) (msg ServerMsg) {
	msg.Code = CodeNormal
	msg.Type = TypeCmd
	msg.Cmd = CmdHeartBeat
	msg.Message = "ping pong"
	msg.Mid = mid
	msg.Data = struct{}{}

	return
}

func AuthFailReply() (reply ServerMsg) {
	reply.Code = CodeAuthFail
	reply.Type = TypeReply
	reply.Message = "认证失败,请检查设置"
	reply.Data = struct{}{}

	return
}

func UndefinedCallback(clientMsg ClientMsg) (reply ServerMsg) {
	reply.Code = CodeUndefinedCallback
	reply.Message = "undefined callback function name"
	reply.Cmd = CmdNormal
	reply.Mid = clientMsg.Mid
	reply.Type = TypeReply
	reply.Data = struct{}{}

	return
}

func Reply(code int, cmd string, mid string, data interface{}) (reply ServerMsg) {
	var message string = "unknown"
	if m, ok := code2MessageMap[code]; ok {
		message = m
	}

	reply.Code = code
	reply.Message = message
	reply.Mid = mid
	reply.Type = TypeReply
	reply.Cmd = cmd
	reply.Data = data

	return
}
