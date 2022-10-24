package course

import (
	"encoding/json"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/chester84/libtools"
	"github.com/gomodule/redigo/redis"
	"tinypro/common/lib/redis/cache"
	"tinypro/common/models"
	"tinypro/common/pkg/teacher"
	"tinypro/common/pogo/reqs"
	"tinypro/common/pogo/resps"
	"tinypro/common/types"
	"sort"
)

func FrontPage(req reqs.PageSelectedInfo) (retList []resps.FrontPage, err error) {
	m := models.Course{}
	o := orm.NewOrm()
	list := make([]models.Course, 0)
	retList = make([]resps.FrontPage, 0)

	pageSize := req.Size
	if req.Size > 100 {
		pageSize = types.DefaultPagesize
	}

	retList, err = getFrontPageListCache(req)
	if err != nil {
		logs.Error("FrontPage err %#v", err)
		return
	}

	// 缓存存在数据直接返回
	if len(retList) > 0 {
		return
	}

	if req.Type == 0 {
		if req.SelectedWeight <= 0 {
			//默认取最新记录
			_, err = o.QueryTable(m.TableName()).
				//Filter("id__in", ids).
				Filter("is_selected", 1).
				OrderBy("-selected_weight").
				Limit(pageSize).
				All(&list)
		} else {
			_, err = o.QueryTable(m.TableName()).
				//Filter("id__in", ids).
				Filter("is_selected", 1).
				Filter("selected_weight__gt", req.SelectedWeight).
				OrderBy("selected_weight").
				Limit(pageSize).
				All(&list)
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].SelectedWeight > list[j].SelectedWeight
		})

		//sql := fmt.Sprintf("SELECT * FROM course_tab WHERE is_selected = 1 AND created_at > ? ORDER BY ID ASC LIMIT ? ")
		//_, err = o.Raw(sql, req.Timestamp, pageSize).QueryRows(&list)
	} else {
		_, err = o.QueryTable(m.TableName()).
			//Filter("id__in", ids).
			Filter("is_selected", 1).
			Filter("selected_weight__lt", req.SelectedWeight).
			OrderBy("-selected_weight").
			All(&list)
	}

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("FrontPage db filter get exception, err: %v", err)
		return
	}

	//retList := make([]resps.FrontPage, 0)

	for _, item := range list {
		data := resps.FrontPage{}
		data.SN = item.ID
		data.Name = item.Name
		data.Logo = item.Logo
		data.EnrollHeaderPic = item.EnrollHeaderPic
		data.EnrollIntroPic = item.EnrollIntroPic
		data.ClassRoomHeaderPic = item.ClassRoomHeaderPic
		data.ClassRoomIntroPic = item.ClassRoomIntroPic
		data.CourseBgPic = item.CourseBgPic
		data.MeetingType = item.MeetingType
		data.MeetingDesc = types.OnlineMeetingMap()[item.MeetingType]
		data.InviteLevel = item.InviteLevel
		data.PersonLimit = item.PersonLimit
		data.SelectedWeight = item.SelectedWeight

		materials := make([]resps.CourseMaterial, 0)
		err = json.Unmarshal([]byte(item.Materials), &materials)
		if err != nil {
			logs.Error("course FrontPage json.Unmarshal Materials, err: %v", err)
			return
		}

		data.CourseMaterials = materials
		data.WorkWxQr = item.WorkWxQR
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

	err = setFrontPageListCache(req, retList)
	if err != nil {
		logs.Error("setFrontPageListCache err: %v", err)
	}

	return
}

func getFrontPageListCache(req reqs.PageSelectedInfo) (retList []resps.FrontPage, err error) {
	key := GetFrontPageListHashKey(req)
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	var retStr string
	retStr, err = redis.String(cacheClient.Do("HGET", FrontPageListHashDomain, key))
	if err != nil && err != redis.ErrNil {
		return
	} else {
		if retStr != "" {
			err = json.Unmarshal([]byte(retStr), &retList)
		} else {
			err = nil
		}
	}
	return
}

func setFrontPageListCache(req reqs.PageSelectedInfo, retList []resps.FrontPage) (err error) {
	key := GetFrontPageListHashKey(req)
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	var jsonB []byte
	jsonB, err = json.Marshal(retList)
	if err != nil {
		return
	}

	_, err = cacheClient.Do("HSET", FrontPageListHashDomain, key, string(jsonB))
	return
}

func PublicCourses(req reqs.PageInfo) (retList []resps.PublicCourses, err error) {
	m := models.Course{}
	o := orm.NewOrm()
	list := make([]models.Course, 0)
	retList = make([]resps.PublicCourses, 0)

	pageSize := req.Size
	if req.Size > 100 {
		pageSize = types.DefaultPagesize
	}

	retList, err = getPublicCourseListCache(req)
	if err != nil {
		logs.Error("publicCourseListCache err %#v", err)
		return
	}

	// 缓存存在数据直接返回
	if len(retList) > 0 {
		return
	}

	if req.Type == 0 {
		if req.SN <= 0 {
			//默认取最新记录
			_, err = o.QueryTable(m.TableName()).
				//Filter("id__in", ids).
				OrderBy("-id").
				Limit(pageSize).
				All(&list)
		} else {
			_, err = o.QueryTable(m.TableName()).
				//Filter("id__in", ids).
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
			Filter("id__lt", req.SN).
			OrderBy("-id").
			All(&list)
	}

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("[course->PublicCourses] db filter get exception, err: %v", err)
		return
	}

	for _, item := range list {
		var teacherIds []int64

		teacherInfo := make([]resps.TeacherInfo, 0)
		err = json.Unmarshal([]byte(item.Teachers), &teacherIds)
		teacherList, errI := teacher.GetTeachersByIdS(teacherIds)
		if errI != nil {
			logs.Error("[course->GetTeachersByIdS] db filter get exception, err: %v", err)
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

		data := resps.PublicCourses{}
		data.SN = item.ID
		data.Name = item.Name
		data.TeacherInfoList = teacherInfo
		data.Logo = item.Logo
		data.EnrollHeaderPic = item.EnrollHeaderPic
		data.EnrollIntroPic = item.EnrollIntroPic
		data.ClassRoomHeaderPic = item.ClassRoomHeaderPic
		data.ClassRoomIntroPic = item.ClassRoomIntroPic
		data.CourseBgPic = item.CourseBgPic
		data.InviteLevel = item.InviteLevel
		data.MeetingType = item.MeetingType
		data.MeetingDesc = types.OnlineMeetingMap()[item.MeetingType]
		data.PersonLimit = item.PersonLimit
		data.CourseMaterials = materials
		data.WorkWxQr = item.WorkWxQR
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

	err = setPublicCourseListCache(req, retList)
	if err != nil {
		logs.Error("setPublicCourseListCache err: %v", err)
	}
	return
}

func getPublicCourseListCache(req reqs.PageInfo) (retList []resps.PublicCourses, err error) {
	key := GetPublicCourseListHashKey(req)
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	var retStr string
	retStr, err = redis.String(cacheClient.Do("HGET", PublicCourseListHashDomain, key))
	if err != nil && err != redis.ErrNil {
		return
	} else {
		if retStr != "" {
			err = json.Unmarshal([]byte(retStr), &retList)
		} else {
			err = nil
		}
	}
	return
}

func setPublicCourseListCache(req reqs.PageInfo, retList []resps.PublicCourses) (err error) {
	key := GetPublicCourseListHashKey(req)
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	var jsonB []byte
	jsonB, err = json.Marshal(retList)
	if err != nil {
		return
	}

	_, err = cacheClient.Do("HSET", PublicCourseListHashDomain, key, string(jsonB))
	return
}

func GetRemindCourseList() (retList []models.Course, err error) {
	m := models.Course{}
	o := orm.NewOrm()
	list := make([]models.Course, 0)
	retList = make([]models.Course, 0)

	now := libtools.GetUnixMillis()

	// 服务端每10秒执行一次程序，这地方有个10秒的误差
	// 把这10秒算进去，这里多加了1秒...
	beforeStart := now + 60*60*1000 - 11*1000
	//beforeStart := now + 60*60*1000
	beforeStartTime := libtools.UnixMsec2Date(beforeStart, "Y-m-d H:i:s")

	_, err = o.QueryTable(m.TableName()).Filter("remind_flag", 0).Filter("start_time__gte", beforeStartTime).All(&list, "id", "start_time")

	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		return
	} else {
		for _, item := range list {
			logs.Debug("diff %d", item.StartTime.UnixMilli()-now)
			// 同样将10秒误差算进来
			diff := 1000*60*60 + 12*1000
			if item.StartTime.UnixMilli()-now <= int64(diff) {
				//如果相差在一个小时之内，
				//也就是提前一个小时提醒
				retList = append(retList, item)
			}
		}
		return retList, nil
	}
}
