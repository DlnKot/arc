package store

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DlnKot/arc/internal/domain"
)

//go:embed deployment-defaults.json
var embeddedDeploymentDefaults []byte

type deploymentDefaults struct {
	Settings    domain.Settings  `json:"settings"`
	Connections []map[string]any `json:"connections"`
}

func loadEmbeddedDeploymentDefaults() deploymentDefaults {
	var parsed deploymentDefaults
	if err := json.Unmarshal(embeddedDeploymentDefaults, &parsed); err != nil {
		return deploymentDefaults{Settings: domain.DefaultSettings()}
	}
	parsed.Connections = normalizeFactoryConnections(parsed.Connections)
	return parsed
}

func normalizeFactoryConnections(connections []map[string]any) []map[string]any {
	out := make([]map[string]any, 0, len(connections))
	for index, connection := range connections {
		normalized := normalizeConnection(connection)
		typeValue := strings.ToLower(asString(normalized["type"]))
		if typeValue == "" {
			continue
		}
		factoryID := strings.TrimSpace(asString(normalized["factoryId"]))
		if factoryID == "" {
			factoryID = fmt.Sprintf("factory-%s-%d", typeValue, index)
		}
		normalized["id"] = "factory:" + factoryID
		normalized["factoryId"] = factoryID
		normalized["type"] = typeValue
		normalized["isDefault"] = true
		normalized["isUserModified"] = false
		if defaultSettings, ok := normalized["defaultSettings"].(map[string]any); ok {
			normalized["clientSettings"] = cloneMap(defaultSettings)
		}
		out = append(out, normalized)
	}
	return out
}

func composeConnections(templates []map[string]any, overrides map[string]map[string]any, users []map[string]any) []map[string]any {
	result := make([]map[string]any, 0, len(templates)+len(users))
	for _, template := range templates {
		item := cloneMap(template)
		if override, ok := overrides[asString(template["factoryId"])]; ok {
			if name := strings.TrimSpace(asString(override["name"])); name != "" {
				item["name"] = name
			}
		}
		result = append(result, item)
	}
	return append(result, cloneConnections(users)...)
}
