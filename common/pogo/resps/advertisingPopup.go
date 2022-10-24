package resps

type AccessPopupResp struct {
	AdvertiseUrl string    `json:"advertise_url"`
	CourseNode   FrontPage `json:"course_node"`
	IsAccess     int       `json:"is_access"`
}
