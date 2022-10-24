package course

import (
	"fmt"
	"tinypro/common/pogo/reqs"
)

const (
	FrontPageListHashDomain    = "tinypro:hash:front-page"
	PublicCourseListHashDomain = "tinypro:hash:public-course"
)

func GetPublicCourseListHashKey(req reqs.PageInfo) (key string) {
	key = fmt.Sprintf("%d:%d:%d", req.Type, req.SN, req.Size)
	return
}

func GetFrontPageListHashKey(req reqs.PageSelectedInfo) (key string) {
	key = fmt.Sprintf("%d:%d:%d", req.Type, req.SelectedWeight, req.Size)
	return
}
