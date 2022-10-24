package enroll

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/adapter/orm"
	"tinypro/common/models"
	"tinypro/common/pkg/teacher"
	"tinypro/common/pogo/reqs"
	"tinypro/common/pogo/resps"
	"github.com/chester84/libtools"
	"tinypro/common/types"
	"sort"
	"time"
)

func MyEnrolls(userObj models.AppUser, req reqs.PageInfo) (retList []resps.MyCourses, err error) {
	m := models.Enroll{}
	o := orm.NewOrm()
	list := make([]models.Enroll, 0)
	retList = make([]resps.MyCourses, 0)
	//now := time.Now().UnixMilli()

	pageSize := req.Size
	if req.Size > 100 {
		pageSize = types.DefaultPagesize
	}

	if req.Type == 0 {
		if req.SN <= 0 {
			//默认取最新记录
			_, err = o.QueryTable(m.TableName()).
				Filter("user_id", userObj.Id).
				OrderBy("-id").
				Limit(pageSize).
				All(&list)
		} else {
			_, err = o.QueryTable(m.TableName()).
				Filter("user_id", userObj.Id).
				Filter("id__gt", req.SN).
				OrderBy("id").
				Limit(pageSize).
				All(&list)
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].ID > list[j].ID
		})
	} else {
		_, err = o.QueryTable(m.TableName()).
			Filter("user_id", userObj.Id).
			Filter("id__lt", req.SN).
			OrderBy("-id").
			Limit(pageSize).
			All(&list)
	}

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("MyEnrolls PublicCourses db filter get exception, err: %v", err)
		return
	}

	for _, item := range list {
		var courseObj models.Course
		err = models.OrmOneByPkId(item.CourseId, &courseObj)
		if err != nil {
			logs.Error("MyEnrolls Course OrmOneByPkId db filter get exception, err: %v", err)
			return
		}

		var teacherIds []int64

		teacherInfo := make([]resps.TeacherInfo, 0)
		err = json.Unmarshal([]byte(courseObj.Teachers), &teacherIds)
		teacherList, errI := teacher.GetTeachersByIdS(teacherIds)
		if errI != nil {
			logs.Error("MyEnrolls GetTeachersByIdS db filter get exception, err: %v", err)
			err = errI
			return
		}
		for _, teacherObj := range teacherList {
			tmp := resps.TeacherInfo{}
			tmp.Url = teacherObj.Header
			tmp.Name = teacherObj.Name
			teacherInfo = append(teacherInfo, tmp)
		}

		materials := make([]resps.CourseMaterial, 0)
		err = json.Unmarshal([]byte(courseObj.Materials), &materials)
		if err != nil {
			logs.Error("course FrontPage json.Unmarshal Materials, err: %v", err)
			return
		}

		data := resps.MyCourses{}
		data.SN = item.ID
		data.Name = courseObj.Name
		data.IsPass = item.IsPass
		data.IsSignIn = item.IsSignIn
		data.TeacherInfoList = teacherInfo
		data.EnrollHeaderPic = courseObj.EnrollHeaderPic
		data.EnrollIntroPic = courseObj.EnrollIntroPic
		data.ClassRoomHeaderPic = courseObj.ClassRoomHeaderPic
		data.ClassRoomIntroPic = courseObj.ClassRoomIntroPic
		data.CourseBgPic = courseObj.CourseBgPic
		data.InviteLevel = courseObj.InviteLevel
		data.MeetingType = courseObj.MeetingType
		data.MeetingDesc = types.OnlineMeetingMap()[courseObj.MeetingType]
		data.PersonLimit = courseObj.PersonLimit
		data.CourseMaterials = materials
		data.WorkWxQr = courseObj.WorkWxQR
		data.LiveStream.MiniProgramUrl = courseObj.MiniProgramUrl
		data.LiveStream.Review = courseObj.Review
		data.LiveStream.MeetingUrl = courseObj.MeetingUrl
		data.LiveStream.MeetingId = courseObj.MeetingID

		data.StartTime = libtools.UnixMsec2Date(courseObj.StartTime.UnixMilli(), "Y-m-d H:i")
		data.EndTime = libtools.UnixMsec2Date(courseObj.EndTime.UnixMilli(), "Y-m-d H:i")

		//if courseObj.StartTime.UnixMilli() > now {
		//	data.ClassRoomStatus = types.ClassRoomEnroll
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomEnroll]
		//} else if courseObj.StartTime.UnixMilli() < now && courseObj.EndTime.UnixMilli() > now {
		//	data.ClassRoomStatus = types.ClassRoomLive
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomLive]
		//} else {
		//	data.ClassRoomStatus = types.ClassRoomEnd
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomEnd]
		//}

		retList = append(retList, data)
	}

	return
}

func MyCourses(userObj models.AppUser, req reqs.PageInfo) (retList []resps.MyCourses, err error) {
	m := models.Enroll{}
	o := orm.NewOrm()
	list := make([]models.Enroll, 0)
	retList = make([]resps.MyCourses, 0)
	//now := time.Now().UnixMilli()

	pageSize := req.Size
	if req.Size > 100 {
		pageSize = types.DefaultPagesize
	}

	if req.Type == 0 {
		if req.SN <= 0 {
			//默认取最新记录
			_, err = o.QueryTable(m.TableName()).
				Filter("user_id", userObj.Id).
				Filter("is_pass", 1).
				OrderBy("-id").
				Limit(pageSize).
				All(&list)
		} else {
			_, err = o.QueryTable(m.TableName()).
				Filter("user_id", userObj.Id).
				Filter("is_pass", 1).
				Filter("id__gt", req.SN).
				OrderBy("id").
				Limit(pageSize).
				All(&list)
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].ID > list[j].ID
		})
	} else {
		_, err = o.QueryTable(m.TableName()).
			Filter("user_id", userObj.Id).
			Filter("is_pass", 1).
			Filter("id__lt", req.SN).
			OrderBy("-id").
			Limit(pageSize).
			All(&list)
	}

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("MyCourses PublicCourses db filter get exception, err: %v", err)
		return
	}

	for _, item := range list {
		var courseObj models.Course
		err = models.OrmOneByPkId(item.CourseId, &courseObj)
		if err != nil {
			logs.Error("MyCourses Course OrmOneByPkId db filter get exception, err: %v", err)
			return
		}

		var teacherIds []int64

		teacherInfo := make([]resps.TeacherInfo, 0)
		err = json.Unmarshal([]byte(courseObj.Teachers), &teacherIds)
		teacherList, errI := teacher.GetTeachersByIdS(teacherIds)
		if errI != nil {
			logs.Error("MyCourses GetTeachersByIdS db filter get exception, err: %v", err)
			err = errI
			return
		}
		for _, teacherObj := range teacherList {
			tmp := resps.TeacherInfo{}
			tmp.Url = teacherObj.Header
			tmp.Name = teacherObj.Name
			teacherInfo = append(teacherInfo, tmp)
		}

		materials := make([]resps.CourseMaterial, 0)
		err = json.Unmarshal([]byte(courseObj.Materials), &materials)
		if err != nil {
			logs.Error("course FrontPage json.Unmarshal Materials, err: %v", err)
			return
		}

		data := resps.MyCourses{}
		data.SN = item.ID
		data.Name = courseObj.Name
		data.IsPass = item.IsPass
		data.IsSignIn = item.IsSignIn
		data.TeacherInfoList = teacherInfo
		data.EnrollHeaderPic = courseObj.EnrollHeaderPic
		data.EnrollIntroPic = courseObj.EnrollIntroPic
		data.ClassRoomHeaderPic = courseObj.ClassRoomHeaderPic
		data.ClassRoomIntroPic = courseObj.ClassRoomIntroPic
		data.CourseBgPic = courseObj.CourseBgPic
		data.InviteLevel = courseObj.InviteLevel
		data.MeetingType = courseObj.MeetingType
		data.MeetingDesc = types.OnlineMeetingMap()[courseObj.MeetingType]
		data.PersonLimit = courseObj.PersonLimit
		data.CourseMaterials = materials
		data.WorkWxQr = courseObj.WorkWxQR
		data.LiveStream.MiniProgramUrl = courseObj.MiniProgramUrl
		data.LiveStream.Review = courseObj.Review
		data.LiveStream.MeetingUrl = courseObj.MeetingUrl
		data.LiveStream.MeetingId = courseObj.MeetingID
		data.StartTime = libtools.UnixMsec2Date(courseObj.StartTime.UnixMilli(), "Y-m-d H:i")
		data.EndTime = libtools.UnixMsec2Date(courseObj.EndTime.UnixMilli(), "Y-m-d H:i")

		//if courseObj.StartTime.UnixMilli() > now {
		//	data.ClassRoomStatus = types.ClassRoomEnroll
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomEnroll]
		//} else if courseObj.StartTime.UnixMilli() < now && courseObj.EndTime.UnixMilli() > now {
		//	data.ClassRoomStatus = types.ClassRoomLive
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomLive]
		//} else {
		//	data.ClassRoomStatus = types.ClassRoomEnd
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomEnd]
		//}

		retList = append(retList, data)
	}

	return
}

func MyHistoryCourses(userObj models.AppUser, req reqs.PageInfo) (retList []resps.MyHistoryCourses, err error) {
	o := orm.NewOrm()
	retList = make([]resps.MyHistoryCourses, 0)
	list := make([]resps.MyHistoryCoursesDB, 0)
	now := time.Now()

	pageSize := req.Size
	if req.Size > 100 {
		pageSize = types.DefaultPagesize
	}

	if req.Type == 0 {
		if req.SN <= 0 {
			//默认取最新记录
			sql := fmt.Sprintf(`select e.id, e.course_id, e.is_pass, e.is_sign_in, c.enroll_header_pic, c.enroll_intro_pic, c.class_room_header_pic, c.class_room_intro_pic, c.course_bg_pic,  c.work_wx_qr, c.materials, c.meeting_url, c.meeting_id, c.mini_program_url, c.review, c.teachers, c.name, c.invite_level,c.meeting_type,c.person_limit,c.start_time,c.end_time from enroll_tab e left join course_tab c on e.course_id = c.id where e.user_id = ? and e.is_pass = 1 and c.end_time < ? order by e.id desc limit ?`)
			_, err = o.Raw(sql, userObj.Id, now, pageSize).QueryRows(&list)
		} else {
			sql := fmt.Sprintf(`select e.id, e.course_id, e.is_pass, e.is_sign_in, c.enroll_header_pic, c.enroll_intro_pic, c.class_room_header_pic, c.class_room_intro_pic, c.course_bg_pic, c.work_wx_qr, c.materials, c.meeting_url, c.meeting_id, c.mini_program_url, c.review, c.teachers, c.name, c.invite_level,c.meeting_type,c.person_limit,c.start_time,c.end_time from enroll_tab e left join course_tab c on e.course_id = c.id where e.id > ? e.user_id = ? and e.is_pass = 1 and c.end_time < ? order by e.id asc limit ?`)
			_, err = o.Raw(sql, req.SN, userObj.Id, now, pageSize).QueryRows(&list)
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].Id > list[j].Id
		})
	} else {
		sql := fmt.Sprintf(`select e.id, e.course_id, e.is_pass, e.is_sign_in, c.enroll_header_pic, c.enroll_intro_pic, c.class_room_header_pic, c.class_room_intro_pic, c.course_bg_pic, c.work_wx_qr, c.materials, c.meeting_url, c.meeting_id, c.mini_program_url, c.review, c.teachers, c.name, c.invite_level,c.meeting_type,c.person_limit,c.start_time,c.end_time from enroll_tab e left join course_tab c on e.course_id = c.id where e.id < ? e.user_id = ? and e.is_pass = 1 and c.end_time < ? order by e.id desc limit ?`)
		_, err = o.Raw(sql, req.SN, userObj.Id, now, pageSize).QueryRows(&list)
	}

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("MyCourses PublicCourses db filter get exception, err: %v", err)
		return
	}

	logs.Debug("list %#v", list)

	for _, item := range list {
		var teacherIds []int64

		teacherInfo := make([]resps.TeacherInfo, 0)
		err = json.Unmarshal([]byte(item.Teachers), &teacherIds)
		teacherList, errI := teacher.GetTeachersByIdS(teacherIds)
		if errI != nil {
			logs.Error("MyCourses GetTeachersByIdS db filter get exception, err: %v", err)
			err = errI
			return
		}
		for _, teacherObj := range teacherList {
			tmp := resps.TeacherInfo{}
			tmp.Url = teacherObj.Header
			tmp.Name = teacherObj.Name
			teacherInfo = append(teacherInfo, tmp)
		}

		materials := make([]resps.CourseMaterial, 0)
		err = json.Unmarshal([]byte(item.Materials), &materials)
		if err != nil {
			logs.Error("course FrontPage json.Unmarshal Materials, err: %v", err)
			return
		}

		data := resps.MyHistoryCourses{}

		data.SN = libtools.Int642Str(item.Id)
		data.Name = item.Name
		data.MeetingType = item.MeetingType
		data.EnrollHeaderPic = item.EnrollHeaderPic
		data.EnrollIntroPic = item.EnrollIntroPic
		data.ClassRoomHeaderPic = item.ClassRoomHeaderPic
		data.ClassRoomIntroPic = item.ClassRoomIntroPic
		data.CourseBgPic = item.CourseBgPic
		data.IsPass = item.IsPass
		data.IsSignIn = item.IsSignIn
		data.InviteLevel = item.InviteLevel
		data.PersonLimit = item.PersonLimit
		data.CourseMaterials = materials
		data.WorkWxQr = item.WorkWxQr
		data.CourseId = item.CourseId
		data.TeacherInfoList = teacherInfo
		data.MeetingDesc = types.OnlineMeetingMap()[item.MeetingType]
		data.LiveStream.MeetingId = item.MeetingId
		data.LiveStream.MeetingUrl = item.MeetingUrl
		data.LiveStream.MiniProgramUrl = item.MiniProgramUrl
		data.LiveStream.Review = item.Review
		data.StartTime = libtools.UnixMsec2Date(item.StartTime.UnixMilli(), "Y-m-d H:i")
		data.EndTime = libtools.UnixMsec2Date(item.EndTime.UnixMilli(), "Y-m-d H:i")

		//if item.StartTime.UnixMilli() > now {
		//	data.ClassRoomStatus = types.ClassRoomEnroll
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomEnroll]
		//} else if item.StartTime.UnixMilli() < now && item.EndTime.UnixMilli() > now {
		//	data.ClassRoomStatus = types.ClassRoomLive
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomLive]
		//} else {
		//	data.ClassRoomStatus = types.ClassRoomEnd
		//	data.ClassRoomStatusDesc = types.ClassRoomStatusEnumMap()[types.ClassRoomEnd]
		//}

		retList = append(retList, data)
	}

	return
}

func GetOneByUserIdCourseId(userId, courseId int64) (m models.Enroll, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(m.TableName()).
		Filter("user_id", userId).
		Filter("course_id", courseId).
		Limit(1).
		One(&m)

	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		return
	} else {
		return m, nil
	}
}

func GetEnrollDetail(userId, courseId int64) (resp resps.MyCourses, err error) {
	var m models.Enroll
	m, err = GetOneByUserIdCourseId(userId, courseId)
	if err != nil {
		return
	}

	var courseObj models.Course
	err = models.OrmOneByPkId(courseId, &courseObj)
	if err != nil {
		return
	}

	var teacherIds []int64

	teacherInfo := make([]resps.TeacherInfo, 0)
	err = json.Unmarshal([]byte(courseObj.Teachers), &teacherIds)
	teacherList, errI := teacher.GetTeachersByIdS(teacherIds)
	if errI != nil {
		logs.Error("GetEnrollDetail GetTeachersByIdS db filter get exception, err: %v", err)
		err = errI
		return
	}
	for _, teacherObj := range teacherList {
		tmp := resps.TeacherInfo{}
		tmp.Url = teacherObj.Header
		tmp.Name = teacherObj.Name
		teacherInfo = append(teacherInfo, tmp)
	}

	materials := make([]resps.CourseMaterial, 0)
	err = json.Unmarshal([]byte(courseObj.Materials), &materials)
	if err != nil {
		logs.Error("GetEnrollDetail json.Unmarshal Materials, err: %v", err)
		return
	}

	resp.SN = m.ID
	resp.Name = courseObj.Name
	resp.IsPass = m.IsPass
	resp.IsSignIn = m.IsSignIn
	resp.TeacherInfoList = teacherInfo
	resp.EnrollHeaderPic = courseObj.EnrollHeaderPic
	resp.EnrollIntroPic = courseObj.EnrollIntroPic
	resp.ClassRoomHeaderPic = courseObj.ClassRoomHeaderPic
	resp.ClassRoomIntroPic = courseObj.ClassRoomIntroPic
	resp.CourseBgPic = courseObj.CourseBgPic
	resp.InviteLevel = courseObj.InviteLevel
	resp.MeetingType = courseObj.MeetingType
	resp.MeetingDesc = types.OnlineMeetingMap()[courseObj.MeetingType]
	resp.PersonLimit = courseObj.PersonLimit
	resp.CourseMaterials = materials
	resp.WorkWxQr = courseObj.WorkWxQR
	resp.LiveStream.MiniProgramUrl = courseObj.MiniProgramUrl
	resp.LiveStream.Review = courseObj.Review
	resp.LiveStream.MeetingUrl = courseObj.MeetingUrl
	resp.LiveStream.MeetingId = courseObj.MeetingID
	resp.StartTime = libtools.UnixMsec2Date(courseObj.StartTime.UnixMilli(), "Y-m-d H:i")
	resp.EndTime = libtools.UnixMsec2Date(courseObj.EndTime.UnixMilli(), "Y-m-d H:i")

	return
}

func SignInCourse(userObj models.AppUser, req reqs.SignInReqT) (m models.Enroll, err error) {
	courseId, _ := libtools.Str2Int64(req.CourseSN)
	m, err = GetOneByUserIdCourseId(userObj.Id, courseId)
	if err != nil {
		return
	} else {
		m.IsSignIn = 1
		_, err = models.OrmUpdate(&m, []string{"IsSignIn"})
	}

	return
}
