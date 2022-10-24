package loginbiz

import (
	"github.com/beego/beego/v2/adapter/logs"
	"tinypro/common/models"
	"tinypro/common/pkg/accesstoken"
	"tinypro/common/pkg/account"
	"tinypro/common/pkg/weixin"
	"tinypro/common/pogo/reqs"
	"tinypro/common/types"
)

func WxOauthLoginOrRegister(req reqs.WxLoginReqT, ip string) (models.AppUser, weixin.SnsJsCode2SessionResponse, string, error) {
	var (
		authData    types.ApiOauthLoginReqT
		wxOauth2Res weixin.Oauth2SilentResponse
		wxUserInfo  weixin.SnsApiUserInfo
		wxJsSession weixin.SnsJsCode2SessionResponse
		user        models.AppUser
		accessToken string
		err         error
	)

	if req.AppSN == weixin.AppSNWxGzh {
		wxOauth2Res, err = weixin.Oauth2Silent(req.Code, req.AppSN)
		if err != nil || wxOauth2Res.Openid == "" {
			logs.Error("[WxOauth2Silent] ip: %s,err: %v", ip, err)
			return user, wxJsSession, accessToken, err
		}

		wxUserInfo, err = weixin.SnsUserInfo(wxOauth2Res.AccessToken, wxOauth2Res.Openid)
		if err != nil {
			logs.Error("[SnsUserInfo] ip: %s, err: %v", ip, err)
			return user, wxJsSession, accessToken, err
		}

		authData = types.ApiOauthLoginReqT{
			AppSN:        req.AppSN,
			OpenOauthPlt: types.OpenOauthWeChat,
			Nickname:     wxUserInfo.Nickname,
			OpenUserID:   wxUserInfo.Unionid,
			OpenAvatar:   wxUserInfo.HeadImgUrl,
			Gender:       wxUserInfo.FixWxSex(),
			WxOpenId:     wxUserInfo.Openid,
			Country:      wxUserInfo.Country,
			Province:     wxUserInfo.Province,
			City:         wxUserInfo.City,
		}
	} else {
		wxJsSession, err = weixin.SnsJsCode2Session(req.Code, req.AppSN)
		if err != nil {
			logs.Error("[SnsJsCode2Session]  err: %#v", err)
			return user, wxJsSession, accessToken, err
		}

		// 兼容 unionid 为空的情况
		if wxJsSession.Unionid == "" {
			wxJsSession.Unionid = wxJsSession.Openid
		}

		var nickname string
		if req.Reginfo.Username != "" {
			nickname = req.Reginfo.Username
		} else {
			nickname = account.GenGuestNickname()
		}

		authData = types.ApiOauthLoginReqT{
			AppSN:        req.AppSN,
			OpenOauthPlt: types.OpenOauthWeChat,
			Nickname:     nickname,
			//Nickname:     req.Reginfo.Username,
			OpenAvatar: req.Reginfo.Avatar,
			OpenUserID: wxJsSession.Unionid,
			Mobile:     req.Reginfo.Mobile,
			WxOpenId:   wxJsSession.Openid,
			Gender:     types.GenderUnknown,
		}
	}

	user, err = account.LoginOrRegister(authData, ip, types.WebApiVersion)
	if err != nil {
		return user, wxJsSession, accessToken, err
	}
	if user.Id <= 0 {
		//新用户，引导去注册
		return user, wxJsSession, accessToken, err
	}
	//account.WriteWxOpenId(user.Id, authData.AppSN, authData.WxOpenId)

	accessToken, err = accesstoken.GenTokenWithCache(user.Id, types.PlatformWxMiniProgram, ip)
	if err != nil {
		logs.Error("[GenTokenWithCache] gen token get exception, ip: %s, accountID: %d, err: %v",
			ip, user.Id, err)
	}
	return user, wxJsSession, accessToken, err
}
