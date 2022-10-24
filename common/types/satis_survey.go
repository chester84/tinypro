package types

type SatisSurveyQ struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type SatisSurveyScore struct {
	S int    `json:"s"`
	D string `json:"d"`
}
