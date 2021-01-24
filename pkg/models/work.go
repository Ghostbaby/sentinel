package models

import "github.com/go-logr/logr"

// WorkResult 探测结果聚合
type WorkResult struct {
	Name string  `json:"name"` // 探测网站名称
	Host string  `json:"host"` // 被探测节点ip
	Max  float64 `json:"max"`  // 最高延迟
	Min  float64 `json:"min"`  // 最低延迟
	Avg  float64 `json:"avg"`  // 平均延迟
}

const (
	PingCnUrl               = "https://www.ping.cn"
	PingCnCheckPath         = "/check"
	PingCnPayloadType       = "ping"
	PingCnPayloadCreateTask = 1
	PingCnPayloadExecTask   = 0
)

type Job struct {
	Log        logr.Logger
	IP         string
	RetryTimes int
}

type Result struct {
	IP   string
	Work *WorkResult
}
