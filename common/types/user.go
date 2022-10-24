package types

type ApiUserBaseProfileResponse struct {
	Nickname   string `json:"nickname"`
	OpenAvatar string `json:"open_avatar"`
	ShareCode  string `json:"share_code"`
}

type ApiTwoTuple struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type ApiLoginUserInfoResponse struct {
	Nickname    string `json:"nickname"`
	OpenAvatar  string `json:"open_avatar"`
	AccessToken string `json:"access_token"`
}
