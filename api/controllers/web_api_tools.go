package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/cerror"
	"tinypro/common/pkg/qqlbs"
	"tinypro/common/pkg/system/config"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type WebApiToolsController struct {
	WebApiBaseController
}

func (c *WebApiToolsController) Prepare() {
	// 调用上一级的 Prepare 方
	c.WebApiBaseController.Prepare()
}

func (c *WebApiToolsController) Address2Geo() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"address": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	type reqT struct {
		Address string `json:"address"`
	}

	var req reqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil || req.Address == "" {
		logs.Error("[Address2Geo] parse request get exception, ip: %s, accountID: %d reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	geo, err := qqlbs.Address2Geocoder(req.Address)
	if err != nil {
		logs.Error("[Address2Geo] call api get exception, ip: %s, accountID: %d reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.CallApiFail)
		return
	}

	type resT struct {
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	}

	res := resT{
		Longitude: fmt.Sprintf(`%f`, geo.Lng),
		Latitude:  fmt.Sprintf(`%f`, geo.Lat),
	}

	c.SuccessResponse(res)
}

func (c *WebApiToolsController) Geo2Address() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"longitude": true,
		"latitude":  true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	type reqT struct {
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	}

	var req reqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil || req.Latitude == "" || req.Longitude == "" {
		logs.Error("[Geo2Address] parse request get exception, ip: %s, accountID: %d reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	_, lngE := libtools.Str2Float64(req.Longitude)
	_, latE := libtools.Str2Float64(req.Latitude)
	if lngE != nil || latE != nil {
		logs.Error("[Geo2Address] request out of range, ip: %s, accountID: %d reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	res, err := qqlbs.Geo2Address(req.Longitude, req.Latitude)
	if err != nil {
		logs.Error("[Geo2Address] call api get exception, ip: %s, accountID: %d reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.CallApiFail)
	}

	c.SuccessResponse(res)
}

func (c *WebApiToolsController) MnpCheckVersion() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"app_version": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	type reqT struct {
		AppVersion string `json:"app_version"`
	}

	var req reqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil || req.AppVersion == "" {
		logs.Error("[MnpCheckVersion] parse request get exception, ip: %s, accountID: %d reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	var conf = types.MnpReleaseCtlConf{}
	confData := config.ValidItemString(types.MnpReleaseCtlKey)
	err = json.Unmarshal([]byte(confData), &conf)
	if err != nil {
		logs.Warning("[MnpCheckVersion] json decode get exception, confData: %s, err: %v", confData, err)
	}

	type resT struct {
		Hit        int    `json:"hit"`
		CtlVersion string `json:"ctl_version"`
	}

	var res = resT{
		CtlVersion: conf.AppVersion,
	}
	if req.AppVersion == conf.AppVersion {
		res.Hit = 1
	}

	c.SuccessResponse(res)
}
