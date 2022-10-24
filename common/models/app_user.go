package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(AppUser))
}

const APP_USER_TABLENAME = "app_user"

type AppUser struct {
	Id       int64 `orm:"pk;"`
	Nickname string
	Email    string
	Mobile   string
	Password string

	OpenOauthPlt int
	OpenToken    string
	OpenUserID   string `orm:"column(open_user_id)"`
	OpenAvatar   string
	OpenOauthMd5 string
	Birthday     string
	Gender       types.GenderEnum
	City         string

	RealName        string // 实名认证
	Company         string
	Position        string
	Residence       string
	RecommendPerson string
	RegisterAt      time.Time
	Status          types.StatusCommonEnum
	Level           int
	Balance         int64
	TotalAmount     int64
	AppVersion      string
	WxOpenId        string
	WxCountry       string
	WxProvince      string
	WxCity          string

	ShareCode   int64
	InviterCode int64
	LastLoginAt time.Time
	LastLoginIP string `orm:"column(last_login_ip)"`
}

func (r *AppUser) TableName() string {
	return APP_USER_TABLENAME
}

func (r *AppUser) Age() int {
	var age = 0
	if r.Birthday != "" {
		exp := strings.Split(r.Birthday, "-")
		if len(exp) > 0 {
			birthYear, _ := libtools.Str2Int(exp[0])
			age = time.Now().Year() - birthYear
			if age < 0 {
				logs.Warning("[AppUser.Age] birthday data abnormal, user: %#v", *r)
				age = 0
			}
		}
	}

	return age
}

// 余额可为负
func (r *AppUser) ChangeBalance(amount int64) (int64, error) {
	o := orm.NewOrm()

	sql := fmt.Sprintf(`UPDATE %s SET balance = balance + (%d) WHERE id = %d`, r.TableName(), amount, r.Id)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		logs.Error("[AppUser->ChangeBalance] db update get exception, SQL: %s, err: %v", strings.Replace(sql, "\n", " ", -1), err)
	} else {
		r.Balance += amount
	}

	return r.Balance, err
}

// 消费总额是累加的正数
func (r *AppUser) AddAmount(amount int64) (int64, error) {
	var err error
	if amount < 0 {
		err = fmt.Errorf("[AppUser->AddAmount] param amount can not be negtive, amount is %d", amount)
		return r.TotalAmount, err
	}
	o := orm.NewOrm()

	sql := fmt.Sprintf(`UPDATE %s SET total_amount = total_amount + (%d) WHERE id = %d`, r.TableName(), amount, r.Id)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		logs.Error("[AppUser->AddAmount] db update get exception, SQL: %s, err: %v", strings.Replace(sql, "\n", " ", -1), err)
	} else {
		r.TotalAmount += amount
	}

	return r.TotalAmount, err
}

// 充值
func (r *AppUser) AddBalanceAmount(amount int64) (int64, int64, error) {
	var err error
	if amount < 0 {
		//充值记录不能为小数
		err = fmt.Errorf("[AppUser->ChangeBalanceAmount] param amount can not be negtive, amount is %d", amount)
		return r.Balance, r.TotalAmount, err
	}
	o := orm.NewOrm()

	sql := fmt.Sprintf(`UPDATE %s SET balance = balance + (%d), total_amount = total_amount + (%d) WHERE id = %d`, r.TableName(), amount, amount, r.Id)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		logs.Error("[AppUser->ChangeBalanceAmount] db update get exception, SQL: %s, err: %v", strings.Replace(sql, "\n", " ", -1), err)
	} else {
		r.Balance += amount
		r.TotalAmount += amount
	}
	return r.Balance, r.TotalAmount, err
}

func (r *AppUser) IsEmptyMobile() bool {
	if r.Mobile == "" {
		return true
	} else {
		return false
	}
}
