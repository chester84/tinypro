package types

type HighChartsSeries struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

type HighChartsSpLine struct {
	XAxis  []string           `json:"xAxis"`
	Series []HighChartsSeries `json:"series"`
}

// 只有两组数据的折线图
type VueEchartsLineG2 struct {
	G1 []interface{} `json:"g1"`
	G2 []interface{} `json:"g2"`
}
