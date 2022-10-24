package weixin

import (
	"fmt"

	"github.com/chester84/libtools"
)

func BuildShareData(sceneNum int) AppMessageShareData {
	var data AppMessageShareData

	data.Title = `亿车互助`
	data.Desc = `一车出事故,众车来互助`
	data.Link = libtools.InternalH5Domain()
	data.ImgUrl = `https://cdn-1302993108.cos.ap-guangzhou.myqcloud.com/img/h5-03/0.png`

	switch sceneNum {

	case 1: // 公示
		data.Title = `亿车互助 - 公示`
		data.Link = fmt.Sprintf(`%s/?p=1`, libtools.InternalH5Domain())

	case 2: // 延保活动
		data.Title = `EV互助 - 零部件延保`
		data.Desc = `一车出故障,众车来互助`
		data.ImgUrl = `https://cdn-1302993108.cos.ap-guangzhou.myqcloud.com/img/ev-share.jpg`
		data.Link = fmt.Sprintf(`%s/?activity=yb-activity-2020`, libtools.InternalH5Domain())
	}

	return data
}
