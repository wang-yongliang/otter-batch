package schedule

import (
	"time"

	"github.com/lishimeng/go-log"
	"github.com/robfig/cron/v3"
)

type Job interface {
	Run() // 任务执行
}
type job struct {
	jobId    string
	spec     string
	entryID  cron.EntryID
	handler  string
	function TaskFunc
	prevTime time.Time
	nextTime time.Time
}

// Run 任务执行
func (j *job) run() {
	callElement := Server.RunJob(j.jobId)
	log.Info("任务【 %s  】 %s 执行结果： code=%d, msg=%s", j.jobId, j.handler, callElement.HandleCode, callElement.HandleMsg)
}

func (j *job) Run() {
	j.run()
}
