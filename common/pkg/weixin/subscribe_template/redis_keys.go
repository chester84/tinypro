package subscribe_template

import "fmt"

const (
	tmplIDCourseRemind      = "UKp_zq3qN4H0tylWHQGcDFbah8UOyAbRzeC2EnCzReY"
	tmplIDCourseEnrollCheck = "h3IybQ9m6kqWtQOZN9WackqKQGfSLEDoQaFEg6iyokQ"
)

const (
	rdsKeySubscribePrefix = `tinypro:set:subscribe`
)

func TmplMap() map[string]string {
	mapTmpl := make(map[string]string)
	mapTmpl[tmplIDCourseRemind] = tmplIDCourseRemind
	mapTmpl[tmplIDCourseEnrollCheck] = tmplIDCourseEnrollCheck
	return mapTmpl
}

func CheckTmplMapKey(tmplId string) string {
	if v, ok := TmplMap()[tmplId]; ok {
		return v
	} else {
		return ""
	}
}

func GetRDSCourseTmplIDKey(tmplId string, courseSn string) string {
	return fmt.Sprintf("%s:%s:%s", rdsKeySubscribePrefix, tmplId, courseSn)
}

func GetRDSRemindCourseKey(courseSn string) string {
	return GetRDSCourseTmplIDKey(tmplIDCourseRemind, courseSn)
}
