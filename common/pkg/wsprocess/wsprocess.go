package wsprocess

import "tinypro/common/types"

type ProcessFunc func(accountID int64, clientMsg types.ClientMsg) (reply types.ServerMsg)

var processFuncBox = map[string]ProcessFunc{}

func Register(name string, funcName ProcessFunc) {
	if funcName == nil {
		panic("[wsprocess->Register] register func is nil")
	}

	if _, ok := processFuncBox[name]; ok {
		panic("register process func twice for func " + name)
	}

	processFuncBox[name] = funcName
}

func ProcessFuncBox() map[string]ProcessFunc {
	return processFuncBox
}
