package types

const (
	HotSearchWordsKey = `hot_search_words`

	SpecialTopicContent = 104102
)

type EduChargeTypeEnum int

const (
	EduChargeTypeFixed   EduChargeTypeEnum = 1
	EduChargeTypeUnfixed EduChargeTypeEnum = 2
)

type HotSearchWordsConf struct {
	Zh string `json:"zh"`
	En string `json:"en"`
}

type EduChargeItemSimple struct {
	ItemSN string `json:"item_sn"`
	Name   string `json:"name"`
	Amount string `json:"amount,omitempty"`

	AmountNum int64  `json:"amount_num,omitempty"`
	Remark    string `json:"remark,omitempty"`
}

type EduPayCreateOrderReq struct {
	AppSN       int                   `json:"app_sn"`
	TradeType   WxPayTradeTypeEnum    `json:"trade_type"`
	ChargeSN    string                `json:"charge_sn"`
	TotalAmount int64                 `json:"total_amount"`
	FixedGroup  []EduChargeItemSimple `json:"fixed_group"`
	UnfixGroup  []EduChargeItemSimple `json:"unfix_group"`
	Remark      string                `json:"remark,omitempty"` // 总备注
}

// 此类内容没有头图和标题摘要
func LandingPageConf() map[int]string {
	return map[int]string{
		// 学校概况 100
		100001: `松雷教育`,
		100002: `校董会/领导团队`,
		100003: `办学`,
		100004: `学校历史`,
		100005: `管理架构`,

		// 松雷人 101
		101001: `寄语`,
		101002: `教师团队`,

		// 新闻 102
		102001: `党团建设党委办公室`,
		102002: `党团建设工会`,
		102003: `党团建设团委`,
		102004: `党团建设少先队`,

		// 松雷特色 103
		103001: `教育理念与宗旨`,
		103002: `德育之窗`,
		103003: `身心健康`,
		103004: `升学指导`,

		// 教学 104
		104000: `教师团队`,

		// 资料 105
		105000: `资料`,

		// 校服 106
		106000: `校服`,

		// 环境 107
		107000: `环境`,

		// 餐饮 108
		108000: `餐饮`,

		// 师资 109
		109000: `师资`,

		// 地图 111
		111000: `地图`,

		// 活动 113
		113101: `走进松雷`,
	}
}

// 呈现在列表页,头图,标题,摘要齐全
func StdContentConf() map[int]string {
	return map[int]string{
		// 学校概况
		100101: `校园分层`,
		100102: `名师风采`,

		// 松雷人
		101101: `校友风采`,

		// 新闻
		102101: `学校新闻`,
		102102: `通知公告`,
		102103: `媒体报道`,

		// 松雷特色
		103101: `课程设置`,
		103102: `学生活动`,
		103103: `在线学习`,

		// 活动 110
		110000: `活动`,

		// 消息 112
		112101: `消息`,

		// 最新资讯 114
		114000: `最新资讯`,
	}
}

// 这一类是热区配置
func HotZoneMapConf() map[int]string {
	return map[int]string{
		200000: `报名`,
		400000: `联系我们`,
		500000: `直播`,
	}
}

// 专题类
func SpecialTopicConf() map[int]string {
	return map[int]string{
		// 教学
		104101: `课程概述`,
	}
}

func Mark2Label(markNum int) string {
	if label, ok := LandingPageConf()[markNum]; ok {
		return label
	}

	if label, ok := StdContentConf()[markNum]; ok {
		return label
	}

	if label, ok := HotZoneMapConf()[markNum]; ok {
		return label
	}

	if label, ok := SpecialTopicConf()[markNum]; ok {
		return label
	}

	return `-`
}
