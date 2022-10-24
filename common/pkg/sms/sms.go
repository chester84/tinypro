package sms

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
	tcsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"

	"tinypro/common/lib/redis/cache"
	"tinypro/common/models"
	"tinypro/common/pkg/tc"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func GetSmsByMobile(mobile string, serviceType types.SmsServiceTypeEnum) (models.Sms, error) {
	var obj = models.Sms{}
	var err error

	o := orm.NewOrm()

	err = o.QueryTable(obj.TableName()).
		Filter("mobile", mobile).
		Filter("verify_status", types.SmsVerifyStatusFail).
		Filter("service_type", serviceType).
		OrderBy("-id").
		Limit(1).
		One(&obj)

	return obj, err
}

func SendCode(mobiles, codes []string, tplID, ip string, serviceType types.SmsServiceTypeEnum) error {
	var err error

	mobilesPtr := make([]*string, 0)
	codesPtr := make([]*string, 0)

	for _, mobile := range mobiles {
		mobileChange := fmt.Sprintf("+86%s", mobile)
		mobilesPtr = append(mobilesPtr, &mobileChange)
	}

	for _, code := range codes {
		codesPtr = append(codesPtr, &code)
	}

	request := tcsms.NewSendSmsRequest()
	request.PhoneNumberSet = mobilesPtr
	request.TemplateID = &tplID
	request.TemplateParamSet = codesPtr
	sign := tc.GetBindMobileSign()
	request.Sign = &sign
	smsSdkAppid := tc.GetSmsSdkId()
	request.SmsSdkAppid = &smsSdkAppid

	resp, err := tc.SmsClient().SendSms(request)
	logs.Debug("resp is %#v", resp)
	jsonResp, _ := libtools.JsonEncode(resp)
	logs.Debug("jsonResp is %s", jsonResp)

	if err != nil ||
		resp == nil ||
		resp.Response == nil ||
		len(resp.Response.SendStatusSet) <= 0 ||
		*resp.Response.SendStatusSet[0].Code != "Ok" {
		err = fmt.Errorf("[SendCode] Send by tc API error: %#v", err)
		return err
	}

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	for i, mobile := range mobiles {
		smsVerifyCode := models.Sms{}
		smsVerifyCode.Vendor = types.Tencent
		smsVerifyCode.ServiceType = serviceType
		smsVerifyCode.Mobile = mobile
		smsVerifyCode.Code = codes[i]
		smsVerifyCode.Content = tplID
		smsVerifyCode.SmsType = types.SmsText
		smsVerifyCode.Expires = ExpireTime
		smsVerifyCode.SendStatus = types.SmsSendSuccess
		smsVerifyCode.Response = resp.ToJsonString()
		smsVerifyCode.ResponseId = *resp.Response.RequestId
		smsVerifyCode.Ip = ip
		smsVerifyCode.CallbackDeliveryStatus = types.SmsCallbackDeliverySuccess
		smsVerifyCode.VerifyStatus = types.SmsVerifyStatusFail
		smsVerifyCode.CreatedAt = time.Now()

		_, err = models.OrmInsert(&smsVerifyCode)
		if err != nil {
			logs.Error("[sms] Send insert exception, data: %#v, err: %v", smsVerifyCode, err)
			return err
		}

		rdsKey := buildSmsVerifyKey(serviceType, mobile)
		_, err = cacheClient.Do("SET", rdsKey, codes[i], "PX", ExpireTime)
		if err != nil {
			logs.Error("[SendCode] redis> SET %s %s PX %d, err: %v", rdsKey, codes[i], ExpireTime, err)
		}
	}

	return err
}

func VerifyCode(mobile, code string, serviceType types.SmsServiceTypeEnum) bool {
	var ret bool
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	rdsKey := buildSmsVerifyKey(serviceType, mobile)
	rdsCode, err := redis.String(cacheClient.Do("GET", rdsKey))
	if err != nil {
		logs.Error("[VerifyCode] redis> GET %s, err: %v", rdsKey, err)
		return ret
	}

	if rdsCode != code {
		logs.Warning("[VerifyCode] VerifyCode not match, rdsCode: %s, code: %s", rdsCode, code)
		return ret
	}

	go func() {
		smsVerifyData, _ := GetSmsByMobile(mobile, serviceType)
		if smsVerifyData.Id > 0 || smsVerifyData.Code == code {
			smsVerifyData.VerifyStatus = types.SmsVerifyStatusEnumSuccess
			smsVerifyData.UpdatedAt = time.Now()

			_, err := models.OrmUpdate(&smsVerifyData, []string{"verify_status", "updated_at"})
			if err != nil {
				logs.Error("[sms] VerifyCode update status err %#v", err)
			}
		} else {
			logs.Warning("[VerifyCode] last record no match, data: %#v, code: %s", smsVerifyData, code)
		}
	}()

	return true
}
