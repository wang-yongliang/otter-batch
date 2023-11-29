package schedule

import (
	"fmt"
	"runtime/debug"

	"github.com/lishimeng/go-log"
)

// TaskFunc 任务执行函数
type TaskFunc func() (code int, msg string)

// Task 任务
type Task struct {
	Id        string
	Name      string
	Param     *RunReq
	fn        TaskFunc
	StartTime int64
	EndTime   int64
}

// Run 运行任务
func (t *Task) Run() CallElement {
	defer func() CallElement {
		if err := recover(); err != nil {
			log.Info(err)
			debug.PrintStack() //堆栈跟踪
			return Callback(Error, fmt.Sprintf("panic: %v", err))
		}
		return Callback(Success, "")
	}()
	code, msg := t.fn()
	return Callback(code, msg)
}

// Info 任务信息
// func (t *Task) Info() string {
// 	return "任务ID[" + utils.Int64ToStr(t.Id) + "]任务名称[" + t.Name + "]"
// }
