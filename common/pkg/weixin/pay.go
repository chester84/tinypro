package weixin

import (
	"encoding/xml"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"tinypro/common/models"
	"tinypro/common/pkg/account"
	"tinypro/common/pkg/payment"
	"time"

	"github.com/chester84/libtools"
	"tinypro/common/types"
)

const (
	apiResponseSuccess = `SUCCESS`

	PayTradeTypeJsApi  types.WxPayTradeTypeEnum = `JSAPI`
	PayTradeTypeNative types.WxPayTradeTypeEnum = `NATIVE`
	PayTradeTypeApp    types.WxPayTradeTypeEnum = `APP`
	PayTradeTypeH5     types.WxPayTradeTypeEnum = `MWEB`
)

type PayCommonEmbed struct {
	XMLName  xml.Name `json:"-" xml:"xml"`
	Appid    string   `json:"appid" xml:"appid"`
	MchId    string   `json:"mch_id" xml:"mch_id"`
	NonceStr string   `json:"nonce_str" xml:"nonce_str"`
	Sign     string   `json:"sign" xml:"sign"` // sign_type, 默认为MD5
	OpenId   string   `json:"openid" xml:"openid"`
}

type PayResponseEmbed struct {
	//XMLName    xml.Name `json:"-" xml:"xml"`
	ReturnCode string `json:"return_code" xml:"return_code"`
	ReturnMsg  string `json:"return_msg" xml:"return_msg"`
	ResultCode string `json:"result_code" xml:"result_code"`
	ErrCode    string `json:"err_code,omitempty" xml:"err_code,omitempty"`
	ErrCodeDes string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty"`
}

func (r *PayResponseEmbed) IsSuccess() bool {
	return r.ReturnCode == apiResponseSuccess && r.ResultCode == apiResponseSuccess
}

// docs: https://pay.weixin.qq.com/wiki/doc/api/wxpay_v2/open/chapter8_1.shtml
type PayUnifiedOrderRequest struct {
	PayCommonEmbed

	Body types.XmlCData `json:"body" xml:"body"`
	//Body string `json:"body" xml:"body"`

	OutTradeNo     string `json:"out_trade_no" xml:"out_trade_no"`
	TotalFee       int64  `json:"total_fee" xml:"total_fee"` // 订单总金额，单位为分
	SpBillCreateIp string `json:"spbill_create_ip" xml:"spbill_create_ip"`
	NotifyUrl      string `json:"notify_url" xml:"notify_url"`

	TradeType types.WxPayTradeTypeEnum `json:"trade_type" xml:"trade_type"`
}

type PayUnifiedOrderResponse struct {
	PayResponseEmbed

	TradeType types.WxPayTradeTypeEnum `json:"trade_type" xml:"trade_type"`

	PrepayId string `json:"prepay_id" xml:"prepay_id"`
	CodeUrl  string `json:"code_url" xml:"code_url"`
}

func payUnifiedOrderApi() string {
	var api string
	if libtools.IsProductEnv() {
		api = `https://api.mch.weixin.qq.com/pay/unifiedorder`
	} else {
		api = `https://api.mch.weixin.qq.com/sandboxnew/pay/unifiedorder`
	}

	return api
}

func fixSandboxFee(fee int64) int64 {
	if libtools.IsProductEnv() {
		return fee
	} else {
		return 101
		//return 888 //!  文档上有,但实际无效 [沙箱支付金额(888)无效，请检查需要验收的case]
		//return 301
	}
}

type PaySampleRes struct {
	AppId     string `json:"app_id"`
	TimeStamp string `json:"time_stamp"`
	NonceStr  string `json:"nonce_str"`
	Package   string `json:"package"`
	PrepayId  string `json:"prepay_id"`
	SignType  string `json:"sign_type"`
	PaySN     int64  `json:"pay_sn,string"`
	CodeUrl   string `json:"code_url"`
	PaySign   string `json:"pay_sign"`
}

//!!! 沙箱环境金额只能传 101

func EduCreatePaymentOrder(userId int64, body string, ip string, eduReq types.EduPayCreateOrderReq) (payObj PaySampleRes, err error) {
	var api = payUnifiedOrderApi()
	var req = PayUnifiedOrderRequest{}

	openId, err := account.GetWxOpenId(userId, eduReq.AppSN)
	if err != nil || openId == "" {
		logs.Error("[EduCreatePaymentOrder] can get wx openid, userId: %d, appSN: %d", userId, eduReq)
		return
	}

	eduPay, err := payment.EduCreateOne(userId, eduReq)
	if err != nil {
		return
	}

	req.Appid = AppID(eduReq.AppSN)
	req.MchId = MerchantID()
	req.NonceStr = libtools.GenerateRandomStr(16)
	req.OpenId = openId
	req.Body = types.XmlCData(body)
	//reqs.Body = body
	req.OutTradeNo = fmt.Sprintf(`%d`, eduPay.Id)
	req.TotalFee = fixSandboxFee(eduReq.TotalAmount)
	req.SpBillCreateIp = ip
	req.TradeType = eduReq.TradeType
	req.NotifyUrl = fmt.Sprintf(`%s/wx/edu/unified-order/notify`, libtools.InternalApiDomain())

	reqMap := libtools.Struct2MapV3(req)
	sign := Signature(reqMap, MerchantKey())
	req.Sign = sign

	reqXmlB, err := xml.Marshal(req)
	if err != nil {
		logs.Error("[EduCreatePaymentOrder] xml encode exception, reqs: %#v, err: %v", req, err)
		return
	}

	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodPOST, api, nil, string(reqXmlB), libtools.DefaultHttpTimeout())
	if err != nil || httpCode != 200 {
		logs.Error("[EduCreatePaymentOrder] call api error. api: %s, reqXML: %s, err: %v", api, string(reqXmlB), err)
		return
	}

	var resp PayUnifiedOrderResponse
	err = xml.Unmarshal(httpBody, &resp)
	if err != nil {
		logs.Error("[EduCreatePaymentOrder] xml decode error, requestXML: %s, responseXML: %s, err: %v", string(reqXmlB), string(httpBody), err)
		return
	}

	logs.Debug("[EduCreatePaymentOrder] api: %s, reqs: %s, res: %s", api, string(reqXmlB), string(httpBody))

	if !resp.IsSuccess() {
		err = fmt.Errorf(`no valid data`)
		return
	}

	payObj.AppId = AppID(eduReq.AppSN)
	payObj.TimeStamp = fmt.Sprintf(`%d`, time.Now().Unix())
	payObj.NonceStr = req.NonceStr
	payObj.Package = fmt.Sprintf(`prepay_id=%s`, resp.PrepayId)
	payObj.PrepayId = resp.PrepayId
	payObj.SignType = "MD5"
	payObj.PaySN = eduPay.Id
	payObj.CodeUrl = resp.CodeUrl

	// see: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
	payMap := map[string]interface{}{
		"appId": payObj.AppId,

		"timeStamp": payObj.TimeStamp,
		//"timestamp": payObj.TimeStamp, // 微支付坑: 服务端此参数要大写

		"nonceStr": payObj.NonceStr,
		"package":  payObj.Package,
		"signType": payObj.SignType,
	}
	//payObj.PaySign = Signature(payMap, Secret(appSN)) // 微支付坑: 需要商户的密钥,而非公众号/小程序密钥
	payObj.PaySign = Signature(payMap, MerchantKey())

	eduPay.WxPrepayId = resp.PrepayId
	_, err = models.OrmUpdate(&eduPay, []string{"WxPrepayId"})
	if err != nil {
		logs.Error("[CreatePaymentOrder] update prepay_id exception, ip: %s, accountID: %d, obj: %#v, err: %v",
			ip, userId, eduPay, err)
	}

	return
}

type UnifiedOrderPayNotifyRequest struct {
	Appid              string `json:"appid" xml:"appid"`
	SubAppid           string `json:"sub_appid" xml:"sub_appid"`
	Attach             string `json:"attach" xml:"attach"` //not required
	BankType           string `json:"bank_type" xml:"bank_type"`
	FeeType            string `json:"fee_type" xml:"fee_type"`
	CashFeeType        string `json:"cash_fee_type,omitempty" xml:"cash_fee_type,omitempty"`               //not required
	CashFee            string `json:"cash_fee,omitempty" xml:"cash_fee,omitempty"`                         //not required
	DeviceInfo         string `json:"device_info,omitempty" xml:"device_info,omitempty"`                   //not required
	ErrCode            string `json:"err_code,omitempty" xml:"err_code,omitempty"`                         //not required
	ErrCodeDes         string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty"`                 //not required
	SettlementTotalFee string `json:"settlement_total_fee,omitempty" xml:"settlement_total_fee,omitempty"` //not required
	IsSubscribe        string `json:"is_subscribe" xml:"is_subscribe"`
	MchId              string `json:"mch_id" xml:"mch_id"`
	NonceStr           string `json:"nonce_str" xml:"nonce_str"`
	Openid             string `json:"openid" xml:"openid"`
	OutTradeNo         string `json:"out_trade_no" xml:"out_trade_no"`
	ResultCode         string `json:"result_code" xml:"result_code"`
	ReturnCode         string `json:"return_code" xml:"return_code"`
	ReturnMsg          string `json:"return_msg" xml:"return_msg"` //not required
	Sign               string `json:"sign" xml:"sign"`
	TimeEnd            string `json:"time_end" xml:"time_end"`
	TotalFee           string `json:"total_fee" xml:"total_fee"`
	CouponFee          string `json:"coupon_fee" xml:"coupon_fee"`
	CouponCount        string `json:"coupon_count" xml:"coupon_count"`
	CouponType         string `json:"coupon_type" xml:"coupon_type"`
	CouponId           string `json:"coupon_id" xml:"coupon_id"`
	TradeType          string `json:"trade_type" xml:"trade_type"`
	TransactionId      string `json:"transaction_id" xml:"transaction_id"`
}

// https://pay.weixin.qq.com/wiki/doc/api/wxpay_v2/open/chapter8_4.shtml
type ApplyRefundRequest struct {
	PayCommonEmbed

	OutTradeNo  string `json:"out_trade_no" xml:"out_trade_no"`
	OutRefundNo string `json:"out_refund_no" xml:"out_refund_no"`
	TotalFee    int64  `json:"total_fee" xml:"total_fee"`   // 订单总金额，单位为分
	RefundFee   int64  `json:"refund_fee" xml:"refund_fee"` // 退款总金额，单位为分
	RefundDesc  string `json:"refund_desc,omitempty" xml:"refund_desc,omitempty"`
	NotifyUrl   string `json:"notify_url" xml:"notify_url"`
}

type ApplyRefundResponse struct {
	PayResponseEmbed

	TransactionId string `json:"transaction_id" xml:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no"`
	OutRefundNo   string `json:"out_refund_no" xml:"out_refund_no"`
	RefundId      string `json:"refund_id" xml:"refund_id"`
	RefundFee     int64  `json:"refund_fee" xml:"refund_fee"` // 退款总金额，单位为分
	TotalFee      int64  `json:"total_fee" xml:"total_fee"`   // 订单总金额，单位为分
	CashFee       int64  `json:"cash_fee" xml:"cash_fee"`     // 现金支付金额，单位为分
}

func GetSandBoxPayMerchantKey() (key string) {
	type SandBoxPayMerchantKeyReq struct {
		XMLName  xml.Name `json:"-" xml:"xml"`
		MchID    string   `json:"mch_id" xml:"mch_id"`
		NonceStr string   `json:"nonce_str" xml:"nonce_str"`
		Sign     string   `json:"sign" xml:"sign"`
	}

	req := SandBoxPayMerchantKeyReq{}
	req.MchID = MerchantID()
	req.NonceStr = libtools.GenerateRandomStr(6)

	reqMap := libtools.Struct2MapV3(req)
	sign := Signature(reqMap, MerchantProdKey())
	req.Sign = sign

	reqXmlB, err := xml.Marshal(req)
	if err != nil {
		logs.Error("[getSandBoxPayMerchantKey] xml encode exception, reqs: %#v, err: %v", req, err)
		return
	}

	api := "https://api.mch.weixin.qq.com/sandboxnew/pay/getsignkey"
	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodPOST, api, nil, string(reqXmlB), libtools.DefaultHttpTimeout())
	if err != nil || httpCode != 200 {
		logs.Error("[getSandBoxPayMerchantKey] call api error. api: %s, reqXML: %s, err: %v", api, string(reqXmlB), err)
		return
	}

	type SandBoxPayMerchantKeyResp struct {
		XMLName        xml.Name `json:"-" xml:"xml"`
		ReturnCode     string   `json:"return_code" xml:"return_code"`
		ReturnMsg      string   `json:"return_msg" xml:"return_msg"`
		SandboxSignKey string   `json:"sandbox_signkey" xml:"sandbox_signkey"`
	}

	resp := SandBoxPayMerchantKeyResp{}
	err = xml.Unmarshal(httpBody, &resp)
	if err != nil {
		logs.Error("[getSandBoxPayMerchantKey] xml decode error, requestXML: %s, responseXML: %s, err: %v", string(reqXmlB), string(httpBody), err)
		return
	}

	logs.Debug("[getSandBoxPayMerchantKey] api: %s, reqs: %s, res: %s", api, string(reqXmlB), string(httpBody))
	if resp.ReturnCode != "SUCCESS" && resp.ReturnMsg != "ok" {
		logs.Error("[getSandBoxPayMerchantKey] return fail, requestXML: %s, responseXML: %s, err: %v", string(reqXmlB), string(httpBody), err)
		return
	} else {
		key = resp.SandboxSignKey
	}

	return
}
