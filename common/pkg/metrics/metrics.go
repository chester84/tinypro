package metrics

import (
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/chester84/libtools"
)

var (
	// !!!注意,此处有些硬编码
	appName, _ = config.String("appname")
	appNameExp = strings.Split(appName, "-")
	Namespace  = "jax"
	Subsystem  = appNameExp[1]
	Hostname   = libtools.Hostname()
)

// 初始化 web_request_total, counter 类型指标,表示接收 http 请求总次数
var WebRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "web_request_total",
		Help:      "Number of uri requests in total",
	},
	// 设置标签
	[]string{"hostname", "url"},
)

// web_request_duration, Histogram类型指标, bucket 代表 duration 的分布区间
var WebRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "web_request_duration",
		Help:      "Web request duration distribution",
		Buckets:   prometheus.LinearBuckets(0, 1, 60), // 0 ~ 60
	},
	// 设置标签
	[]string{"hostname", "url"},
)

var RealTimeOnlineNum = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Namespace:   Namespace,
		Subsystem:   Subsystem,
		Name:        "real_time_online_num",
		Help:        "Number of app real-time online.",
		ConstLabels: prometheus.Labels{"destination": "RealTimeOnlineNum"},
	})

func init() {
	// 注册监控指标
	prometheus.MustRegister(WebRequestTotal)
	prometheus.MustRegister(WebRequestDuration)
	prometheus.MustRegister(RealTimeOnlineNum)
}
