package event

import "encoding/json"

type MDC struct {
	Tenant       string `json:"tenantId"`
	Subscription string `json:"subscriptionId"`
	Trace        string `json:"traceId"`
}

type Event struct {
	Timestamp string `json:"timestamp"`
	Logger    string `json:"loggerName"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	MDC       MDC    `json:"mdc"`
	Line      string
}

func ToJson(row string) (Event, bool) {

	var js Event
	if json.Unmarshal([]byte(row), &js) == nil {
		return js, true
	}

	return Event{}, false
}
