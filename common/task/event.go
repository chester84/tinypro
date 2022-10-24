package task

import (
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/lib/redis/storage"
	"tinypro/common/task/event/runner"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type Event struct {
}

func init() {
	Register("event", "异步事件消费处理", &Event{})
}

func (r *Event) Cancel() {
	logs.Informational("[task->Event] cancel task")

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	lockKey := eventLock
	_, _ = storageClient.Do("DEL", lockKey)
}

func (r *Event) Start() {
	logs.Info("[task->Event] start task")

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	// +1 分布式锁 {{{
	lockKey := eventLock
	lock, err := storageClient.Do("SET", lockKey, libtools.GetUnixMillis(), "NX")
	if err != nil || lock == nil {
		logs.Error("[task->Event] process is working, so, I will exit, err: %v", err)
		// ***! // 很重要!
		close(done)
		return
	}

	for {
		if cancelled() {
			logs.Info("[task->Event] receive exit cmd.")
			break
		}

		go func() {
			for {
				storageClientHeart := storage.RedisStorageClient.Get()
				TaskHeartBeat(storageClientHeart, lockKey)
				storageClientHeart.Close()
				time.Sleep(5 * time.Minute)
			}
		}()

		// 消费队列
		logs.Info("[consumeEventQueue] consume queue")
		var wg sync.WaitGroup
		// 可视情况加工作 goroutine 数,一期只开2个
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go consumeEventQueue(&wg, i)
		}

		// 主 goroutine,等待工作 goroutine 正常结束
		wg.Wait()
	}

	// -1 正常退出时,释放锁 }}}
	_, _ = storageClient.Do("DEL", lockKey)

	logs.Info("[task->Event] politeness exit.")
}

func consumeEventQueue(wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	logs.Info("It will do consumeEventQueue, workerID:", workerID)

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	queueName := eventQueue
	for {
		if cancelled() {
			logs.Info("[consumeEventQueue] receive exit cmd, workID:", workerID)
			break
		}

		qValueByte, err := storageClient.Do("RPOP", queueName)
		if err != nil {
			logs.Error("[consumeEventQueue] RPOP error workID:%d, err:%v", workerID, err)
		}

		// 没有可供消费的数据
		if qValueByte == nil {
			logs.Info("[consumeEventQueue] no data for consume, I will sleep moments, workID:", workerID)
			time.Sleep(time.Second)
			continue
		}

		str := string(qValueByte.([]byte))
		queueValToCmd, _ := libtools.Str2Int64(str)
		if queueValToCmd == types.TaskExitCmd {
			logs.Info("[consumeEventQueue] receive exit cmd, I will exit after jobs done. workID:", workerID, ", queueVal:", queueValToCmd)
			// ***! // 很重要!
			close(done)
			break
		}

		// 真正开始工作了
		addCurrentData(str, libtools.GetUnixMillis())
		_, _ = runner.ParseAndRun(qValueByte.([]byte))
		removeCurrentData(str)
	}
}
