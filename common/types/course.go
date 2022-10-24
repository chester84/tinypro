package types

type ClassRoomStatusEnum int

const (
	ClassRoomEnroll ClassRoomStatusEnum = 1
	ClassRoomLive   ClassRoomStatusEnum = 2
	ClassRoomEnd    ClassRoomStatusEnum = 3
)

func ClassRoomStatusEnumMap() map[ClassRoomStatusEnum]string {
	return map[ClassRoomStatusEnum]string{
		ClassRoomEnroll: "开放报名",
		ClassRoomLive:   "直播中",
		ClassRoomEnd:    "已结束",
	}
}

type MeetingEnum int

const (
	OnlineMeeting MeetingEnum = 1
)

func OnlineMeetingMap() map[MeetingEnum]string {
	return map[MeetingEnum]string{
		OnlineMeeting: "腾讯会议",
	}
}

func SatisScoreMap() map[int]string {
	return map[int]string{
		1: "非常不满意",
		2: "比较不满意",
		3: "一般",
		4: "比较满意",
		5: "非常满意",
	}
}
