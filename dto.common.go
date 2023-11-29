package schedule

// CallElement 执行器执行完任务后，回调任务结果时使用
type CallElement struct {
	HandleCode int    `json:"handleCode"` //200表示正常,500表示失败
	HandleMsg  string `json:"handleMsg"`
}

// RunReq 触发任务请求参数
type RunReq struct {
	JobID           int64  `json:"jobId"`           // 任务ID
	ExecutorHandler string `json:"executorHandler"` // 任务标识
}

const (
	Success = 200
	Error   = 500
)
