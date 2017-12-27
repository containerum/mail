package clickhouse

import "time"

type LogRecord struct {
	Method       string        `json:"method"`
	RequestTime  time.Time     `json:"request_time"`
	RequestSize  uint64        `json:"request_size"`
	ResponseSize uint64        `json:"response_size"`
	User         string        `json:"user"`
	Path         string        `json:"path"`
	Latency      time.Duration `json:"latency"`
	ID           string        `json:"id"`

	Status          uint   `json:"status"`
	Upstream        string `json:"upstream"`
	UserAgent       string `json:"user_agent"`
	Fingerprint     string `json:"fingerprint"`
	RequestHeaders  []int8 `json:"request_headers"`
	RequestBody     []int8 `json:"request_body"`
	ResponseHeaders []int8 `json:"response_headers"`
	ResponseBody    []int8 `json:"response_body"`
	GatewayID       string `json:"gateway_id"`
}
