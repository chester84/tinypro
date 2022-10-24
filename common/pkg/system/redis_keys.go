package system

import "fmt"

const (
	rdsKeyPvPrefix    = `hm:stat:pv`     // 访问次数,记数器
	rdsKeyUvBoxPrefix = `hm:stat:uv-box` // 记录去重的人数

	rdsKey2Keep = `hm:stat:2-keep` // 次留
	rdsKey7Keep = `hm:stat:7-keep` // 7留

	rdsKeyRiskSamePlacePrefix = `hm:hash:risk-rules:sp`
)

func BuildPvKey(date string) string {
	return fmt.Sprintf(`%s:%s`, rdsKeyPvPrefix, date)
}

func BuildUvBoxKey(date string) string {
	return fmt.Sprintf(`%s:%s`, rdsKeyUvBoxPrefix, date)
}

func Build2KeepKey(date string) string {
	return fmt.Sprintf(`%s:%s`, rdsKey2Keep, date)
}

func Build7KeepKey(date string) string {
	return fmt.Sprintf(`%s:%s`, rdsKey7Keep, date)
}

func BuildRiskSamePlace(userId int64) string {
	return fmt.Sprintf(`%s:%d`, rdsKeyRiskSamePlacePrefix, userId)
}
