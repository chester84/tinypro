package advertisingpopup

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"tinypro/common/lib/redis/cache"
	"tinypro/common/models"
	"tinypro/common/pogo/resps"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func GetOneById(id int64) (obj models.AdvertisingPopup, err error) {
	err = models.OrmOneByPkId(id, &obj)
	return
}

func GetUserAccessAdvertisingPopup(userObj models.AppUser, advertiseId int64) (resp resps.AccessPopupResp, err error) {
	obj := models.AdvertisingPopup{}
	err = models.OrmOneByPkId(advertiseId, &obj)
	if err != nil {
		err = fmt.Errorf("GetUserAccessAdvertisingPopup get empty advertising, err %#v", err)
		logs.Error(err)
		return
	}
	// 未配置广告弹窗
	if obj.ID <= 0 {
		return
	}

	courseObj := models.Course{}
	err = models.OrmOneByPkId(obj.CourseId, &courseObj)
	if err != nil {
		err = fmt.Errorf("GetUserAccessAdvertisingPopup get empty Course, err %#v", err)
		logs.Error(err)
		return
	}

	//配置的课程不存在
	if courseObj.ID <= 0 {
		return
	}

	// 开关关闭
	if obj.Switch == 0 {
		// 1表示访问过了
		resp.IsAccess = 1
		return
	}

	rdsKey := GetAccessPopupRdsKey(advertiseId, userObj.Id)
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	cValue, err := cacheClient.Do("GET", rdsKey)
	if err != nil {
		err = fmt.Errorf("GetUserAccessAdvertisingPopup cacheClient GET err, %#v", err)
		return
	}

	if cValue == nil {
		resp.IsAccess = 0
		_, err = cacheClient.Do("SETEX", rdsKey, types.DaySecond, 1)
		if err != nil {
			err = fmt.Errorf("GetUserAccessAdvertisingPopup cacheClient GET err, %#v", err)
			return
		}
		resp.AdvertiseUrl = obj.Url
		resp.CourseNode.SN = courseObj.ID
		resp.CourseNode.Name = courseObj.Name
		resp.CourseNode.Logo = courseObj.Logo
		resp.CourseNode.EnrollHeaderPic = courseObj.EnrollHeaderPic
		resp.CourseNode.EnrollIntroPic = courseObj.EnrollIntroPic
		resp.CourseNode.ClassRoomHeaderPic = courseObj.ClassRoomHeaderPic
		resp.CourseNode.ClassRoomIntroPic = courseObj.ClassRoomIntroPic
		resp.CourseNode.CourseBgPic = courseObj.CourseBgPic
		resp.CourseNode.MeetingType = courseObj.MeetingType
		resp.CourseNode.MeetingDesc = types.OnlineMeetingMap()[courseObj.MeetingType]
		resp.CourseNode.InviteLevel = courseObj.InviteLevel
		resp.CourseNode.PersonLimit = courseObj.PersonLimit
		resp.CourseNode.SelectedWeight = courseObj.SelectedWeight

		materials := make([]resps.CourseMaterial, 0)
		err = json.Unmarshal([]byte(courseObj.Materials), &materials)
		if err != nil {
			logs.Error("GetUserAccessAdvertisingPopup FrontPage json.Unmarshal Materials, err: %v", err)
			return
		}

		resp.CourseNode.CourseMaterials = materials
		resp.CourseNode.WorkWxQr = courseObj.WorkWxQR
		resp.CourseNode.StartTime = libtools.UnixMsec2Date(courseObj.StartTime.UnixMilli(), "Y-m-d H:i")
		resp.CourseNode.EndTime = libtools.UnixMsec2Date(courseObj.EndTime.UnixMilli(), "Y-m-d H:i")
	} else {
		resp.IsAccess = 1
	}

	return
}
