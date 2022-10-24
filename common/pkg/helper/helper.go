package helper

import (
	strip "github.com/grokify/html-strip-tags-go"

	"tinypro/common/models"
	"tinypro/common/pkg/account"
	"tinypro/common/pkg/admin"
	"tinypro/common/types"
)

func OperatorName(pkID int64) string {
	if pkID <= 0 {
		return "-"
	}

	if 1 == pkID {
		return "admin"
	}

	bizSN, err := models.ParseBizSNFromPkID(pkID)
	if err != nil {
		return "-"
	}

	var name = ""

	switch bizSN {
	case types.AccountSystem:
		name = admin.OperatorName(pkID)

	case types.AppUserBiz:
		name = account.AppUserNickname(pkID)
	}

	return name
}

func LongStrDisplay(str string, length int) (show string) {
	if length < 8 || length > 32 {
		length = 8
	}

	// 过滤掉html标签,如果有的话,因为简要展示的地方本身也不需要html标签
	str = strip.StripTags(str)

	words := ([]rune)(str)
	wordsLen := len(words)
	if wordsLen < length {
		show = str
	} else {
		show = string(words[0:length]) + "..."
	}

	return
}
