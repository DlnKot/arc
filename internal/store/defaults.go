package store

import "github.com/DlnKot/arc/internal/domain"

func defaultStoreData(settings domain.Settings) domain.StoreData {
	return domain.StoreData{
		Settings:                   settings,
		ConnectionsUser:            []map[string]any{},
		DefaultConnectionOverrides: map[string]map[string]any{},
	}
}
