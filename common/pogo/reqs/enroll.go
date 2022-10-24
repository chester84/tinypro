package reqs

type EnrollReqT struct {
	CourseSN        string `json:"course_sn"`
	RealName        string `json:"real_name"`
	Mobile          string `json:"mobile"`
	Company         string `json:"company"`
	Position        string `json:"position"`
	Residence       string `json:"residence"`
	RecommendPerson string `json:"recommend_person"`
}

type SignInReqT struct {
	CourseSN string `json:"course_sn"`
}
