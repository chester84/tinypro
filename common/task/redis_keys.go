package task

import "tinypro/common/types"

const (
	// 异步事件处理
	eventLock  = "tinypro:lock:event"
	eventQueue = types.EventTaskRdsKey

	templateLock  = "tinypro:lock:template"
	templateQueue = "tinypro:queue:template"

	autoCloseLock  = `tinypro:lock:auto-close`
	autoCloseQueue = `tinypro:queue:auto-close`

	remindAttendClassLock  = "tinypro:lock:remind-attend-class"
	remindAttendClassQueue = "tinypro:queue:remind-attend-class"
)
