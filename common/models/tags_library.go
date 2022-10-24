package models

import (
	"github.com/beego/beego/v2/client/orm"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(TagsLibrary))
}

const TAGS_LIBRARY_TABLENAME = "tags_library"

type TagsLibrary struct {
	Id        int64 `orm:"pk;"`
	Name      string
	Weight    int    // 权重,用于标签排序
	Img       string // 标签图片
	Detail    string // 详细描述
	MarkNum   int    // 标签标记数量
	JoinNum   int    // 参与数
	Heat      int    // 热度
	CreatedBy int64
	CreatedAt int64
	Status    types.StatusCommonEnum // 1: 上线使用中; 2: 下线停用; 0: 异常
	LastOpBy  int64                  // 最后操作员
	LastOpAt  int64                  // 最后操作时间
}

func (r *TagsLibrary) TableName() string {
	return TAGS_LIBRARY_TABLENAME
}
