package schedule

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/lishimeng/go-log"
	"github.com/robfig/cron/v3"
)

type Schedule interface {
	AddJob(executorHandler string, spec string, f TaskFunc) (err error)
	RunJob(jobId string) CallElement
}
type schedule struct {
	cron           *cron.Cron
	mu             sync.RWMutex
	regList        []RegTask
	jobList        jobList
	runningJobList taskList
}

type RegTask struct {
	Handler  string
	CronSpec string
	Func     TaskFunc
}

func (s *schedule) AddJob(executorHandler string, spec string, f TaskFunc) (err error) {
	return s.addJob(executorHandler, spec, f)
}

// RunJob 运行任务
func (e *schedule) RunJob(jobId string) CallElement {
	return e.runJob(jobId)
}

func (s *schedule) run() (err error) {
	if len(s.regList) == 0 {
		Server = s
		return errors.New("no task registered")
	}
	for _, task := range s.regList {
		err = s.AddJob(task.Handler, task.CronSpec, task.Func)
		if err != nil {
			log.Info(err)
			return
		}
	}
	return
}

func (s *schedule) addJob(executorHandler string, spec string, f TaskFunc) (err error) {
	if len(executorHandler) == 0 || len(spec) == 0 || f == nil {
		return errors.New("invalid param")
	}
	jobId := uuid.New().String()
	job := job{
		jobId:    jobId,
		spec:     spec,
		handler:  executorHandler,
		function: f,
	}
	entryID, err := s.cron.AddJob(spec, &job)
	if err != nil {
		log.Info(err)
		log.Info("任务[ %s ]定时格式错误，添加失败。handler: %s", jobId, executorHandler)
		return
	}
	entry := s.cron.Entry(entryID)
	job.entryID = entryID
	job.nextTime = entry.Next
	s.jobList.Set(jobId, &job)
	log.Info("任务初始化[成功]jobId:%s, entryID: %d, handler: %s, 下次执行时间：%s", jobId, entryID, executorHandler, job.nextTime.Format("2006-01-02 15:04:05"))
	return
}

// 运行一个任务
func (s *schedule) runJob(jobId string) CallElement {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.jobList.Exists(jobId) {
		return Callback(Error, fmt.Sprintf("任务[ %s ]没有注册", jobId))
	}
	job := s.jobList.Get(jobId)
	task := Task{
		Id:   jobId,
		Name: job.handler,
		fn:   job.function,
	}
	task.Id = job.jobId
	task.Name = job.handler
	task.fn = job.function

	s.runningJobList.Set(task.Id, &task)
	log.Info("任务[ %s ] %s 开始执行:", task.Id, task.Name)
	res := task.Run()
	// 更新任务状态
	s.updateJob(jobId)
	return res
}

func (s *schedule) updateJob(jobId string) {
	job := s.jobList.Get(jobId)
	entry := s.cron.Entry(job.entryID)
	job.prevTime = job.nextTime
	job.nextTime = entry.Next
	s.jobList.Set(jobId, job)
	log.Info("任务执行结束：jobId:%s, entryID: %d, handler: %s, 下次执行时间：%s", jobId, job.entryID, job.function, job.nextTime.Format("2006-01-02 15:04:05"))
	// 删除执行中的任务
	if s.runningJobList.Exists(jobId) {
		s.runningJobList.Del(jobId)
	}
}

func Callback(code int, msg string) CallElement {
	return CallElement{
		HandleCode: code,
		HandleMsg:  msg,
	}
}
