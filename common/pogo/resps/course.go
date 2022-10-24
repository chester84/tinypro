package resps

import (
	"tinypro/common/types"
	"time"
)

type CourseMaterial struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Type string `json:"type"`
}

type FrontPage struct {
	SN                 int64             `json:"sn,string"`
	Name               string            `json:"name"`
	Logo               string            `json:"logo"`
	EnrollHeaderPic    string            `json:"enroll_header_pic"`
	EnrollIntroPic     string            `json:"enroll_intro_pic"`
	ClassRoomHeaderPic string            `json:"class_room_header_pic"`
	ClassRoomIntroPic  string            `json:"class_room_intro_pic"`
	CourseBgPic        string            `json:"course_bg_pic"`
	StartTime          string            `json:"start_time"`
	EndTime            string            `json:"end_time"`
	MeetingType        types.MeetingEnum `json:"meeting_type"`
	MeetingDesc        string            `json:"meeting_desc"`
	InviteLevel        string            `json:"invite_level"`
	PersonLimit        int               `json:"person_limit"`
	CourseMaterials    []CourseMaterial  `json:"course_materials"`
	WorkWxQr           string            `json:"work_wx_qr"`
	SelectedWeight     int               `json:"selected_weight"`
	//ClassRoomStatus     types.ClassRoomStatusEnum `json:"class_room_status"`
	//ClassRoomStatusDesc string                    `json:"class_room_status_desc"`
}

type TeacherInfo struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type PublicCourses struct {
	SN                 int64             `json:"sn,string"`
	Name               string            `json:"name"`
	TeacherInfoList    []TeacherInfo     `json:"teacher_info_list"`
	Logo               string            `json:"logo"`
	EnrollHeaderPic    string            `json:"enroll_header_pic"`
	EnrollIntroPic     string            `json:"enroll_intro_pic"`
	ClassRoomHeaderPic string            `json:"class_room_header_pic"`
	ClassRoomIntroPic  string            `json:"class_room_intro_pic"`
	CourseBgPic        string            `json:"course_bg_pic"`
	StartTime          string            `json:"start_time"`
	EndTime            string            `json:"end_time"`
	InviteLevel        string            `json:"invite_level"`
	MeetingType        types.MeetingEnum `json:"meeting_type"`
	MeetingDesc        string            `json:"meeting_desc"`
	PersonLimit        int               `json:"person_limit"`
	CourseMaterials    []CourseMaterial  `json:"course_materials"`
	WorkWxQr           string            `json:"work_wx_qr"`
	//ClassRoomStatus     types.ClassRoomStatusEnum `json:"class_room_status"`
	//ClassRoomStatusDesc string                    `json:"class_room_status_desc"`
}

type LastEnrollInfo struct {
	RealName        string `json:"real_name"`
	Mobile          string `json:"mobile"`
	Company         string `json:"company"`
	Position        string `json:"position"`
	Residence       string `json:"residence"`
	RecommendPerson string `json:"recommend_person"`
}

type UserEnroll struct {
	IsEnroll int `json:"is_enroll"`
	LastEnrollInfo
}

type LiveStreamItem struct {
	MeetingUrl     string `json:"meeting_url"`
	MeetingId      string `json:"meeting_id"`
	MiniProgramUrl string `json:"mini_program_url"`
	Review         string `json:"review"`
}

type MyCourses struct {
	SN                 int64             `json:"sn,string"`
	Name               string            `json:"name"`
	EnrollHeaderPic    string            `json:"enroll_header_pic"`
	EnrollIntroPic     string            `json:"enroll_intro_pic"`
	ClassRoomHeaderPic string            `json:"class_room_header_pic"`
	ClassRoomIntroPic  string            `json:"class_room_intro_pic"`
	TeacherInfoList    []TeacherInfo     `json:"teacher_info_list"`
	CourseBgPic        string            `json:"course_bg_pic"`
	StartTime          string            `json:"start_time"`
	EndTime            string            `json:"end_time"`
	InviteLevel        string            `json:"invite_level"`
	MeetingType        types.MeetingEnum `json:"meeting_type"`
	MeetingDesc        string            `json:"meeting_desc"`
	IsPass             int               `json:"is_pass"`
	IsSignIn           int               `json:"is_sign_in"`
	PersonLimit        int               `json:"person_limit"`
	CourseMaterials    []CourseMaterial  `json:"course_materials"`
	WorkWxQr           string            `json:"work_wx_qr"`
	LiveStream         LiveStreamItem    `json:"live_stream"`
	//ClassRoomStatus     types.ClassRoomStatusEnum `json:"class_room_status"`
	//ClassRoomStatusDesc string                    `json:"class_room_status_desc"`
}

type MyHistoryCoursesDB struct {
	Id                 int64             `json:"id"`
	Name               string            `json:"name"`
	EnrollHeaderPic    string            `json:"enroll_header_pic"`
	EnrollIntroPic     string            `json:"enroll_intro_pic"`
	ClassRoomHeaderPic string            `json:"class_room_header_pic"`
	ClassRoomIntroPic  string            `json:"class_room_intro_pic"`
	Teachers           string            `json:"teachers"`
	CourseId           string            `json:"course_id"`
	IsPass             int               `json:"is_pass"`
	IsSignIn           int               `json:"is_sign_in"`
	CourseBgPic        string            `json:"course_bg_pic"`
	StartTime          time.Time         `json:"start_time"`
	EndTime            time.Time         `json:"end_time"`
	InviteLevel        string            `json:"invite_level"`
	MeetingType        types.MeetingEnum `json:"meeting_type"`
	MeetingDesc        string            `json:"meeting_desc"`
	PersonLimit        int               `json:"person_limit"`
	Materials          string            `json:"materials"`
	WorkWxQr           string            `json:"work_wx_qr"`
	MeetingUrl         string            `json:"meeting_url"`
	MeetingId          string            `json:"meeting_id"`
	MiniProgramUrl     string            `json:"mini_program_url"`
	Review             string            `json:"review"`
}

type MyHistoryCourses struct {
	SN                 string            `json:"sn"`
	Name               string            `json:"name"`
	EnrollHeaderPic    string            `json:"enroll_header_pic"`
	EnrollIntroPic     string            `json:"enroll_intro_pic"`
	ClassRoomHeaderPic string            `json:"class_room_header_pic"`
	ClassRoomIntroPic  string            `json:"class_room_intro_pic"`
	TeacherInfoList    []TeacherInfo     `json:"teacher_info_list"`
	CourseId           string            `json:"course_id"`
	IsPass             int               `json:"is_pass"`
	IsSignIn           int               `json:"is_sign_in"`
	CourseBgPic        string            `json:"course_bg_pic"`
	StartTime          string            `json:"start_time"`
	EndTime            string            `json:"end_time"`
	InviteLevel        string            `json:"invite_level"`
	MeetingType        types.MeetingEnum `json:"meeting_type"`
	MeetingDesc        string            `json:"meeting_desc"`
	PersonLimit        int               `json:"person_limit"`
	CourseMaterials    []CourseMaterial  `json:"course_materials"`
	WorkWxQr           string            `json:"work_wx_qr"`
	LiveStream         LiveStreamItem    `json:"live_stream"`
	//ClassRoomStatus     types.ClassRoomStatusEnum `json:"class_room_status"`
	//ClassRoomStatusDesc string                    `json:"class_room_status_desc"`
}
