package enrollbiz

import (
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/chester84/libtools"
	"tinypro/common/models"
	"tinypro/common/pogo/reqs"
	"tinypro/common/pogo/resps"
	"time"
)

func GetLastEnroll(userObj models.AppUser, courseID int64) (ret resps.UserEnroll, err error) {
	o := orm.NewOrm()
	obj := models.Enroll{}

	err = o.QueryTable(obj.TableName()).
		Filter("user_id", userObj.Id).
		Filter("course_id", courseID).
		OrderBy("-id").
		One(&obj)

	if err != nil {
		if err.Error() != orm.ErrNoRows.Error() {
			return
		} else {
			//如果记录为空，那么没有报名
			err = nil
			ret.IsEnroll = 0
		}
	} else {
		ret.IsEnroll = 1
	}

	err = o.QueryTable(obj.TableName()).
		Filter("user_id", userObj.Id).
		OrderBy("-id").
		One(&obj)

	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		return
	} else {
		err = nil
	}

	mobile := ""
	if obj.Mobile == "" {
		mobile = userObj.Mobile
	} else {
		mobile = obj.Mobile
	}

	ret.Mobile = mobile
	ret.RealName = obj.RealName
	ret.Company = obj.Company
	ret.Position = obj.Position
	ret.Residence = obj.Residence
	ret.RecommendPerson = obj.RecommendPerson

	return
}

func EnrollCourse(req reqs.EnrollReqT, userId int64) (err error) {
	now := time.Now()
	o := orm.NewOrm()
	obj := models.Enroll{}
	courseId, _ := libtools.Str2Int64(req.CourseSN)

	err = o.QueryTable(obj.TableName()).Filter("user_id", userId).Filter("course_id", courseId).Limit(1).One(&obj)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		return
	} else {
		if err == nil {
			obj.UpdatedAt = now
			_, err = models.OrmUpdate(&obj, []string{"UpdatedAt"})
		} else {
			enrollObj := models.Enroll{
				CreatedAt:       now,
				UpdatedAt:       now,
				UserID:          userId,
				CourseId:        courseId,
				RealName:        req.RealName,
				Mobile:          req.Mobile,
				Company:         req.Company,
				Position:        req.Position,
				Residence:       req.Residence,
				RecommendPerson: req.RecommendPerson,
			}
			_, err = models.OrmInsert(&enrollObj)
			if err != nil {
				return
			}
		}
	}

	return
}

func IsPassEnroll(userId, courseID int64) (ret int, err error) {
	o := orm.NewOrm()
	obj := models.Enroll{}

	err = o.QueryTable(obj.TableName()).
		Filter("user_id", userId).
		Filter("course_id", courseID).
		Filter("is_pass", 1).
		Limit(1).
		One(&obj)

	if err != nil {
		if err.Error() != orm.ErrNoRows.Error() {
			return
		} else {
			//如果记录为空，那么没有报名
			err = nil
			ret = 0
		}
	} else {
		ret = 1
	}

	return
}
