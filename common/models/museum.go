package models

import "tinypro/common/types"

// 	orm.RegisterModel(new(Museum))

const MUSEUM_TABLENAME = "museum"

type Museum struct {
	Id            int64                  `orm:"pk;"`
	Subject       string                 // 主题
	Image         string                 // 头图
	Summary       string                 // 摘要
	Content       string                 // 具休内容
	DiggCount     int                    // 点赞数
	CommentCount  int                    // 评论数
	FavoriteCount int                    // 收藏数
	Pin           int                    // 置顶
	Weight        int                    // 排序权重
	Status        types.StatusCommonEnum // 1: 上线; 0: 下线
	CreatedBy     int64                  // 记录创建者
	CreatedAt     int64                  // 记录创建时间
	LastOpBy      int64                  // 最后操作员
	LastOpAt      int64                  // 最后操作时间
}

func (r *Museum) TableName() string {
	return MUSEUM_TABLENAME
}
