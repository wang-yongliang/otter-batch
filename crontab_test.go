package schedule

import (
	"fmt"
	"github.com/lishimeng/go-log"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

type testJob struct {
	name string
}
// Run 任务执行
func (j *testJob) run() {
	log.Info("测试任务0000", time.Now())
}

func (j *testJob) Run() {
	j.run()
}
func TestCrontab(t *testing.T) {
	crontab := cron.New(cron.WithSeconds()) //精确到秒
	//定义定时器调用的任务函数
	task1 := func() {
		fmt.Println("hello world1111", time.Now())
	}
	task2 := func() {
		fmt.Println("hello world2222", time.Now())
	}
	task3 := func() {
		fmt.Println("hello world3333", time.Now())
	}
	//定时任务
	spec1 := "*/20 * * * * ?" //cron表达式，每五秒一次
	spec2 := "*/10 * 10 * * ?" //cron表达式，每10秒一次
	spec3 := "1 * * * * ?" //cron表达式，每分钟一次
	// 添加定时任务,
	crontab.AddFunc(spec1, task1)
	crontab.AddFunc(spec2, task2)

	tj := testJob{name: "ryker"}
	// 启动定时器
	crontab.Start()
	e3, _ := crontab.AddFunc(spec3, task3)
	e4, _ := crontab.AddJob(spec3, &tj)

	// 打印下次运行时间
	//for _, e := range crontab.Entries() {
	//	fmt.Printf("entry time: %d, %s\n", e.ID, e.Next.Format("2006-01-02 15:04:05"))
	//}
	fmt.Printf("entry time e3: %d, %s\n", e3, crontab.Entry(e3).Next.Format("2006-01-02 15:04:05"))
	fmt.Printf("entry time e4: %d, %s\n", e4, crontab.Entry(e4).Next.Format("2006-01-02 15:04:05"))
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	// 根据实际情况进行控制
	select {} //阻塞主线程停止
}
