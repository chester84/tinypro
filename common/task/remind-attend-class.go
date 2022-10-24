package task

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/chester84/libtools"
	"github.com/gomodule/redigo/redis"
	"tinypro/common/lib/redis/storage"
	"tinypro/common/models"
	"tinypro/common/pkg/course"
	"tinypro/common/pkg/weixin"
	"tinypro/common/pkg/weixin/subscribe_template"
	"tinypro/common/types"
	"sync"
	"time"
)

type RemindAttendClass struct {
}

func init() {
	Register("remind-attend-class", "上课提醒", &RemindAttendClass{})
}

func (r *RemindAttendClass) Cancel() {
	logs.Informational("remind-attend-class cancel task")

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	lockKey := remindAttendClassLock
	_, _ = storageClient.Do("DEL", lockKey)
}

func (r *RemindAttendClass) Start() {
	logs.Info("remind-attend-class start task")

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	// +1 分布式锁 {{{
	lockKey := remindAttendClassLock
	lock, err := storageClient.Do("SET", lockKey, libtools.GetUnixMillis(), "NX")
	if err != nil || lock == nil {
		logs.Error("RemindAttendClass process is working, so, I will exit, err: %v", err)
		// ***! // 很重要!
		close(done)
		return
	}

	queueName := remindAttendClassQueue
	for {
		//10秒执行一次
		time.Sleep(time.Second * 8)

		if cancelled() {
			logs.Info("RemindAttendClass receive exit cmd.")
			break
		}

		TaskHeartBeat(storageClient, lockKey)

		// 1. 创建任务队列
		logs.Info("RemindAttendClass produceQueue")
		queueLen, err := redis.Int64(storageClient.Do("LLEN", queueName))
		logs.Debug("queueLen:", queueLen, ", err:", err)

		if err != nil {
			logs.Error("RemindAttendClass redis get exception `LLEN %s` error. err: %v", queueName, err)
			break
		} else if queueLen == 0 {
			// 队列是空
			var list []models.Course
			list, err = course.GetRemindCourseList()
			if err != nil {
				logs.Error("RemindAttendClass GetRemindCourseList err: %v", err)
				break
			}
			for _, item := range list {
				_, err = storageClient.Do("LPUSH", queueName, item.ID)
				if err != nil {
					logs.Error("RemindAttendClass redis> LPUSH %s %d, err: %v", queueName, item.ID, err)
				}
			}
		}

		var wg sync.WaitGroup
		// 可视情况加工作 goroutine 数
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go consumeRemindAttendClassQueue(&wg, i)
		}

		// 主 goroutine,等待工作 goroutine 正常结束
		wg.Wait()
	}

	// -1 正常退出时,释放锁 }}}
	_, err = storageClient.Do("DEL", lockKey)
	if err != nil {
		logs.Error("RemindAttendClass redis> DEL %s, err: %v", lockKey, err)
	}

	logs.Info("RemindAttendClass politeness exit.")
}

func consumeRemindAttendClassQueue(wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	logs.Info("RemindAttendClass will do consumeRemindAttendClassQueue, workerID:", workerID)

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	queueName := remindAttendClassQueue

	for {
		if cancelled() {
			logs.Info("RemindAttendClass receive exit cmd, workID:", workerID)
			break
		}

		qValue, err := redis.String(storageClient.Do("RPOP", queueName))
		if err != nil && err.Error() != redis.ErrNil.Error() {
			logs.Error("RemindAttendClass RPOP error workID:%d, err:%v", workerID, err)
			break
		}

		// 没有可供消费的数据
		if qValue == "" {
			logs.Info("RemindAttendClass no data for consume, so exit, workID:", workerID)
			break
		}

		cmdID, _ := libtools.Str2Int64(qValue)
		if cmdID == types.TaskExitCmd {
			logs.Info("RemindAttendClass receive exit cmd, I will exit after jobs done. workID:", workerID, ", cmdID:", cmdID)
			// ***! // 很重要!
			close(done)
			break
		}

		// 真正开始工作了
		addCurrentData(qValue, cmdID)
		handleRemindAttendClassMessage(cmdID, workerID)
		removeCurrentData(qValue)
	}
}

func handleRemindAttendClassMessage(dataID int64, workerID int) {
	logs.Debug(`handleRemindAttendClassMessage enters, dataID %d, err: %d`, dataID, workerID)

	if dataID <= 0 {
		logs.Warning("RemindAttendClass get zero id, workerID: %d", workerID)
		return
	}

	var courseObj models.Course
	err := models.OrmOneByPkId(dataID, &courseObj)
	if err != nil {
		logs.Error(`handleRemindAttendClassMessage OrmOneByPkId course unexpect data, dataId: %d, err: %v`, dataID, err)
		return
	}

	data := make(map[string]weixin.SubscribeMessageSendDataItem)
	data, err = subscribe_template.BuildRemindAttendClassMsg(courseObj)
	if err != nil {
		logs.Error(`handleRemindAttendClassMessage BuildRemindAttendClassMsg err, dataId: %d, err: %v`, dataID, err)
		return
	}

	logs.Debug(`handleRemindAttendClassMessage SubscribeMessageSendDataItem, data:%v`, data)

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	key := subscribe_template.GetRDSRemindCourseKey(libtools.Int642Str(courseObj.ID))
	cursor := 0

	logs.Debug(`handleRemindAttendClassMessage GetRDSRemindCourseKey, key:%s`, key)

	for {
		arr := make([]string, 0)
		var reply []interface{}
		// 每次尝试取100条，这个100条，只是我们给予的建议，redis会根据实际情况返回数据量的
		reply, err = redis.Values(storageClient.Do("SSCAN", key, cursor, "COUNT", 100))
		if err != nil && err != redis.ErrNil {
			logs.Error(`handleRemindAttendClassMessage redis.Values err, key %s, err: %v`, key, err)
			return
		}
		if _, err = redis.Scan(reply, &cursor, &arr); err != nil && err != redis.ErrNil {
			logs.Error(`handleRemindAttendClassMessage redis.Scan err, key %s, err: %v`, key, err)
			return
		}

		logs.Debug("handleRemindAttendClassMessage cursor %d", cursor)
		logs.Debug("handleRemindAttendClassMessage arr %#v", arr)

		for _, userID := range arr {
			userId, _ := libtools.Str2Int64(userID)

			var userObj models.AppUser
			err = models.OrmOneByPkId(userId, &userObj)
			if err != nil {
				logs.Error(`handleRemindAttendClassMessage OrmOneByPkId appuser err, userId %d, err: %v`, userId, err)
				return
			}

			logs.Debug("handleRemindAttendClassMessage userObj %v", userObj)

			//logs.Debug(data)
			page := fmt.Sprintf("/pages/detail/courseclass?sn=%d", courseObj.ID)
			err = weixin.SendSubscribeMessage(weixin.AppSNWxMng, userObj, courseObj, weixin.SubAttendClass, page, data)
			if err != nil {
				logs.Warning(`handleRemindAttendClassMessage SendSubscribeMessage err: %v`, err)
			}

			// 遍历伴随着srem会有问题...
			// 直接在最外层del掉key值即可
			//_, err = storageClient.Do("SREM", key, userID)
			//if err != nil {
			//	logs.Error(`handleRemindAttendClassMessage SendSubscribeMessage err: %v`, err)
			//}
		}

		// 再次为0时，集合中的数据全部取出了
		if cursor == 0 {
			logs.Debug(`handleRemindAttendClassMessage GetRDSRemindCourseKey, cursor==0 key:%s`, key)
			_, err = storageClient.Do("DEL", key)
			if err != nil {
				logs.Error(`handleRemindAttendClassMessage SendSubscribeMessage DEL, key:%s, err: %v`, key, err)
			}

			logs.Debug(`handleRemindAttendClassMessage courseObj.Remind, courseObj:%v`, courseObj)
			courseObj.RemindFlag = 1
			_, err = models.OrmUpdate(&courseObj, []string{"remind_flag"})
			if err != nil {
				logs.Error(`handleRemindAttendClassMessage OrmUpdate course unexpect data, dataId: %d, err: %v`, dataID, err)
				return
			}

			logs.Debug(`handleRemindAttendClassMessage courseObj after update, courseObj:%v`, courseObj)

			break
		}
	}

	logs.Info("RemindAttendClass finish handle, workerID: %d", workerID)
}
