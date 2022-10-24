package satissurvey

import (
	"github.com/beego/beego/v2/client/orm"
	"tinypro/common/models"
	"tinypro/common/pogo/reqs"
	"github.com/chester84/libtools"
	"tinypro/common/types"
	"time"
)

func GetOneByUserIdCourseId(userId, courseId int64) (m models.SatisSurvey, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(m.TableName()).Filter("user_id", userId).Filter("course_id", courseId).Limit(1).One(&m)

	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		return
	} else {
		err = nil
	}

	return
}

func UserGiveScore(userObj models.AppUser, req reqs.QScore) (err error) {
	m := models.SatisSurvey{}
	courseId, _ := libtools.Str2Int64(req.CourseSn)
	m, err = GetOneByUserIdCourseId(userObj.Id, courseId)
	if err != nil {
		return
	}

	if m.ID <= 0 {
		now := time.Now()
		m.UserId = userObj.Id
		m.CourseId, _ = libtools.Str2Int64(req.CourseSn)
		m.Q1 = req.Q1
		m.Q2 = req.Q2
		m.Q3 = req.Q3
		m.Q4 = req.Q4
		m.Q5 = req.Q5
		m.Q6 = req.Q6
		m.Q7 = req.Q7
		m.Q8 = req.Q8
		m.Q9 = req.Q9
		m.Q10 = req.Q10
		m.CreatedAt = now
		m.UpdatedAt = now

		_, err = models.OrmInsert(&m)
	}

	return
}

func SatisQuestion() []types.SatisSurveyQ {
	return []types.SatisSurveyQ{
		types.SatisSurveyQ{Q: "q1", A: "课程知识体系清晰、内容扎实"},
		types.SatisSurveyQ{Q: "q2", A: "课程内容及案例实用性、启发性强"},
		types.SatisSurveyQ{Q: "q3", A: "课程很好地将学术理论和管理实践相结合"},
		types.SatisSurveyQ{Q: "q4", A: "老师严格要求同学们积极参与学习，并成功地引导他们进行互动讨论"},
		types.SatisSurveyQ{Q: "q5", A: "老师认真考虑了同学们就课程内容提出的问题和观点"},
		types.SatisSurveyQ{Q: "q6", A: "老师的敬业精神"},
		types.SatisSurveyQ{Q: "q7", A: "我学到了与实际工作相关的技能、工具和思考方法"},
		types.SatisSurveyQ{Q: "q8", A: "开拓视野、体验多元"},
		types.SatisSurveyQ{Q: "q9", A: "增进思想交流、结交新朋友"},
		types.SatisSurveyQ{Q: "q10", A: "整体课程体验"},
	}
}

func SurveyScoreConfig() []types.SatisSurveyScore {
	return []types.SatisSurveyScore{
		types.SatisSurveyScore{S: 1, D: "非常不满意"},
		types.SatisSurveyScore{S: 2, D: "比较不满意"},
		types.SatisSurveyScore{S: 3, D: "一般"},
		types.SatisSurveyScore{S: 4, D: "比较满意"},
		types.SatisSurveyScore{S: 5, D: "非常满意"},
	}
}
