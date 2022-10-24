package qqlbs

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
)

func Address2Geocoder(address string) (Location, error) {
	var location Location
	var err error

	api := fmt.Sprintf(`%s/?address=%s&key=%s`, geocoderApi, libtools.UrlEncode(address), ApiKey())
	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, api, nil, "", libtools.DefaultHttpTimeout())
	if err != nil {
		logs.Error("[Address2Geocoder] call api fail, api: %s, httpCode: %s, err: %v", api, httpCode, err)
		return location, err
	}

	res := GeocoderResponse{}
	err = json.Unmarshal(httpBody, &res)
	if err != nil {
		logs.Error("[Address2Geocoder] json decode exception, api: %s, httpCode: %s, httpBody: %s, err: %v", api, httpCode, string(httpBody), err)
		return location, err
	}

	if res.Status != 0 {
		logs.Error("[Address2Geocoder] response is not the expected value, api: %s, httpCode: %s, httpBody: %s, err: %v", api, httpCode, string(httpBody), err)
		err = fmt.Errorf(`call api status is on 0`)
		return location, err
	}

	location = res.Result.Location

	return location, err
}

func LocationIPGeo(ip string) (Location, error) {
	var location Location
	var err error

	api := fmt.Sprintf(`%s?ip=%s&key=%s`, locationIPApi, ip, ApiKey())
	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, api, nil, "", libtools.DefaultHttpTimeout())
	if err != nil {
		logs.Error("[LocationIPGeo] call api fail, api: %s, httpCode: %s, err: %v", api, httpCode, err)
		return location, err
	}

	res := LocationIPResponse{}
	err = json.Unmarshal(httpBody, &res)
	if err != nil {
		logs.Error("[LocationIPGeo] json decode exception, api: %s, httpCode: %s, httpBody: %s, err: %v", api, httpCode, string(httpBody), err)
		return location, err
	}

	if res.Status != 0 {
		logs.Error("[LocationIPGeo] response is not the expected value, api: %s, httpCode: %s, httpBody: %s, err: %v", api, httpCode, string(httpBody), err)
	}

	location = res.Result.Location

	return location, err
}

func Geo2Address(lnt, lat string) (AddressComponent, error) {
	var (
		address AddressComponent
		err     error
	)

	api := fmt.Sprintf(`%s/?location=%s,%s&key=%s`, geocoderApi, lat, lnt, ApiKey())
	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, api, nil, "", libtools.DefaultHttpTimeout())
	if err != nil {
		logs.Error("[Geo2Address] call api fail, api: %s, httpCode: %s, err: %v", api, httpCode, err)
		return address, err
	}

	res := Geo2AddressResponse{}
	err = json.Unmarshal(httpBody, &res)
	if err != nil {
		logs.Error("[Geo2Address] json decode exception, api: %s, httpCode: %s, httpBody: %s, err: %v", api, httpCode, string(httpBody), err)
		return address, err
	}

	if res.Status != 0 {
		logs.Error("[Geo2Address] response is not the expected value, api: %s, httpCode: %s, httpBody: %s, err: %v", api, httpCode, string(httpBody), err)
		err = fmt.Errorf(`call api status is on 0`)
		return address, err
	}

	address = res.Result.AddressComponent

	address.Address = res.Result.Address
	address.Recommend = res.Result.FormattedAddresses.Recommend

	return address, err
}
