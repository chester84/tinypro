package wsprocess

import (
	"tinypro/common/types"
)

func init() {
	Register(types.ActionAuth, Auth)
}

func Auth(accountID int64, clientMsg types.ClientMsg) (reply types.ServerMsg) {
	reply.Type = types.TypeReply
	reply.Mid = clientMsg.Mid

	reply.Code = types.CodeNormal
	reply.Message = `认证成功`
	reply.Action = types.ActionAuth
	reply.Data = struct{}{}

	return
}
