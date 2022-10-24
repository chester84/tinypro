package models

import (
	"tinypro/common/types"

	"github.com/beego/beego/v2/client/orm"
)

// 	orm.RegisterModel(new(I18nStringMapping))

const I18N_STRING_MAPPING_TABLENAME = "i18n_string_mapping"

type I18nStringMapping struct {
	Id        int64 `orm:"pk;" json:"id"`
	SrcString string
	DstString string
	Language  types.LanguageTypeEnum `json:"language"`
	CreatedBy int64
	CreatedAt int64
	LastOpAt  int64
	LastOpBy  int64
}

func (r *I18nStringMapping) TableName() string {
	return I18N_STRING_MAPPING_TABLENAME
}

func (r *I18nStringMapping) GetDstString() (string, error) {
	o := orm.NewOrm()
	err := o.QueryTable(r.TableName()).Filter("SrcString", r.SrcString).Filter("Language", r.Language).One(r)
	return r.DstString, err
}
