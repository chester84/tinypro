package i18n

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"tinypro/common/lib/redis/storage"
	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type LangItem map[string]string
type LangConfig map[string]LangItem

var langSupportConf = map[string]string{
	types.LangEnUS: types.LanguageTypeEnglishDisplay,
	types.LangZhCN: types.LanguageTypeChineseDisplay,
}

// 翻译函数,如果没有配置语言包,显示原始数据
func T(lang, src string) string {
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	langT := types.Short2LangType(lang)
	key := stringMappingKey(src, langT)

	// read from redis
	v, err := redis.String(storageClient.Do("HGET", rdsKeyStringMapping, key))
	if err != nil {
		if err.Error() != types.RedigoNilReturned {
			logs.Error("[i18n->T] redis> HGET %s %s, err: %v", rdsKeyStringMapping, key, err)
		} else {
			logs.Notice("redis> HGET %s %s, err: %v", rdsKeyStringMapping, key, err)
		}
	} else {
		return v
	}

	// read from db
	one := models.I18nStringMapping{
		SrcString: src,
		Language:  langT,
	}
	// 未添加时 err不为nil
	dst, err := one.GetDstString()
	if err != nil {
		logs.Warn("[T] src: %v, lang: %v, err: %v", src, lang, err)
		dst = src
	}

	_, err = storageClient.Do("HSET", rdsKeyStringMapping, key, dst)
	if err != nil {
		logs.Error("[i18n->T] redis> HSET %s %s %s, err: %v", rdsKeyStringMapping, key, dst, err)
	}
	return dst
}

// 取系统支持的语言包
func LangSupportConf() map[string]string {
	return langSupportConf
}

func IsExist(lang string) bool {
	if _, ok := langSupportConf[lang]; ok {
		return true
	}

	return false
}

func StringMappingList(condBox map[string]interface{}, page int, pageSize int) (list []models.I18nStringMapping, total int64, err error) {
	obj := models.I18nStringMapping{}
	o := orm.NewOrm()

	var sql string
	var whereBox []string

	sqlCount := "SELECT COUNT(id) AS total"
	sqlQuery := fmt.Sprintf("SELECT *")
	from := fmt.Sprintf("FROM %s", obj.TableName())

	var where string
	if v, ok := condBox["src_string"]; ok {
		whereBox = append(whereBox, fmt.Sprintf(`src_string = '%s'`, libtools.Escape(v.(string))))
	}

	if v, ok := condBox["language"]; ok {
		whereBox = append(whereBox, fmt.Sprintf(`language = %v`, v.(types.LanguageTypeEnum)))
	}

	if len(whereBox) > 0 {
		where = strings.Join(whereBox, " AND ")
		where = " where " + where
	}

	var orderBy = "ORDER BY id DESC"
	if v, ok := condBox["order_by"]; ok {
		orderBy = v.(string)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = types.DefaultPagesize
	}
	offset := (page - 1) * pageSize

	limit := fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)

	sql = fmt.Sprintf("%s %s %s", sqlCount, from, where)
	errCount := o.Raw(sql).QueryRow(&total)
	if errCount != nil {
		logs.Error("[COUNT-SQL] db get exception, SQL: %s, err: %v", sql, errCount)
	}

	sql = fmt.Sprintf("%s %s %s %s %s", sqlQuery, from, where, orderBy, limit)
	logs.Debug("[SELECT-SQL]: %s", sql)
	_, err = o.Raw(sql).QueryRows(&list)
	if err != nil {
		logs.Error("[SELECT-SQL] db get exception, SQL: %s, err: %v", sql, err)
	}

	return
}

func UpdateMappingCache(src string, lang types.LanguageTypeEnum, dst string) {
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	key := stringMappingKey(src, lang)
	_, err := storageClient.Do("HSET", rdsKeyStringMapping, key, dst)
	if err != nil {
		logs.Error("[UpdateMappingCache] redis> HSET %s %s %s, err: %v", rdsKeyStringMapping, key, dst, err)
	}

	return
}

func InsertOrUpdate(opUid int64, one models.I18nStringMapping) error {
	var err error
	o := orm.NewOrm()

	check := models.I18nStringMapping{}
	err = o.QueryTable(check.TableName()).Filter("language", one.Language).Filter("src_string", one.SrcString).One(&check)
	if err != nil {
		// 新增
		if err.Error() != types.EmptyOrmStr {
			logs.Error("[InsertOrUpdate] check data get exception, err: %v", err)
			return err
		}

		_, err = models.OrmInsert(&one)
		if err != nil {
			logs.Error("[InsertOrUpdate] insert exception, data: %#v err: %v", one, err)
			return err
		}
	} else {
		origin := check

		check.DstString = one.DstString
		check.LastOpBy = opUid
		check.LastOpAt = libtools.GetUnixMillis()
		_, err = models.OrmUpdate(&check, []string{"dst_string", "last_op_by", "last_op_at"})
		if err != nil {
			logs.Error("[InsertOrUpdate] update exception, data: %#v, err: %v", check, err)
			return err
		}

		models.OpLogWrite(opUid, check.Id, models.OpCodeUpI18nMapping, check.TableName(), origin, check)
	}

	UpdateMappingCache(one.SrcString, one.Language, one.DstString)

	return nil
}
