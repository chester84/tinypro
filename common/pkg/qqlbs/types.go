// https://lbs.qq.com/service/webService/webServiceGuide/webServiceGeocoder
// https://lbs.qq.com/service/webService/webServiceGuide/webServiceIp
package qqlbs

import "github.com/chester84/libtools"

const (
	geocoderApi   = `https://apis.map.qq.com/ws/geocoder/v1`
	locationIPApi = `https://apis.map.qq.com/ws/location/v1/ip`
)

func ApiKey() string {
	var keyBox = []string{
		`K67BZ-BCILS-JMRO4-6KDOT-2RUGF-54FAX`,
		`PE7BZ-3IMCW-C37R2-ODE32-232CS-KVBRT`,
		`WEABZ-WMI6W-INXRK-RCQ5R-QKAPV-76F5U`,
		`ROOBZ-J4L6P-TTCD3-LUO22-Y5V52-E3F7I`,
		`7NIBZ-TQH62-VWKUW-C5JX6-LZC2E-XCF5B`,
	}

	return keyBox[libtools.GenerateRandom(0, 1000)%len(keyBox)]
}

type Location struct {
	Lat float64 `json:"lat"` // 纬度 latitude
	Lng float64 `json:"lng"` // 经度 longitude
}

type GeocoderResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Location Location `json:"location"`
	} `json:"result"`
}

type LocationIPResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Location Location `json:"location"`
	} `json:"result"`
}

type AddressComponent struct {
	Address      string `json:"address"`
	Recommend    string `json:"recommend"`
	Nation       string `json:"nation"`
	Province     string `json:"province"`
	City         string `json:"city"`
	District     string `json:"district"`
	Street       string `json:"street"`
	StreetNumber string `json:"street_number"`
}

type Geo2AddressResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`

	Result struct {
		Address string `json:"address"`

		FormattedAddresses struct {
			Recommend string `json:"recommend"`
		} `json:"formatted_addresses"`

		AddressComponent AddressComponent `json:"address_component"`
	} `json:"result"`
}
