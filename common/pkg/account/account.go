package account

import (
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"tinypro/common/cerror"
	"tinypro/common/lib/redis/cache"
	"tinypro/common/lib/redis/storage"
	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

const (
	guestNicknamePrefix = `Guest`
)

func buildOpenOauthMd5(plt int, openUid string) string {
	return libtools.Md5(fmt.Sprintf(`%d$%s`, plt, openUid))
}

func LoginOrRegister(req types.ApiOauthLoginReqT, ip, appVersion string) (models.AppUser, error) {
	var user models.AppUser
	var err error

	logs.Debug("LoginOrRegister req %#v", req)

	if req.OpenUserID == "" || req.Nickname == "" || req.OpenOauthPlt < types.OpenOauthWeChat {
		err = fmt.Errorf("get empty user_id or name")
		return user, err
	}

	var oauthMd5 = buildOpenOauthMd5(req.OpenOauthPlt, req.OpenUserID)
	one, errQ := tryFindOneByCond(req.WxOpenId, oauthMd5)
	if errQ != nil {
		if errQ.Error() == types.EmptyOrmStr {
			// 新用户,同时传过来了mobile不为空，自动注册
			if req.Mobile != "" {
				city := libtools.GeoIpCityZhCN(ip)
				if city == libtools.EmptyRecord {
					city = ""
				}

				// 之前device.GenerateBizId()
				// 这个方法有可能引起redis变慢

				u := models.AppUser{
					Nickname:     req.Nickname,
					OpenUserID:   req.OpenUserID,
					OpenToken:    req.OpenToken,
					OpenOauthPlt: req.OpenOauthPlt,
					OpenAvatar:   req.OpenAvatar,
					OpenOauthMd5: oauthMd5,
					Mobile:       req.Mobile,
					Email:        req.Email,
					Gender:       req.Gender,
					City:         city,
					RegisterAt:   time.Now(),
					LastLoginAt:  time.Now(),
					LastLoginIP:  ip,
					Status:       types.StatusValid,
					Birthday:     "1990-01-01",
					ShareCode:    GenerateInviteCode(),
					AppVersion:   appVersion,
					WxOpenId:     req.WxOpenId,
					WxCountry:    req.Country,
					WxProvince:   req.Province,
					WxCity:       req.City,
				}

				if req.OpenOauthPlt == types.OpenOauthGuest {
					u.Nickname = GenGuestNickname()
					u.Gender = types.GenderUnknown
				}

				pkId, errI := models.OrmInsert(&u)
				if errI != nil {
					err = errI
				} else {
					user = u
					user.Id = pkId
				}
			}

		} else {
			// 到此分支应该是异常
			err = errQ
		}
	} else {
		origin := one
		// 更新老用户最后登陆数据,历史登陆数据呢? TODO
		one.LastLoginAt = time.Now()
		one.LastLoginIP = ip
		upCols := []string{
			"last_login_at", "last_login_ip",
		}

		var hasChange bool
		if req.Nickname != "" && req.Nickname != one.Nickname && !strings.HasPrefix(req.Nickname, guestNicknamePrefix) {
			one.Nickname = req.Nickname
			upCols = append(upCols, "nickname")
			hasChange = true
		}
		if req.OpenAvatar != "" && req.OpenAvatar != one.OpenAvatar {
			one.OpenAvatar = req.OpenAvatar
			upCols = append(upCols, "open_avatar")
			hasChange = true
		}
		if req.Country != "" && req.Country != one.WxCountry {
			one.WxCountry = req.Country
			upCols = append(upCols, "wx_country")
			hasChange = true
		}
		if req.Province != "" && req.Province != one.WxProvince {
			one.WxProvince = req.Province
			upCols = append(upCols, "wx_province")
			hasChange = true
		}
		if req.City != "" && req.City != one.WxCity {
			one.WxCity = req.City
			upCols = append(upCols, "wx_city")
			hasChange = true
		}
		if req.Mobile != "" && req.Mobile != one.Mobile {
			one.Mobile = req.Mobile
			upCols = append(upCols, "mobile")
			hasChange = true
		}

		if req.WxOpenId != "" && req.WxOpenId != one.WxOpenId {
			one.WxOpenId = req.WxOpenId
			upCols = append(upCols, "wx_open_id")
			hasChange = true
		}

		logs.Debug("RegisterOrLogin debug req !!!!! %#v", req)

		_, errU := models.OrmUpdate(&one, upCols)
		if errU != nil {
			err = errU
			logs.Error("[RegisterOrLogin] update exception, one: %#v, err: %v", one, err)
		}

		//logs.Debug("hasChange: %#v", hasChange)
		if hasChange {
			models.OpLogWrite(one.Id, one.Id, models.OpCodeUpAppUser, one.TableName(), origin, one)
		}

		user = one
	}

	WriteWxOpenId(user.Id, req.AppSN, req.WxOpenId)

	return user, err
}

func RegisterOrLogin(req types.ApiOauthLoginReqT, ip, appVersion string) (models.AppUser, error) {
	var user models.AppUser
	var err error

	if req.OpenUserID == "" || req.Nickname == "" || req.OpenOauthPlt < types.OpenOauthWeChat {
		err = fmt.Errorf("get empty user_id or name")
		return user, err
	}

	var oauthMd5 = buildOpenOauthMd5(req.OpenOauthPlt, req.OpenUserID)
	one, errQ := tryFindOneByCond(req.WxOpenId, oauthMd5)
	if errQ != nil {
		if errQ.Error() == types.EmptyOrmStr {
			// 新用户,自动注册
			//id, _ := device.GenerateBizId(types.AppUserBiz)
			//目前不需要这种redis生成主键id
			city := libtools.GeoIpCityZhCN(ip)
			if city == libtools.EmptyRecord {
				city = ""
			}

			u := models.AppUser{
				Nickname:     req.Nickname,
				OpenUserID:   req.OpenUserID,
				OpenToken:    req.OpenToken,
				OpenOauthPlt: req.OpenOauthPlt,
				OpenAvatar:   req.OpenAvatar,
				OpenOauthMd5: oauthMd5,
				Email:        req.Email,
				Gender:       req.Gender,
				City:         city,
				RegisterAt:   time.Now(),
				LastLoginAt:  time.Now(),
				LastLoginIP:  ip,
				Status:       types.StatusValid,
				Birthday:     "1990-01-01",
				ShareCode:    GenerateInviteCode(),
				AppVersion:   appVersion,
				WxCountry:    req.Country,
				WxProvince:   req.Province,
				WxCity:       req.City,
			}

			if req.OpenOauthPlt == types.OpenOauthGuest {
				u.Nickname = GenGuestNickname()
				u.Gender = types.GenderUnknown
			}

			pkId, errI := models.OrmInsert(&u)
			if errI != nil {
				err = errI
			} else {
				u.Id = pkId
				user = u
			}
		} else {
			// 到此分支应该是异常
			err = errQ
		}

	} else {
		origin := one
		// 更新老用户最后登陆数据,历史登陆数据呢? TODO
		one.LastLoginAt = time.Now()
		one.LastLoginIP = ip
		upCols := []string{
			"last_login_at", "last_login_ip",
		}

		var hasChange bool
		if req.Nickname != "" && req.Nickname != one.Nickname && !strings.HasPrefix(req.Nickname, guestNicknamePrefix) {
			one.Nickname = req.Nickname
			upCols = append(upCols, "nickname")
			hasChange = true
		}
		if req.OpenAvatar != "" && req.OpenAvatar != one.OpenAvatar {
			one.OpenAvatar = req.OpenAvatar
			upCols = append(upCols, "open_avatar")
			hasChange = true
		}
		if req.Country != "" && req.Country != one.WxCountry {
			one.WxCountry = req.Country
			upCols = append(upCols, "wx_country")
			hasChange = true
		}
		if req.Province != "" && req.Province != one.WxProvince {
			one.WxProvince = req.Province
			upCols = append(upCols, "wx_province")
			hasChange = true
		}
		if req.City != "" && req.City != one.WxCity {
			one.WxCity = req.City
			upCols = append(upCols, "wx_city")
			hasChange = true
		}

		_, errU := models.OrmUpdate(&one, upCols)
		if errU != nil {
			err = errU
			logs.Error("[RegisterOrLogin] update exception, one: %#v, err: %v", one, err)
		}

		//logs.Debug("hasChange: %#v", hasChange)
		if hasChange {
			models.OpLogWrite(one.Id, one.Id, models.OpCodeUpAppUser, one.TableName(), origin, one)
		}

		user = one
	}

	WriteWxOpenId(user.Id, req.AppSN, req.WxOpenId)

	return user, err
}

func tryFindOneByCond(wxOpenId, openOauthMd5 string) (models.AppUser, error) {
	var obj = models.AppUser{}
	var err error

	o := orm.NewOrm()

	err = o.QueryTable(obj.TableName()).
		Filter("open_oauth_md5", openOauthMd5).
		One(&obj)
	if err != nil {
		if err != orm.ErrNoRows {
			logs.Error("[tryFindOneByCond] db exception, openOauthMd5: %s,  err: %v", openOauthMd5, err)
		}

		wxOpenIdObj := models.WxOpenId{}
		err = o.QueryTable(wxOpenIdObj.TableName()).Filter("open_id", wxOpenId).One(&wxOpenIdObj)
		if err != nil {
			if err != orm.ErrNoRows {
				logs.Error("[tryFindOneByCond] db exception, wxOpenId: %s,  err: %v", wxOpenId, err)
			}

			return obj, err
		}

		err = models.OrmOneByPkId(wxOpenIdObj.UserId, &obj)
		if err != nil {
			logs.Warning("[tryFindOneByCond] get one user exception, userId: %d,  err: %v", wxOpenIdObj.UserId, err)
		} else {
			// OrmOneByPkId升级过，
			// 此处需要从新赋一个 orm.ErrNoRows
			err = orm.ErrNoRows
		}

	}

	return obj, err
}

func OneAppUserByMobile(mobile string) (models.AppUser, error) {
	var obj = models.AppUser{}
	var err error

	if mobile == "" {
		err = fmt.Errorf(`[OneAppUserByMobile] intput mobile is empty`)
		return obj, err
	}

	o := orm.NewOrm()

	err = o.QueryTable(obj.TableName()).Filter("mobile", mobile).One(&obj)

	return obj, err
}

func OneAppUserByWatchDeviceId(deviceId int64) (models.AppUser, error) {
	var obj = models.AppUser{}
	var err error

	if deviceId == 0 {
		err = fmt.Errorf(`[OneAppUserByWatchDeviceId] intput deviceId is empty`)
		return obj, err
	}

	o := orm.NewOrm()

	err = o.QueryTable(obj.TableName()).Filter("watch_device_id", deviceId).One(&obj)

	return obj, err
}

func ResetAppUserNickname(pkID int64, nickname string) {
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	hashKey := rdsKeyAppUserNickname

	if pkID > 0 && nickname != "" {
		_, err := storageClient.Do("HSET", hashKey, pkID, nickname)
		if err != nil {
			logs.Error("[ResetAppUserNickname] redis> HSET %s %d %s, err: %v", hashKey, pkID, nickname)
		}
	}
}

func AppUserNickname(pkID int64) string {
	if pkID <= 0 {
		return "-"
	}

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	hashKey := rdsKeyAppUserNickname

	var name string
	value, err := redis.String(storageClient.Do("HGET", hashKey, pkID))
	if err != nil && err != redis.ErrNil {
		logs.Error("[AppUserNickname] redis get exception, redis> HGET %s %d, err: %v", hashKey, pkID, err)
	}

	if value != "" {
		name = value
		return name
	} else {
		var one models.AppUser
		err := models.OrmOneByPkId(pkID, &one)
		if err != nil {
			logs.Warning("[AppUserNickname] can find data, pkID: %d, err: %v", pkID, err)
			return "-"
		}

		name = one.Nickname
		_, err = storageClient.Do("HSET", hashKey, pkID, name)
		if err != nil {
			logs.Error("[AppUserNickname] redis get exception, redis> HSET %s %d %s, err: %v", hashKey, pkID, name, err)
		}
	}

	if name == "" {
		name = "-"
	}

	return name
}

func CheckIsSignIn(accountID int64) bool {
	if accountID <= 0 {
		logs.Error("input is 0")
		return false
	}

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	timeNow := libtools.GetUnixMillis()
	ex := 3600 * 24

	cKey := fmt.Sprintf(`%s:%s:%d`, rdsKeySignInPrefix,
		libtools.UnixMsec2Date(timeNow, "Y-m-d"), accountID)

	_, err := redis.String(cacheClient.Do("SET", cKey, timeNow, "EX", ex, "NX"))
	if err != nil {
		logs.Warning("[CheckIsSignIn] redis> SET %s %d EX %d NX, err: %v", cKey, timeNow, ex, err)
		if err.Error() == types.RedigoNilReturned {
			return true
		}
	}

	return false
}

func BindMobile(accountId int64, mobile, inviterCode string) cerror.ErrCode {
	var (
		err     error
		oneUser models.AppUser
	)

	inviterCodeI, _ := libtools.Str2Int64(inviterCode)

	_, err = OneAppUserByMobile(mobile)
	if err != nil {
		err = models.OrmOneByPkId(accountId, &oneUser)
		if err != nil {
			logs.Error("[BindMobile] oneUser does not exist, err %#v", err)
			return cerror.InvalidAccount
		}

		oneUser.Mobile = mobile
		// 1. 只生效一次; 2. 不能是自己
		if oneUser.InviterCode == 0 && oneUser.ShareCode != inviterCodeI {
			oneUser.InviterCode = inviterCodeI
		}
		_, errU := models.OrmUpdate(&oneUser, []string{"mobile", "InviterCode"})
		if errU != nil {
			logs.Error("[BindMobile] update oneUser err %#v", errU)
			return cerror.ServiceDbOpFail
		}
	} else {
		return cerror.MobileAlreadyBind
	}

	return cerror.CodeSuccess
}

func GenGuestNickname() string {
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	// 新自动登陆游客的昵称
	num := libtools.GenerateRandom(1, 6)
	guestNum, err := redis.Int64(storageClient.Do("INCRBY", rdsKeyOauthGuest, num))
	if err != nil {
		logs.Error("[genGuestNickname] redis> INCRBY %s %d, err: %v", rdsKeyOauthGuest, num, err)
	}
	return fmt.Sprintf(`%s-%d`, guestNicknamePrefix, guestNum)
}

func InviteFriends(inviterCode int64) (list []models.AppUser, err error) {
	var obj = models.AppUser{}

	if inviterCode <= 0 {
		err = fmt.Errorf(`[InviteFriends] intput inviterCode is empty`)
		return
	}

	o := orm.NewOrm()

	_, err = o.QueryTable(obj.TableName()).
		Filter("inviter_code", inviterCode).
		All(&list)

	return
}

func GenerateInviteCode() int64 {
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	setKey := rdsKeyInviteCodePool

	setCount, err := redis.Int(storageClient.Do("SCARD", setKey))
	if err != nil && err != redis.ErrNil {
		logs.Error("[GenOneInviteCode] redis> SCARD %s, err: %v", setKey, err)
		return 0
	}

	if setCount == 0 {
		for i := 10000; i <= 9999999; i++ {
			_, err := storageClient.Do("SADD", setKey, i)
			if err != nil {
				logs.Error("[GenOneInviteCode] redis> SADD %s %d, err: %v", setKey, i, err)
				return 0
			}
		}
	}

	inviteSN, err := redis.String(storageClient.Do("SPOP", setKey))
	if err != nil {
		logs.Error("[GenOneInviteCode] redis> SPOP %s, err: %v", setKey, err)
		return 0
	}

	inviteCode, _ := libtools.Str2Int64(inviteSN)
	return inviteCode
}
