package models

// `admin`
import (
	"github.com/beego/beego/v2/client/orm"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(Admin))
}

const ADMIN_TABLENAME string = "admin"

type Admin struct {
	Id            int64 `orm:"pk;"`
	Email         string
	Nickname      string
	Password      string
	Status        types.StatusBlockEnum
	RegisterTime  int64 `orm:"column(register_time)"`
	LastLoginTime int64 `orm:"column(last_login_time)"`
}

// 此处声明为指针方法,并不会修改传入的对象,只是为了省去拷贝对象的开消

// 当前模型对应的表名
func (r *Admin) TableName() string {
	return ADMIN_TABLENAME
}
