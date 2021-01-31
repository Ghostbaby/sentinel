package models

import "github.com/go-logr/logr"

// WorkResult 探测结果聚合
type WorkResult struct {
	Provider string           `json:"provider"` // 探测网站名称
	Host     string           `json:"host"`     // 被探测节点ip
	Name     string           `json:"name"`     // 线路名称
	Max      float64          `json:"max"`      // 最高延迟
	Min      float64          `json:"min"`      // 最低延迟
	Avg      float64          `json:"avg"`      // 平均延迟
	Loss     float64          `json:"lost"`     // 丢包率
	Details  []*ResultDetails `json:"details"`  // 详细测速信息
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

type ScopeResult struct {
	IP   string
	Work *WorkResult
}

type ResultDetails struct {
	Name    string  `json:"name"`
	Loss    float64 `json:"loss"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Avg     float64 `json:"avg"`
	Area    string  `json:"area"`
	IspName string  `json:"isp_name"`
}
