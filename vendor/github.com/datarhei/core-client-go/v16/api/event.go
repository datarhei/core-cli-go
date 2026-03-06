package api

type LogEvent struct {
	Timestamp int64  `json:"ts" format:"int64"`
	Level     string `json:"level"`
	Component string `json:"event"`
	Message   string `json:"message"`
	Caller    string `json:"caller"`
	CoreID    string `json:"core_id"`

	Data map[string]string `json:"data"`
}

type LogEventFilter struct {
	Component string            `json:"event"`
	Message   string            `json:"message"`
	Level     string            `json:"level"`
	Caller    string            `json:"caller"`
	CoreID    string            `json:"core_id"`
	Data      map[string]string `json:"data"`
}

type LogEventFilters struct {
	Filters []LogEventFilter `json:"filters"`
}

type MediaEvent struct {
	Action    string   `json:"action"`
	Name      string   `json:"name,omitempty"`
	Names     []string `json:"names,omitempty"`
	Timestamp int64    `json:"ts"`
}

type ProcessEvent struct {
	ProcessID string `json:"pid"`
	Domain    string `json:"domain"`
	Type      string `json:"type"`
	Line      string `json:"line"`
	Progress  any    `json:"progress"`
	Timestamp int64  `json:"ts"`
	CoreID    string `json:"core_id"`
}

type ProcessEventFilter struct {
	ProcessID string `json:"pid"`
	Domain    string `json:"domain"`
	Type      string `json:"type"`
	CoreID    string `json:"core_id"`
}

type ProcessEventFilters struct {
	Filters []ProcessEventFilter `json:"filters"`
}
