package domain

import "time"

type StoreData struct {
	Settings                   Settings                  `json:"settings"`
	ConnectionsUser            []map[string]any          `json:"connectionsUser"`
	DefaultConnectionOverrides map[string]map[string]any `json:"defaultConnectionOverrides,omitempty"`
}

type PingMetrics struct {
	LostPercent float64 `json:"lossPercent"`
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

type AnalyticsEvent struct {
	Type      string         `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      map[string]any `json:"data,omitempty"`
}

type AnalyticsStats struct {
	LaunchesByType     map[string]int `json:"launchesByType"`
	TabsViewCount      map[string]int `json:"tabsViewCount"`
	NetworkCheckCount  int            `json:"networkCheckCount"`
	HelpViewsBySection map[string]int `json:"helpViewsBySection"`
	TotalLaunches      int            `json:"totalLaunches"`
	Errors             []string       `json:"errors"`
}

type AnalyticsSession struct {
	SessionID       string           `json:"sessionId"`
	ClientID        string           `json:"clientId"`
	SessionStart    time.Time        `json:"sessionStart"`
	Events          []AnalyticsEvent `json:"events"`
	Stats           AnalyticsStats   `json:"stats"`
	SessionEnd      time.Time        `json:"sessionEnd,omitempty"`
	SessionDuration int64            `json:"sessionDuration"`
}

type AnalyticsData struct {
	Sessions []AnalyticsSession `json:"sessions"`
}
