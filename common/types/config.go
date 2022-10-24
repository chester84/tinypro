package types

const (
	MnpReleaseCtlKey = `mnp_release_ctl`
)

type MnpReleaseCtlConf struct {
	AppVersion string `json:"app_version"` // 4位字符串版本 1.1234.1234.1234
	NumVersion int64  `json:"num_version"`

	Notice string `json:"notice"` // 提示内容

	LastOpBy int64 `json:"last_op_by"`
	LastOpAt int64 `json:"last_op_at"`
}
