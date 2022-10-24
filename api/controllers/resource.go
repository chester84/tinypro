package controllers

import (
	"bytes"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"

	"tinypro/common/lib/redis/cache"
	"tinypro/common/pkg/metrics"
	"tinypro/common/pkg/tc"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type ResourceController struct {
	beego.Controller

	CurrentRouter string // 当前路由
	IP            string

	beginTime time.Time
}

func (c *ResourceController) Prepare() {
	c.IP = c.Ctx.Input.IP()
	c.CurrentRouter = `/open`

	// 量化接口并发量
	metrics.WebRequestTotal.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.CurrentRouter,
	}).Inc()

	c.beginTime = time.Now()

	if !libtools.IsProductEnv() {
		// 联调打印原始数据
		logs.Notice(">>> router: %s, ip: %s", c.Ctx.Input.URL(), c.IP)
	}
}

func (c *ResourceController) Finish() {
	// 量化接口性能
	duration := time.Since(c.beginTime)
	metrics.WebRequestDuration.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.CurrentRouter,
	}).Observe(duration.Seconds())
}

func (c *ResourceController) Resource() {
	rid := c.Ctx.Input.Param(":rid")

	if len(rid) < 4 {
		c.Ctx.Output.Status = 404
		_ = c.Ctx.Output.Body([]byte(types.Html404))
		return
	}

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	cKey := tc.GenTemporaryUrlRdsKey(rid)
	s3key, err := redis.String(cacheClient.Do("GET", cKey))
	if err != nil {
		if err != redis.ErrNil {
			logs.Error("[Resource] redis> GET %s, ip: %s, router: %s, err: %v",
				cKey, c.IP, c.Ctx.Input.URL(), err)
		}
		c.Ctx.Output.Status = 404
		_ = c.Ctx.Output.Body([]byte(types.Html404))
		return
	}

	resp, err := tc.HeadWithPrivate(s3key)
	if err != nil {
		logs.Warning("[Resource] can not find, ip: %s, router: %s, s3key: %s, err: %v",
			c.IP, c.Ctx.Input.URL(), s3key, err)
		c.Ctx.Output.Status = 404
		_ = c.Ctx.Output.Body([]byte(types.Html404))
		return
	}

	resourceETag := resp.Header.Get("ETag")
	etag := c.Ctx.Request.Header.Get("If-None-Match")
	//logs.Debug("If-None-Match:", etag)
	if etag == resourceETag {
		c.Ctx.Output.Status = 304
		return
	}

	var buf = new(bytes.Buffer)
	err = tc.DownloadPrivate2Stream(s3key, buf)
	if err != nil {
		logs.Warning("[Resource] can not download, ip: %s, router: %s, s3key: %s, err: %v",
			c.IP, c.Ctx.Input.URL(), s3key, err)
		c.Ctx.Output.Status = 404
		_ = c.Ctx.Output.Body([]byte(types.Html404))
		return
	}

	c.Ctx.Output.Header("Content-Type", resp.Header.Get("Content-Type"))
	if resourceETag != "" {
		c.Ctx.Output.Header("Etag", resp.Header.Get("ETag"))
	}
	_ = c.Ctx.Output.Body(buf.Bytes())
}
