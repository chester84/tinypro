package types

type MsgEnum int

const (
	NewCourse  MsgEnum = 1
	ClassStart MsgEnum = 2
)

func MsgTypeMap() map[MsgEnum]string {
	return map[MsgEnum]string{
		NewCourse:  "上新提醒",
		ClassStart: "课程提醒",
	}
}
