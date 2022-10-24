package task

import (
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"tinypro/common/lib/redis/storage"
	"tinypro/common/models"
	"tinypro/common/pkg/payment"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type AutoClose struct {
}

func init() {
	Register("auto-close", "自动关闭支付订单", &AutoClose{})
}

func (r *AutoClose) Cancel() {
	logs.Informational("[task->AutoClose] cancel task")

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	lockKey := autoCloseLock
	_, _ = storageClient.Do("DEL", lockKey)
}

func (r *AutoClose) Start() {
	logs.Info("[task->AutoClose] start task")

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	// +1 分布式锁 {{{
	lockKey := autoCloseLock
	lock, err := storageClient.Do("SET", lockKey, libtools.GetUnixMillis(), "NX")
	if err != nil || lock == nil {
		logs.Error("[task->AutoClose] process is working, so, I will exit, err: %v", err)
		// ***! // 很重要!
		close(done)
		return
	}

	queueName := autoCloseQueue
	for {
		if cancelled() {
			logs.Info("[task->AutoClose] receive exit cmd.")
			break
		}

		TaskHeartBeat(storageClient, lockKey)

		// 1. 创建任务队列
		logs.Info("[task->AutoClose] produceTemplateQueue")
		queueLen, err := redis.Int64(storageClient.Do("LLEN", queueName))
		logs.Debug("queueLen:", queueLen, ", err:", err)

		if err != nil {
			logs.Error("[task->AutoClose] redis get exception `LLEN %s` error. err: %v", queueName, err)
			break
		} else if queueLen == 0 {
			// 队列是空
			box := payment.EduFetchNeedCloseOrder()
			for _, payId := range box {
				_, err = storageClient.Do("LPUSH", queueName, payId)
				if err != nil {
					logs.Error("[task->AutoClose] redis> LPUSH %s %d, err: %v", queueName, payId, err)
				}
			}

			if len(box) <= 0 {
				time.Sleep(time.Second)
			}
		}

		var wg sync.WaitGroup
		// 可视情况加工作 goroutine 数
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go consumeAutoCloseQueue(&wg, i)
		}

		// 主 goroutine,等待工作 goroutine 正常结束
		wg.Wait()
	}

	// -1 正常退出时,释放锁 }}}
	_, err = storageClient.Do("DEL", lockKey)
	if err != nil {
		logs.Error("[task->AutoClose] redis> DEL %s, err: %v", lockKey, err)
	}

	logs.Info("[task->AutoClose] politeness exit.")
}

func consumeAutoCloseQueue(wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	logs.Info("It will do consumeAutoCloseQueue, workerID:", workerID)

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	queueName := autoCloseQueue

	for {
		if cancelled() {
			logs.Info("[consumeAutoCloseQueue] receive exit cmd, workID:", workerID)
			break
		}

		qValue, err := redis.String(storageClient.Do("RPOP", queueName))
		if err != nil {
			logs.Error("[consumeAutoCloseQueue] RPOP error workID:%d, err:%v", workerID, err)
		}

		// 没有可供消费的数据
		if qValue == "" {
			logs.Info("[consumeAutoCloseQueue] no data for consume, so exit, workID:", workerID)
			break
		}

		cmdID, _ := libtools.Str2Int64(qValue)
		if cmdID == types.TaskExitCmd {
			logs.Info("[consumeAutoCloseQueue] receive exit cmd, I will exit after jobs done. workID:", workerID, ", cmdID:", cmdID)
			// ***! // 很重要!
			close(done)
			break
		}

		// 真正开始工作了
		addCurrentData(qValue, cmdID)
		handleAutoClose(cmdID, workerID)
		removeCurrentData(qValue)
	}
}

func handleAutoClose(dataId int64, workerId int) {
	if dataId <= 0 {
		logs.Warning("[handleAutoClose] get zero id, workerId: %d", workerId)
		return
	}

	var payObj models.EduPayment
	err := models.OrmOneByPkId(dataId, &payObj)
	if err != nil {
		logs.Error(`[handleAutoClose] get unexpect data, dataId: %d, err: %v`, dataId, err)
		return
	}

	if payObj.Status != types.PaymentStatusCreated || payObj.ClosedAt != 0 {
		logs.Error(`[handleAutoClose] data status unexpcet, data: %#v`, payObj)
		return
	}

	payObj.Status = types.PaymentStatusClosed
	payObj.ClosedAt = libtools.GetUnixMillis()
	_, err = models.OrmUpdate(&payObj, []string{`Status`, `ClosedAt`})
	if err != nil {
		logs.Error(`[handleAutoClose] db update pay data exception, data: %#v, err: %v`, payObj, err)
	}

	logs.Info("[handleAutoClose] finish handel, workerId: %d, dataId: %d", workerId, dataId)
}
