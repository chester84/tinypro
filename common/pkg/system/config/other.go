package config

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/types"
)

func PlaceApiSwitch() int {
	st, err := ValidItemInt(PlaceApiSwitchKey)
	if err != nil {
		_, errDb := Create(PlaceApiSwitchKey, "0", types.SystemConfigItemTypeInt64, 200, `Place Api 开关, 1: 开启; 0: 关闭`, 0)
		if errDb != nil {
			logs.Error("[PlaceApiSwitch] init get exception, key: %s, err: %v", PlaceApiSwitchKey, errDb)
		}
	}

	return st
}

func DefaultSearchRadius() int {
	dftValue, err := ValidItemInt(DefaultSearchRadiusKey)
	if err != nil {
		dftValue = 5000
		_, errDb := Create(DefaultSearchRadiusKey, fmt.Sprintf(`%d`, dftValue), types.SystemConfigItemTypeInt64, 201, `App 默认搜索半径`, 0)
		if errDb != nil {
			logs.Error("[DefaultSearchRadius] init get exception, key: %s, err: %v", DefaultSearchRadiusKey, errDb)
		}
	}

	return dftValue
}

func ExchangeRateIndianRupee2USDollar() float64 {
	dfValue, err := ValidItemFloat64(ExchangeRateIndianRupee2USDollarKey)
	if err != nil || dfValue < 0.001 {
		dfValue = 0.013
		_, errDB := Create(ExchangeRateIndianRupee2USDollarKey, fmt.Sprintf(`%f`, dfValue), types.SystemConfigItemTypeFloat64, 202, `印度卢比兑美元的汇率`, 0)

		if errDB != nil {
			logs.Error("[ExchangeRateIndianRupee2USDollar] init get exception, key: %s, err: %v", ExchangeRateIndianRupee2USDollarKey, errDB)
		}
	}

	return dfValue
}

func PayStrategy() string {
	dfValue := ValidItemString(PayStrategyKey)
	if dfValue == "" {
		dfValue = `{
    "pay_strategy": "round-robin",
    "random": [
        {"name": "payDockpay", "rate": 80},
 		{"name": "payWintec", "rate": 20}
    ],
    "round-robin": ["payerMax"]
}
`
		/*
			{"pay_strategy":"round-robin","round-robin":["payXPay"],"random":[{"name":"payShareIt","rate":80},{"name":"payDockpay","rate":20}]}
		*/
		_, errDB := Create(PayStrategyKey, fmt.Sprintf(`%s`, dfValue), types.SystemConfigItemTypeString, 100, `第三方支付选择`, 0)

		if errDB != nil {
			logs.Error("[PayStrategy] init get exception, key: %s, err: %v", PayStrategyKey, errDB)
		}
	}

	return dfValue
}
