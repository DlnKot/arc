package store

import "strings"

func normalizeConnections(connections []map[string]any) []map[string]any {
	if len(connections) == 0 {
		return []map[string]any{}
	}

	result := make([]map[string]any, 0, len(connections))
	for _, connection := range connections {
		result = append(result, normalizeConnection(connection))
	}
	return result
}

func normalizeConnection(connection map[string]any) map[string]any {
	out := cloneMap(connection)
	if out == nil {
		out = map[string]any{}
	}

	out["id"] = strings.TrimSpace(asString(out["id"]))
	out["factoryId"] = strings.TrimSpace(asString(out["factoryId"]))
	out["type"] = strings.TrimSpace(asString(out["type"]))
	out["name"] = strings.TrimSpace(asString(out["name"]))
	out["host"] = strings.TrimSpace(asString(out["host"]))
	out["desktopPool"] = strings.TrimSpace(asString(out["desktopPool"]))
	out["storeUrl"] = strings.TrimSpace(asString(out["storeUrl"]))
	out["username"] = strings.TrimSpace(asString(out["username"]))
	out["description"] = strings.TrimSpace(asString(out["description"]))

	if _, ok := out["isUserModified"]; !ok {
		out["isUserModified"] = false
	}
	if _, ok := out["isDefault"]; !ok {
		out["isDefault"] = false
	}

	return out
}
