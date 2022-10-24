package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(SatisSurvey))
}

const SATIS_SURVEY_TABLENAME = "satis_survey_tab"

type SatisSurvey struct {
	ID        int64     `orm:"pk;column(id)" json:"id"` // 主键ID
	UserId    int64     `json:"user_id"`
	CourseId  int64     `json:"course_id"`
	Q1        int       `json:"q1"`
	Q2        int       `json:"q2"`
	Q3        int       `json:"q3"`
	Q4        int       `json:"q4"`
	Q5        int       `json:"q5"`
	Q6        int       `json:"q6"`
	Q7        int       `json:"q7"`
	Q8        int       `json:"q8"`
	Q9        int       `json:"q9"`
	Q10       int       `json:"q10"`
	UpdatedAt time.Time `orm:"type(datetime);precision(3)" json:"updated_at"`
	CreatedAt time.Time `orm:"type(datetime);precision(3)" json:"created_at"` // 创建时间

}

func (r *SatisSurvey) TableName() string {
	return SATIS_SURVEY_TABLENAME
}
