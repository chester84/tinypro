package models

// `admin`
import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(AdvertisingPopup))
}

const ADVERTISINGPOPUP_TABLENAME string = "advertising_popup_tab"

type AdvertisingPopup struct {
	ID        int64     `orm:"pk;column(id)" json:"id"` // 主键ID
	CourseId  int64     `json:"course_id" form:"course_id"`
	Switch    int       `json:"switch" form:"switch"`
	Name      string    `json:"name" form:"name"`
	Url       string    `json:"url" form:"url"`
	GoToUrl   string    `json:"go_to_url" form:"go_to_url"`
	CreatedAt time.Time `orm:"type(datetime);precision(3)" json:"created_at"` // 创建时间
	UpdatedAt time.Time // 更新时间

}

// 当前模型对应的表名
func (r *AdvertisingPopup) TableName() string {
	return ADVERTISINGPOPUP_TABLENAME
}
