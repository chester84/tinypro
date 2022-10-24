package i18n

import (
	"fmt"

	"tinypro/common/types"
)

const rdsKeyStringMapping = "tinypro:hash:string-mapping"

func stringMappingKey(src string, enum types.LanguageTypeEnum) string {
	return fmt.Sprintf("%s_%d", src, enum)
}
