package reqs

type WxLoginReqT struct {
	Code    string `json:"code"`
	AppSN   int    `json:"app_sn"`
	Reginfo struct {
		Avatar   string `json:"avatar"`
		Username string `json:"username"`
		Mobile   string `json:"mobile"`
	} `json:"reginfo"`
}
