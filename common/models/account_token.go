package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/lib/device"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(AccountToken))
}

const (
	ACCOUNT_TOKEN_TABLENAME string = "account_token"

	tokenExpire int64 = 2592000000 * 12 // 1年,毫秒数
)

type AccountToken struct {
	Id          int64  `orm:"pk;"`
	AccountId   int64  `orm:"column(account_id)"`
	AccessToken string `orm:"column(access_token)"`
	TokenIp     string `orm:"column(token_ip)"`
	Expires     int64
	Status      types.StatusCommonEnum
	Platform    string
	Ctime       int64
	Utime       int64
}

func (r *AccountToken) TableName() string {
	return ACCOUNT_TOKEN_TABLENAME
}

func GenerateAccountToken(accountId int64, platform string, ip string) (string, error) {
	bizId, _ := device.GenerateBizId(types.AccessTokenBiz)
	accessToken := libtools.Md5(fmt.Sprintf("%dchester@gmail.com%d@%s", bizId, time.Now().UnixNano(), platform))

	var expires = libtools.GetUnixMillis() + types.DayMillisecond
	switch platform {
	case types.PlatformH5:
		expires = libtools.GetUnixMillis() + types.DayMillisecond*30

	case types.PlatformWxMiniProgram:
		expires = libtools.GetUnixMillis() + types.DayMillisecond*30

	case types.PlatformAndroid, types.PlatformWatch:
		expires = libtools.GetUnixMillis() + tokenExpire
	}

	atIns := AccountToken{
		AccountId:   accountId,
		AccessToken: accessToken,
		TokenIp:     ip,
		Expires:     expires,
		Status:      types.StatusValid,
		Platform:    platform,
		Ctime:       libtools.GetUnixMillis(),
		Utime:       libtools.GetUnixMillis(),
	}

	o := orm.NewOrm()
	_, err := o.Insert(&atIns)

	return accessToken, err
}

func GetValidTokenByAccountId(accountId int64, platform string) (string, error) {
	var atIns = AccountToken{}
	o := orm.NewOrm()
	err := o.QueryTable(atIns.TableName()).
		Filter("account_id", accountId).
		Filter("platform", platform).
		Filter("status", types.StatusValid).
		OrderBy("-id").
		Limit(1).
		One(&atIns)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[GetValidTokenByAccountId] sql error err:%v", err)
		return "", err
	}

	if atIns.Expires < libtools.GetUnixMillis() {
		atIns.Status = types.StatusInvalid
		_, _ = OrmUpdate(&atIns, []string{"status"})
		return "", nil
	}

	return atIns.AccessToken, nil
}

func GetAccessTokenInfo(token string) (AccountToken, error) {
	var atIns = AccountToken{}
	o := orm.NewOrm()
	err := o.QueryTable(atIns.TableName()).Filter("access_token", token).One(&atIns)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[GetAccessTokenInfo] sql error err:%v", err)
		return atIns, err
	}

	return atIns, nil
}

func UpdateAccessTokenStatusByAccountId(accountId int64, status types.StatusCommonEnum) error {
	var atIns = AccountToken{}

	expires := libtools.GetUnixMillis()
	if status == types.StatusValid {
		expires = libtools.GetUnixMillis() + tokenExpire
	}

	o := orm.NewOrm()
	_, err := o.QueryTable(atIns.TableName()).Filter("account_id", accountId).Update(map[string]interface{}{
		"status":  status,
		"expires": expires,
		"utime":   libtools.GetUnixMillis(),
	})

	return err
}

// 账户下有效的token集合
func AccountValidToken(accountId int64) (int64, []AccountToken, error) {
	var list []AccountToken
	var m AccountToken

	// 构建查询对象
	qb, _ := orm.NewQueryBuilder(libtools.DBDriver())
	qb.Select("*").
		From(m.TableName()).
		Where("account_id = ? AND status = ? AND expires > ?").
		OrderBy("id").
		Desc()

	// 导出 SQL 语句
	sql := qb.String()

	// 执行 SQL 语句
	o := orm.NewOrm()
	num, err := o.Raw(sql, accountId, types.StatusValid, libtools.GetUnixMillis()).QueryRows(&list)

	return num, list, err
}
