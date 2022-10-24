package accesstoken

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/lib/redis/cache"
	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

// 缓存相关的 token

func buildTokeCacheKey(platform, token string) (key string) {
	switch platform {
	case types.PlatformH5:
		key = fmt.Sprintf("%s:%s", rdsKeyWebApiTokenPrefix, token)
	case types.PlatformWxMiniProgram:
		key = fmt.Sprintf("%s:%s", rdsKeyMiniProgramApiTokenPrefix, token)
	case types.PlatformAdm:
		key = fmt.Sprintf("%s:%s", rdsKeyWebAdmTokenPrefix, token)
	case types.PlatformWatch:
		key = fmt.Sprintf(`%s:%s`, rdsKeyWatchTokenPrefix, token)
	default:
		key = fmt.Sprintf("%s:%s", rdsKeyAppApiTokenPrefix, token)
	}

	return
}

// 调用方只用关心是否有效,不用关心具体原因
func IsValidAccessToken(platform, token string) (bool, int64) {
	if token == "" {
		logs.Warning("access token is empty")
		return false, 0
	}

	cKey := buildTokeCacheKey(platform, token)
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	cValue, err := cacheClient.Do("GET", cKey)
	//logs.Debug("cValue:", cValue, ", err:", err)
	if err != nil || cValue == nil {
		logs.Warning("access token key DOES NOT EXISTS, cKey:", cKey)
		return false, 0
	}

	var tokenInfo models.AccountToken
	err = json.Unmarshal(cValue.([]byte), &tokenInfo)
	if err != nil {
		// 说明有缓存数据,但内容有问题,消除之
		CleanTokenCache(platform, token)
		logs.Warning("json decode has wrong, please checkout. cKey:", cKey, ", cValue:", string(cValue.([]byte)))
		return false, 0
	}

	if tokenInfo.AccountId <= 0 || tokenInfo.Status != types.StatusValid || tokenInfo.Expires < libtools.GetUnixMillis() {
		// 无效数据,消除之
		CleanTokenCache(platform, token)
		if tokenInfo.Expires < libtools.GetUnixMillis() {
			// 如果过期了，把状态置为失效
			tokenInfo.Status = types.StatusInvalid
			_, _ = models.OrmUpdate(&tokenInfo, []string{"status"})
		}
		logs.Warning("cache data is invalid, please checkout. cKey:", cKey, ", cValue:", string(cValue.([]byte)))
		return false, 0
	}

	return true, tokenInfo.AccountId
}

func GenTokenWithCache(accountId int64, platform string, ip string) (string, error) {
	if libtools.IsProductEnv() {
		// 登陆新设置,踢掉其他设置上的登陆态
		kickOffOtherToken(accountId, platform)
	}

	// 先尝试查询，看是否存在有效token
	token, err := models.GetValidTokenByAccountId(accountId, platform)
	if err != nil {
		logs.Error("GetValidTokenByAccountId err. accountId:", accountId, ", platform:", platform, ", ip:", ip, "err %#v", err)
		return "", err
	}

	if token != "" {
		return token, nil
	}

	// 没有有效token，再创建token
	token, err = models.GenerateAccountToken(accountId, platform, ip)
	if err != nil {
		logs.Error("can NOT create account_token. accountId:", accountId, ", platform:", platform, ", ip:", ip)
		return "", err
	}

	tokenInfo, err := models.GetAccessTokenInfo(token)
	if err != nil {
		logs.Error("can NOT find token info. token:", token)
		return "", err
	}

	cKey := buildTokeCacheKey(platform, token)
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	bson, err := json.Marshal(tokenInfo)
	expires := tokenInfo.Expires - libtools.GetUnixMillis()
	if err != nil || expires < 0 {
		logs.Error("json encode model data has wrong. tokenInfo:", tokenInfo, ", expires:", expires)
		return "", err
	}

	// 安全忽略
	_, _ = cacheClient.Do("SET", cKey, string(bson), "PX", expires)

	return token, nil
}

func kickOffOtherToken(accountID int64, platform string) {
	_, tokenBox, _ := models.AccountValidToken(accountID)
	for _, token := range tokenBox {
		// 有可能多端都要相互踢!!!
		if token.Platform != platform {
			continue
		}

		CleanTokenCache(platform, token.AccessToken)

		token.Status = types.StatusInvalid
		token.Utime = libtools.GetUnixMillis()
		_, err := models.OrmUpdate(&token, []string{"Status", "Utime"})
		if err != nil {
			logs.Error("[kickOffOtherToken] update exception, data: %#v, err: %v", token, err)
		}
	}
}

// 登出操作
func CleanTokenCache(platform, token string) {
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	cKey := buildTokeCacheKey(platform, token)
	// 理论不会出错,直接忽略返回
	_, err := cacheClient.Do("DEL", cKey)
	if err != nil {
		logs.Error("[CleanTokenCache] redis> DEL %s, err: %v", cKey, err)
	}
}

func GetUserIdByToken(token string) (obj models.AccountToken, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(obj.TableName()).
		Filter("access_token", token).
		OrderBy("-id").
		Limit(1).
		One(&obj)

	return
}
