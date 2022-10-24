package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"tinypro/common/models"
	"tinypro/common/types"

	"tinypro/common/lib/redis/storage"
	"tinypro/common/pkg/weixin"
	"github.com/chester84/libtools"

	"github.com/beego/beego/v2/core/logs"
)

type WXNotifyController struct {
	SuccessResp string
	FailResp    string
	OpenApiController
}

func (c *WXNotifyController) Prepare() {
	// 调用上一级的 Prepare 方
	c.SuccessResp = `
	<xml>
        <return_code>SUCCESS</return_code>
        <return_msg>OK</return_msg>
	</xml>
	`

	c.FailResp = `
	<xml>
        <return_code>FAIL</return_code>
        <return_msg>FAIL</return_msg>
	</xml>
	`

	c.OpenApiController.Prepare()
}

func (c *WXNotifyController) EntrustNotify() {
	c.Ctx.Output.Header("Content-Type", "application/xml; charset=utf-8")
	c.Ctx.Output.Status = 200

	type req struct {
		ReturnCode   string `xml:"return_code" json:"return_code"`
		ResultCode   string `xml:"result_code" json:"result_code"`
		Sign         string `xml:"sign" json:"sign"`
		MchId        string `xml:"mch_id" json:"mch_id"`
		ContractCode string `xml:"contract_code" json:"contract_code"`
		Openid       string `xml:"openid" json:"openid"`
		PlanId       string `xml:"plan_id" json:"plan_id"`
		ChangeType   string `xml:"change_type" json:"change_type"`
		OperateTime  string `xml:"operate_time" json:"operate_time"`
		ContractId   string `xml:"contract_id" json:"contract_id"`
	}
	reqPara := req{}
	err := xml.Unmarshal(c.Ctx.Input.RequestBody, &reqPara)
	logs.Warn("EntrustNotify c.Ctx.Input.RequestBody is %s", string(c.Ctx.Input.RequestBody))
	logs.Warn("EntrustNotify reqPara is %#v", reqPara)
	if err != nil {
		logs.Error("[UnifiedOrderPayNotify] xml.Unmarshal err: %v", err)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()
	// 发生并发回调，只能第一个处理，其余的全部退出
	lockKey := fmt.Sprintf("%s:%s", "wx-entrust-notify", reqPara.Openid)
	lock, err := storageClient.Do("SET", lockKey, libtools.GetUnixMillis(), "NX")
	if err != nil || lock == nil {
		logs.Error("[EntrustNotify] is working, so, I will exit, err: %v", err)
		c.Ctx.WriteString(c.FailResp)
		return
	}
	defer func() {
		_, _ = storageClient.Do("DEL", lockKey)
	}()

	bson, err := json.Marshal(reqPara)
	if err != nil {
		logs.Error("[EntrustNotify] json.Marshal(reqPara) err: %v", err)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	var paraMap = make(map[string]interface{})
	err = json.Unmarshal(bson, &paraMap)
	if err != nil {
		logs.Error("[EntrustNotify] json.Unmarshal(bson, &paraMap) err: %v", err)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	//验签
	calcSign := weixin.Signature(paraMap, weixin.MerchantKey())
	if calcSign != reqPara.Sign {
		logs.Error("[EntrustNotify] calcSign != reqPara.Sign calcSign:%s, Sign:%s", calcSign, reqPara.Sign)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	c.Ctx.WriteString(c.SuccessResp)
	return
}

func (c *WXNotifyController) EduUnifiedOrderPayNotify() {
	c.Ctx.Output.Header("Content-Type", "application/xml; charset=utf-8")
	c.Ctx.Output.Status = 200

	reqPara := weixin.UnifiedOrderPayNotifyRequest{}
	err := xml.Unmarshal(c.Ctx.Input.RequestBody, &reqPara)
	logs.Warn("EduUnifiedOrderPayNotify c.Ctx.Input.RequestBody is %s", string(c.Ctx.Input.RequestBody))
	logs.Warn("EduUnifiedOrderPayNotify reqPara is %#v", reqPara)
	if err != nil {
		logs.Error("[UnifiedOrderPayNotify] xml.Unmarshal err: %v", err)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()
	// 发生并发回调，只能第一个处理，其余的全部退出
	lockKey := fmt.Sprintf("%s:%s", "edu-wx-unified-order-pay-notify", reqPara.Openid)
	lock, err := storageClient.Do("SET", lockKey, libtools.GetUnixMillis(), "NX")
	if err != nil || lock == nil {
		logs.Error("[EduUnifiedOrderPayNotify] is working, so, I will exit, err: %v", err)
		c.Ctx.WriteString(c.FailResp)
		return
	}
	defer func() {
		_, _ = storageClient.Do("DEL", lockKey)
	}()

	bson, err := json.Marshal(reqPara)
	if err != nil {
		logs.Error("[EduUnifiedOrderPayNotify] json.Marshal(reqPara) err: %v", err)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	var paraMap = make(map[string]interface{})
	err = json.Unmarshal(bson, &paraMap)
	if err != nil {
		logs.Error("[EduUnifiedOrderPayNotify] json.Unmarshal(bson, &paraMap) err: %v", err)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	//验签
	calcSign := weixin.Signature(paraMap, weixin.MerchantKey())
	if calcSign != reqPara.Sign {
		logs.Error("[EduUnifiedOrderPayNotify] calcSign != reqPara.Sign calcSign:%s, Sign:%s", calcSign, reqPara.Sign)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	//// 业务
	eduPayId, _ := libtools.Str2Int64(reqPara.OutTradeNo)
	var eduPay models.EduPayment
	err = models.OrmOneByPkId(eduPayId, &eduPay)
	if err != nil {
		logs.Error("[EduUnifiedOrderPayNotify] payment data empty, payId: %d, err: %#v", eduPayId, err)
		c.Ctx.WriteString(c.FailResp)
		return
	}

	eduPay.Status = types.PaymentStatusSuccess
	eduPay.WxTransactionId = reqPara.TransactionId
	eduPay.CallbackAt = libtools.GetUnixMillis()
	_, err = models.OrmUpdate(&eduPay, []string{"status", "wx_transaction_id", "callback_at"})
	if err != nil {
		logs.Error("[UnifiedOrderPayNotify] recharge OrmUpdate fail, payId: %d, err: %#v", eduPayId, err)
	}

	invokeLog := models.WxCallback{
		ReqUrl:    fmt.Sprintf("%s%s", libtools.InternalApiDomain(), c.Ctx.Input.URL()),
		ReqType:   types.WxCallbackUnifiedOrder,
		ReqParams: string(c.Ctx.Input.RequestBody),
		RespCode:  200,
		Resp:      c.SuccessResp,
		CreatedAt: libtools.GetUnixMillis(),
	}
	_, err = models.OrmInsert(&invokeLog)
	if err != nil {
		logs.Error("[EduUnifiedOrderPayNotify] WxCallback OrmInsert invokeLog err %#v", err)
	}

	c.Ctx.WriteString(c.SuccessResp)
	return
}

func (c *WXNotifyController) MsgCallback() {
	echoStr := c.GetString(`echostr`)

	invokeLog := models.WxCallback{
		ReqUrl:    fmt.Sprintf("%s%s", libtools.InternalApiDomain(), c.Ctx.Input.URL()),
		ReqType:   types.WxMsgCallback,
		ReqParams: string(c.Ctx.Input.RequestBody),
		RespCode:  200,
		Resp:      echoStr,
		CreatedAt: libtools.GetUnixMillis(),
	}

	_, err := models.OrmInsert(&invokeLog)
	if err != nil {
		logs.Error("[MsgCallback] WxCallback OrmInsert invokeLog err %#v", err)
	}

	c.Ctx.WriteString(echoStr)
	return
}
