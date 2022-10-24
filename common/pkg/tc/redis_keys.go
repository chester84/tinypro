package tc

import "fmt"

const (
	rdsKeyCosTemporaryUrlPrefix = `tinypro:cache:cos-temporary-url`
)

func GenTemporaryUrlRdsKey(rid string) string {
	return fmt.Sprintf(`%s:%s`, rdsKeyCosTemporaryUrlPrefix, rid)
}
