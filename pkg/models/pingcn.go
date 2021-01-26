package models

type PingCnPayload struct {
	Host       string `json:"host"`        // 需要探测的主机
	Type       string `json:"type"`        // ping
	CreateTask int    `json:"create_task"` // 0代表执行 1代表创建任务
	TaskID     string `json:"task_id"`     // 创建任务之后返回task_id
}

// PingCnTask ping.cn 任务编号获取结果
type PingCnTask struct {
	Code int         `json:"code"` // 1代表执行成功
	Data *PingCnData `json:"data"` // 返回任务数据
}

type PingCnData struct {
	TaskID string `json:"taskID"` // 返回任务编号
}

// PingCnResult ping.cn 查询结果
type PingCnResult struct {
	Code int               `json:"code"` // 1代表执行成功
	Data *PingCnResultData `json:"data"` // 返回探测结果
}

type PingCnResultInfo struct {
	NodeName   string `json:"node_name"`
	Packets    int    `json:"packets"`
	Received   int    `json:"received"`
	PacketLoss int    `json:"packet_loss"`
}

type PingCnResultData struct {
	Cost     float64              `json:"cost"`     // 返回执行耗时
	InitData PingCnResultInitData `json:"initData"` // 聚合数据
}

type PingCnResultInitData struct {
	MinMaxAvg PingCnResultMinMaxAvg `json:"minMaxAvg"` // 返回最低、最高、平均延迟
	Result    []*PingCnResultInfo   `json:"result"`
}

type PingCnResultMinMaxAvg struct {
	Max PingCnResultCost `json:"max"` // 最高延迟
	Min PingCnResultCost `json:"min"` // 最低延迟
	Avg PingCnResultCost `json:"avg"` // 平均延迟
}

type PingCnResultCost struct {
	NodeName string  `json:"node_name"` // 区域名称，类似 西藏山南移动
	Cost     float64 `json:"cost"`      // 具体延迟
}
