package admin

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"tinypro/common/lib/redis/cache"
	"tinypro/common/lib/redis/storage"
	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type AdmForDisplay struct {
	models.Admin
}

func Add(admin *models.Admin) (id int64, err error) {
	o := orm.NewOrm()

	id, err = o.Insert(admin)

	return
}

func UpdateStatus(adminId int64, status types.StatusBlockEnum) (num int64, err error) {
	if adminId <= 1 {
		err = errors.New("参数不正确")
		return
	}

	obj := models.Admin{
		Id:     adminId,
		Status: status,
	}
	o := orm.NewOrm()
	num, err = o.Update(&obj, "status")

	return
}

// Update 更新指定角色的属性
// 不含属性校验
// 内部自动更新 Utime
func Update(m *models.Admin, om *models.Admin, cols []string) (num int64, err error) {

	if m.Id <= 0 {
		err = errors.New("Update ID must exist and >0")
		return
	}

	o := orm.NewOrm()

	num, err = o.Update(m, cols...)
	if num > 0 {
		if m.Nickname != om.Nickname {
			ClearNameCache(m.Id)
		}
	}
	return
}

// GetNameByID 根据adminID 获取用户名
func GetNameByID(adminID int64) string {
	//logs.Debug("adminID: %d", adminID)
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	hashKey := rdsKeyOperatorName

	var name = `-`
	valueByte, err := storageClient.Do("HGET", hashKey, adminID)
	//logs.Debug("valueByte:", valueByte, ", err:", err)
	if err == nil && valueByte != nil {
		name = string(valueByte.([]byte))
	} else {
		admin, err := OneByUid(adminID)
		if err != nil {
			return "无效的操作员"
		}

		name = admin.Nickname
		_, _ = storageClient.Do("HSET", hashKey, adminID, name)
	}

	return name
}

// ClearNameCache 清除用户名缓存
func ClearNameCache(adminID int64) {
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()
	hashKey := rdsKeyOperatorName
	_, _ = storageClient.Do("HDEL", hashKey, adminID)
}

// OperatorName 取操作员的名字
func OperatorName(opUid int64) string {
	if 0 == opUid {
		return "-"
	}

	return GetNameByID(opUid)
}

func AddLoginLog(adminUID int64, ip string) (int64, error) {
	obj := models.AdminLoginLog{AdminUID: adminUID, IP: ip, Ctime: libtools.GetUnixMillis()}

	return models.OrmInsert(&obj)
}

func OneByUid(id int64) (models.Admin, error) {
	admin := &models.Admin{Id: id}
	o := orm.NewOrm()

	err := o.Read(admin)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[OneByUid] read one get exception, id: %d, err: %v", id, err)
	}

	return *admin, err
}

func OneByNickName(nickName string) (models.Admin, error) {
	var admin models.Admin

	o := orm.NewOrm()

	err := o.QueryTable(models.ADMIN_TABLENAME).Filter("nickname", nickName).One(&admin)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[OneAdminByNickName] sql error err:%v", err)
	}

	return admin, err
}

func OneByEmail(email string) (models.Admin, error) {
	var admin models.Admin

	o := orm.NewOrm()

	err := o.QueryTable(models.ADMIN_TABLENAME).Filter("email", email).One(&admin)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[OneAdminByEmail] sql error err:%v", err)
	}

	return admin, err
}

func CheckLoginIsValid(email string, password string) bool {
	admin, err := OneByEmail(email)
	//logs.Debug("admin:", admin)
	if err != nil || admin.Id <= 0 {
		logs.Warning("email and info does not exist:", email)
		return false
	}

	cipherText := libtools.PasswordEncrypt(password, admin.RegisterTime)
	if cipherText == admin.Password {
		return true
	}

	logs.Warning("User information is incorrect, email:", email)
	return false
}

func UpdateLastLoginTime(id int64) {
	admin := models.Admin{
		Id:            id,
		LastLoginTime: libtools.GetUnixMillis(),
	}
	o := orm.NewOrm()

	_, _ = o.Update(&admin, "last_login_time")
}

// 改
func UpdateAll(admin models.Admin) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Update(&admin)
	if err != nil {
		logs.Error("model Admin UpdateRepayPlan failed.", err)
	}

	return
}

func buildLoginCaptchaRdsKey(cookieValue string) string {
	return fmt.Sprintf(`%s:%s`, rdsKeyAdmLoginCaptcha, cookieValue)
}

func SetLoginCaptcha(cookieValue, captchaValue string) {
	rdsKey := buildLoginCaptchaRdsKey(cookieValue)

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	_, err := cacheClient.Do("SET", rdsKey, captchaValue, "EX", 600)
	if err != nil {
		logs.Error("[SetLoginCaptcha] redis> SET %s %s EX 600, err: %v", rdsKey, captchaValue, err)
	}
}

func VerifyLoginCaptcha(cookieValue, captchaValue string) bool {
	rdsKey := buildLoginCaptchaRdsKey(cookieValue)

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	cValue, err := redis.String(cacheClient.Do("GET", rdsKey))
	if err != nil {
		if err != redis.ErrNil {
			logs.Error("[VerifyLoginCaptcha] redis> GET %s, err: %v", rdsKey, err)
		} else {
			logs.Warning("[VerifyLoginCaptcha] get empty data, cookieValue, %s, captchaValue: %s, rdsKey: %s",
				cookieValue, captchaValue, rdsKey)
		}

		return false
	}

	if len(cValue) > 0 && cValue == captchaValue {
		return true
	}

	return false
}
