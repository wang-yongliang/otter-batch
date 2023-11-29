package schedule

import (
	"github.com/lishimeng/go-log"
	"github.com/robfig/cron/v3"
)

// Server 全局调度器
var Server Schedule

type Option func(s *schedule)

// WithTask 任务ID、处理函数、cron表达式
func WithTask(rt ...RegTask) Option {
	return func(s *schedule) {
		s.regList = append(s.regList, rt...)
	}
}

func Init(opts ...Option) (err error) {
	s := &schedule{}
	s.cron = cron.New(cron.WithSeconds())
	s.cron.Start()
	s.regList = make([]RegTask, 0)
	s.jobList = jobList{
		data: make(map[string]*job),
	}
	s.runningJobList = taskList{
		data: make(map[string]*Task),
	}
	for _, opt := range opts {
		opt(s)
	}
	err = s.run()
	if err != nil {
		Server = s
		return
	}
	log.Info("调度器初始化定时任务...共 %d 个", len(s.regList))
	Server = s
	return
}
