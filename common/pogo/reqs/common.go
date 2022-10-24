package reqs

type PageInfo struct {
	// type为0，表示要最新数据
	// type为1，表示要老数据
	Type int `json:"type"`
	Size int `json:"size"`
	SN   int `json:"sn"`
}

// 镜像是用weight值
type PageSelectedInfo struct {
	// type为0，表示要最新数据
	// type为1，表示要老数据
	Type           int `json:"type"`
	Size           int `json:"size"`
	SelectedWeight int `json:"selected_weight"`
}
