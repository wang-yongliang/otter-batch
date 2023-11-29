package schedule

const baseApi = "/api"

const (
	ActionRun      = baseApi + "/run"  // 启动任务
	ActionKill     = baseApi + "/kill" // 终止任务
	ActionLog      = baseApi + "/log"  // 任务日志
	ActionBeat     = baseApi + "/log"  // 心跳检测
	ActionIdleBeat = baseApi + "/log"  // 忙碌检测
)
