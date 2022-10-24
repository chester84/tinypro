package subscribe_template

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/chester84/libtools"
	"tinypro/common/models"
	"tinypro/common/pkg/teacher"
	"tinypro/common/pkg/weixin"
	"tinypro/common/types"
)

func BuildRemindAttendClassMsg(courseObj models.Course) (data map[string]weixin.SubscribeMessageSendDataItem, err error) {
	data = make(map[string]weixin.SubscribeMessageSendDataItem)

	data["thing2"] = weixin.SubscribeMessageSendDataItem{
		Value: courseObj.Name,
	}

	address := ""
	if courseObj.MeetingType == types.OnlineMeeting {
		address = types.OnlineMeetingMap()[courseObj.MeetingType]
	}
	data["thing10"] = weixin.SubscribeMessageSendDataItem{
		Value: address,
	}

	var teacherIds []int64
	teacherInfo := make([]string, 0)
	err = json.Unmarshal([]byte(courseObj.Teachers), &teacherIds)
	teacherList, errI := teacher.GetTeachersByIdS(teacherIds)
	if errI != nil {
		logs.Error("BuildRemindAttendClassMsg db filter get exception, err: %v", err)
		err = errI
		return
	}
	for _, teacherObj := range teacherList {
		teacherInfo = append(teacherInfo, teacherObj.Name)
	}
	teacherStr := libtools.ArrayToString(teacherInfo, " ")

	data["thing15"] = weixin.SubscribeMessageSendDataItem{
		Value: teacherStr,
	}

	minute := "00"
	if courseObj.StartTime.Minute() < 10 {
		minute = fmt.Sprintf("0%d", courseObj.StartTime.Minute())
	} else {
		minute = fmt.Sprintf("%d", courseObj.StartTime.Minute())
	}

	datetime := fmt.Sprintf("%d年%d月%d日 %d:%s", courseObj.StartTime.Year(), courseObj.StartTime.Month(), courseObj.StartTime.Day(), courseObj.StartTime.Hour(), minute)
	data["time32"] = weixin.SubscribeMessageSendDataItem{
		Value: datetime,
	}

	return
}
