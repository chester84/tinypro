package reqs

type SubscribeTmpl struct {
	CourseSN string   `json:"course_sn"`
	Ids      []string `json:"ids"`
}
