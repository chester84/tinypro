package models

import (
	"github.com/beego/beego/v2/client/orm"
	"tinypro/common/types"
	"time"
)

func init() {
	orm.RegisterModel(new(Course))
}

const COURSE_TABLENAME = "course_tab"

// Course 结构体
type Course struct {
	ID                 int64             `orm:"pk;column(id)" json:"id"`                       // 主键ID
	CreatedAt          time.Time         `orm:"type(datetime);precision(3)" json:"created_at"` // 创建时间
	UpdatedAt          time.Time         // 更新时间
	Name               string            `json:"name" form:"name"`
	Logo               string            `json:"logo" form:"logo"`
	StartTime          time.Time         `json:"start_time" form:"start_time"`
	EndTime            time.Time         `json:"end_time" form:"end_time"`
	InviteLevel        string            `json:"invite_level" form:"invite_level"`
	CourseType         int               `json:"course_type" form:"course_type"`
	RemindFlag         int               `json:"remind_flag" form:"remind_flag"`
	PersonLimit        int               `json:"person_limit" form:"person_limit"`
	Teachers           string            `json:"teachers" form:"teachers"`
	EnrollHeaderPic    string            `json:"enroll_header_pic"`
	EnrollIntroPic     string            `json:"enroll_intro_pic"`
	ClassRoomHeaderPic string            `json:"class_room_header_pic"`
	ClassRoomIntroPic  string            `json:"class_room_intro_pic"`
	CourseBgPic        string            `json:"course_bg_pic"`
	WorkWxQR           string            `orm:"column(work_wx_qr)" json:"work_wx_qr" form:"work_wx_qr"`
	Materials          string            `json:"materials" form:"materials"`
	IsSelected         *int              `json:"is_selected" form:"is_selected"`
	SelectedWeight     int               `json:"selected_weight" form:"selected_weight"`
	MeetingType        types.MeetingEnum `json:"meeting_type" form:"meeting_type"`
	MeetingUrl         string            `json:"meeting_url" form:"meeting_url"`
	MeetingID          string            `orm:"column(meeting_id)" json:"meeting_id" form:"meeting_id"`
	MiniProgramUrl     string            `json:"mini_program_url" form:"mini_program_url"`
	Review             string            `json:"review" form:"review"`
}

func (r *Course) TableName() string {
	return COURSE_TABLENAME
}
