package types

const (
	EsRepairFactoryGeo   = `repair_factory-geo`
	EsRepairFactoryGeoV1 = `repair_factory-geo-v1`
	EsRepairFactoryGeoV2 = `repair_factory-geo-v2`
)

type EsRepairFactoryGeoItem struct {
	SN       string `json:"sn"`
	Location string `json:"location"`
}

type ApiNearbyReqT struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
	Radius    int    `json:"radius"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	Exclude   int    `json:"exclude"`
}

type ApiRepairFactoryItem struct {
	SN           string `json:"sn"`
	Name         string `json:"name"`
	CoverImg     string `json:"cover_img"` // 封面图
	Publicity    string `json:"publicity"`
	ServicePhone string `json:"service_phone"` // 服务电话
	ServiceTime  string `json:"service_time"`  // 服务时间
	Address      string `json:"address"`       // 地址
	Detail       string `json:"detail"`        // 详细说明
	Distance     string `json:"distance"`      // 距离

	Tags []TagTwoTupleS `json:"tags"`
}

type ApiEv4sItem struct {
	SN           string `json:"sn"`
	Name         string `json:"name"`
	CoverImg     string `json:"cover_img"` // 封面图
	Publicity    string `json:"publicity"`
	ServicePhone string `json:"service_phone"` // 服务电话
	ServiceTime  string `json:"service_time"`  // 服务时间
	Address      string `json:"address"`       // 地址
	Detail       string `json:"detail"`        // 详细说明
	Distance     string `json:"distance"`      // 距离

	Tags []TagTwoTupleS `json:"tags"`
}
