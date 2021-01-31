package models

type ScpResult struct {
	IP      string `json:"ip"` // 被探测节点ip
	Name    string `json:"name"`
	IsReady bool   `json:"is_ready"`
}
