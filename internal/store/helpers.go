package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func asString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case nil:
		return ""
	default:
		return fmt.Sprint(typed)
	}
}

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	payload, _ := json.Marshal(input)
	var output map[string]any
	_ = json.Unmarshal(payload, &output)
	if output == nil {
		return map[string]any{}
	}
	return output
}

func cloneConnections(input []map[string]any) []map[string]any {
	payload, _ := json.Marshal(input)
	var output []map[string]any
	_ = json.Unmarshal(payload, &output)
	if output == nil {
		return []map[string]any{}
	}
	return output
}

func cloneOverrides(input map[string]map[string]any) map[string]map[string]any {
	payload, _ := json.Marshal(input)
	var output map[string]map[string]any
	_ = json.Unmarshal(payload, &output)
	if output == nil {
		return map[string]map[string]any{}
	}
	return output
}

func mergeMaps(base map[string]any, override map[string]any) map[string]any {
	result := cloneMap(base)
	for key, value := range override {
		switch typed := value.(type) {
		case map[string]any:
			current, _ := result[key].(map[string]any)
			result[key] = mergeMaps(current, typed)
		default:
			result[key] = typed
		}
	}
	return result
}

func newConnectionID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.NewString())
}
