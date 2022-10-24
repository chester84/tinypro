package esclient

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/olivere/elastic/v7"
)

var (
	instance *elastic.Client
)

func init() {
	initClient()
}

func initClient() {
	// 注册`esclient`
	esHosts, _ := config.String("es_hosts")
	var hosts []string
	_ = json.Unmarshal([]byte(esHosts), &hosts)

	if len(hosts) == 0 {
		logs.Error("Es initClient host is empty")
		return
	}
	logs.Info("es.initClient host:%v", hosts)

	var err error
	instance, err = elastic.NewClient(
		elastic.SetURL(hosts...),
		elastic.SetSniff(false))

	if err != nil {
		logs.Error("[esclient] client init err:%v.", err)
		panic(err)
	}
}

func Client() *elastic.Client {
	if instance == nil {
		initClient()
	}

	return instance
}
