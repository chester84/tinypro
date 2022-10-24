package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(Enroll))
}

const ENROLL_TABLENAME = "enroll_tab"

// Course 结构体
type Enroll struct {
	ID              int64     `orm:"pk;column(id)" json:"id"` // 主键ID
	CreatedAt       time.Time // 创建时间
	UpdatedAt       time.Time // 更新时间
	UserID          int64     `orm:"column(user_id)" json:"user_id" form:"user_id"`
	CourseId        int64     `json:"course_id" form:"course_id"`
	IsPass          int       `json:"is_pass" form:"is_pass"`
	IsSignIn        int       `json:"is_sign_in" form:"is_sign_in"`
	Reason          string    `json:"reason" form:"reason"`
	RealName        string    `json:"real_name" form:"real_name"`
	Mobile          string    `json:"mobile" form:"mobile"`
	Company         string    `json:"company" form:"company"`
	Position        string    `json:"position" form:"position"`
	Residence       string    `json:"residence" form:"residence"`
	RecommendPerson string    `json:"recommend_person" form:"recommend_person"`
}

func (r *Enroll) TableName() string {
	return ENROLL_TABLENAME
}
