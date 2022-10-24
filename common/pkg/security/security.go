package security

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"tinypro/common/lib/redis/cache"
	"github.com/chester84/libtools"
)

func PassStrongPreventRepeatedEntry(router, traceID, ip string) bool {
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	originTraceID := traceID

	if len(traceID) != 32 {
		traceID = libtools.Md5(traceID)
	}

	ex := 600
	cKey := fmt.Sprintf(`%s:%s`, rdsRepeatedEntryPrefix, traceID)
	_, err := redis.String(cacheClient.Do("SET", cKey, ip, "EX", ex, "NX"))
	if err != nil {
		if err != redis.ErrNil {
			logs.Error("[PassStrongPreventRepeatedEntry] redis> SET %s %s EX %d NX, router: %s, traceID: %s, ip: %s, err: %v",
				cKey, ip, ex, router, traceID, ip, err)
		} else {
			logs.Warning("[PassStrongPreventRepeatedEntry] hit, ip: %s, router: %s, traceID: %s",
				ip, router, originTraceID)
		}

		return false
	}

	return true
}
