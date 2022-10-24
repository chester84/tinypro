package advertisingpopup

import (
	"fmt"
	"github.com/chester84/libtools"
)

const (
	rdsKeyAccessPopup = `tinypro:string:access-popup`
)

func GetAccessPopupRdsKey(advertiseId int64, userId int64) (rdsKey string) {
	todayUnix := libtools.NaturalDay(0)
	today := libtools.UnixMsec2Date(todayUnix, "Y-m-d")
	rdsKey = fmt.Sprintf("%s:%d:%d:%s", rdsKeyAccessPopup, advertiseId, userId, today)
	return
}
