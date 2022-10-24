package models

import "github.com/beego/beego/v2/client/orm"

func init() {
	orm.RegisterModel(new(CasbinRule))
}

const CASBIN_RULE_TABLENAME = "casbin_rule"

type CasbinRule struct {
	Id    int64 `orm:"pk;"`
	PType string
	V0    string `json:"v0"`
	V1    string `json:"v1"`
	V2    string
	V3    string
	V4    string
	V5    string
}

func (r *CasbinRule) TableName() string {
	return CASBIN_RULE_TABLENAME
}
