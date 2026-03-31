package domain

type StoreData struct {
	Settings                   Settings                  `json:"settings"`
	ConnectionsUser            []map[string]any          `json:"connectionsUser"`
	DefaultConnectionOverrides map[string]map[string]any `json:"defaultConnectionOverrides,omitempty"`
}

type PingMetrics struct {
	LossPercent float64 `json:"lossPercent"`
	AvgMs       float64 `json:"avgMs"`
	MinMs       float64 `json:"minMs"`
	MaxMs       float64 `json:"maxMs"`
	Raw         string  `json:"raw,omitempty"`
	Error       string  `json:"error,omitempty"`
}

type PingEvaluation struct {
	Status         string `json:"status"`
	Label          string `json:"label"`
	Recommendation string `json:"recommendation,omitempty"`
}
