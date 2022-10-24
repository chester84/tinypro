// 定制 mysql

package cmysql

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"

	_ "github.com/go-sql-driver/mysql"

	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func init() {
	if !libtools.IsProductEnv() {
		orm.Debug = true
	}

	dbType, _ := config.String("db_type")
	dbCharset, _ := config.String("db_charset")

	_ = orm.RegisterDriver(dbType, orm.DRMySQL)

	dbHost, _ := config.String("db_host")
	dbPort, _ := config.String("db_port")
	dbName, _ := config.String("db_name")
	dbUser, _ := config.String("db_user")
	dbPwd, _ := config.String("db_pwd")

	fmt.Printf("types.OrmDataBase: %s\n", types.OrmDataBase)
	_ = orm.RegisterDataBase(types.OrmDataBase, dbType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName, dbCharset))
	db, _ := orm.GetDB(types.OrmDataBase)
	db.SetConnMaxLifetime(time.Hour)
	orm.SetMaxIdleConns(types.OrmDataBase, 128)
	orm.SetMaxOpenConns(types.OrmDataBase, 128)
}
