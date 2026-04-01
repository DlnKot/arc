package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/DlnKot/arc/internal/domain"
	"github.com/DlnKot/arc/internal/logging"
)

type Service struct {
	mu        sync.Mutex
	appName   string
	storePath string
	data      domain.StoreData
	defaults  deploymentDefaults
	logger    logging.Logger
}

func New(appName string, logger logging.Logger) *Service {
	defaults := loadEmbeddedDeploymentDefaults()
	return &Service{appName: appName, data: domain.StoreData{Settings: defaults.Settings, ConnectionsUser: []map[string]any{}, DefaultConnectionOverrides: map[string]map[string]any{}}, defaults: defaults, logger: logger}
}

func (s *Service) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	storeDir := filepath.Join(configDir, s.appName)
	if err := os.MkdirAll(storeDir, 0o755); err != nil {
		return err
	}

	s.storePath = filepath.Join(storeDir, "config.json")
	content, err := os.ReadFile(s.storePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.data = defaultStoreData(s.defaults.Settings)
			if s.logger != nil {
				s.logger.Infof("store not found, creating default config at %s", s.storePath)
			}
			return s.saveLocked()
		}
		return err
	}

	if len(content) > 1<<20 {
		return errors.New("config.json is larger than 1MB")
	}

	var raw domain.StoreData
	if err := json.Unmarshal(content, &raw); err != nil {
		return err
	}

	users := normalizeConnections(raw.ConnectionsUser)
	overrides := cloneOverrides(raw.DefaultConnectionOverrides)

	merged := domain.SettingsFromMap(s.defaults.Settings.ToMap())
	if raw.Settings.User.Domain != "" || raw.Settings.User.Username != "" {
		merged.User = raw.Settings.User
	}
	if raw.Settings.Rdp.Resolution != "" {
		merged.Rdp = raw.Settings.Rdp
	}
	if raw.Settings.Horizon.AppName != "" || raw.Settings.Horizon.DesktopProtocol != "" {
		merged.Horizon = raw.Settings.Horizon
	}
	if raw.Settings.Citrix.AccountName != "" || raw.Settings.Citrix.ResourceName != "" {
		merged.Citrix = raw.Settings.Citrix
	}
	if raw.Settings.General.MinimizeToTray || raw.Settings.General.StartMinimized {
		merged.General = raw.Settings.General
	}
	if raw.Settings.NetworkCheck.LatencyThresholdMs > 0 {
		merged.NetworkCheck = raw.Settings.NetworkCheck
	}
	if raw.Settings.Updates.InternalServerURL != "" || raw.Settings.Updates.AutoCheck || raw.Settings.Updates.InstallOnQuit || raw.Settings.Updates.UseGithub {
		merged.Updates = raw.Settings.Updates
	}

	s.data = domain.StoreData{
		Settings:                   merged,
		ConnectionsUser:            users,
		DefaultConnectionOverrides: overrides,
	}
	if s.logger != nil {
		s.logger.Infof("store loaded from %s", s.storePath)
	}

	return nil
}

func (s *Service) GetConnections() []map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return composeConnections(s.defaults.Connections, s.data.DefaultConnectionOverrides, s.data.ConnectionsUser)
}

func (s *Service) SaveConnection(connection map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalized := normalizeConnection(connection)
	if factoryID := strings.TrimSpace(asString(normalized["factoryId"])); factoryID != "" {
		for _, template := range s.defaults.Connections {
			if asString(template["factoryId"]) != factoryID {
				continue
			}
			if s.data.DefaultConnectionOverrides == nil {
				s.data.DefaultConnectionOverrides = map[string]map[string]any{}
			}
			name := strings.TrimSpace(asString(normalized["name"]))
			if name == "" || name == asString(template["name"]) {
				delete(s.data.DefaultConnectionOverrides, factoryID)
			} else {
				s.data.DefaultConnectionOverrides[factoryID] = map[string]any{"name": name}
			}
			return s.saveLocked()
		}
	}

	normalized["isDefault"] = false
	normalized["factoryId"] = ""
	normalized["isUserModified"] = true
	id := asString(normalized["id"])
	if id == "" {
		normalized["id"] = newConnectionID()
		id = asString(normalized["id"])
	}

	updated := false
	for index, existing := range s.data.ConnectionsUser {
		if asString(existing["id"]) == id {
			s.data.ConnectionsUser[index] = normalized
			updated = true
			break
		}
	}

	if !updated {
		s.data.ConnectionsUser = append(s.data.ConnectionsUser, normalized)
	}

	return s.saveLocked()
}

func (s *Service) DeleteConnection(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := make([]map[string]any, 0, len(s.data.ConnectionsUser))
	for _, connection := range s.data.ConnectionsUser {
		if asString(connection["id"]) != id {
			filtered = append(filtered, connection)
		}
	}
	s.data.ConnectionsUser = filtered
	return s.saveLocked()
}

func (s *Service) GetSettings() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data.Settings.ToMap()
}

func (s *Service) SaveSettings(settings map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Settings = domain.SettingsFromMap(mergeMaps(s.defaults.Settings.ToMap(), settings))
	return s.saveLocked()
}

func (s *Service) ResetDefaultConnections() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.DefaultConnectionOverrides = map[string]map[string]any{}
	return s.saveLocked()
}

func (s *Service) saveLocked() error {
	payload, err := json.MarshalIndent(domain.StoreData{
		Settings:                   s.data.Settings,
		ConnectionsUser:            normalizeConnections(s.data.ConnectionsUser),
		DefaultConnectionOverrides: cloneOverrides(s.data.DefaultConnectionOverrides),
	}, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := s.storePath + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		_ = os.Remove(s.storePath)
	}
	if err := os.Rename(tmpPath, s.storePath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if s.logger != nil {
		s.logger.Infof("store saved to %s", s.storePath)
	}
	return nil
}
