package analytics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/DlnKot/arc/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	mu       sync.Mutex
	session  *domain.AnalyticsSession
	clientID string
	dataFile string
	logger   Logger
}

type Logger interface {
	Infof(format string, args ...any)
}

func New(dataDir string, logger Logger) *Service {
	clientID := loadOrCreateClientID(dataDir)
	dataFile := filepath.Join(dataDir, "session.json")

	s := &Service{
		clientID: clientID,
		dataFile: dataFile,
		logger:   logger,
	}
	s.startSession()
	return s
}

func (s *Service) startSession() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.session = &domain.AnalyticsSession{
		SessionID:    uuid.New().String(),
		ClientID:     s.clientID,
		SessionStart: time.Now().UTC(),
		Events:       []domain.AnalyticsEvent{},
		Stats: domain.AnalyticsStats{
			LaunchesByType:     map[string]int{"rdp": 0, "horizon": 0, "citrix": 0, "vpn": 0},
			TabsViewCount:      map[string]int{"connections": 0, "settings": 0, "network": 0, "help": 0},
			NetworkCheckCount:  0,
			HelpViewsBySection: map[string]int{"helpdesk": 0, "indeed": 0, "remote-types": 0, "chatbot": 0},
			TotalLaunches:      0,
			Errors:             []string{},
		},
	}
	s.logInfo("analytics session started: %s", s.session.SessionID)
}

func (s *Service) EndSession() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.session == nil {
		return
	}

	s.session.SessionEnd = time.Now().UTC()
	s.session.SessionDuration = s.session.SessionEnd.Unix() - s.session.SessionStart.Unix()
	s.logInfo("analytics session ended: %s, duration: %ds", s.session.SessionID, s.session.SessionDuration)

	s.saveSession()
	s.session = nil
}

func (s *Service) TrackEvent(eventType string, data map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.session == nil {
		return
	}

	event := domain.AnalyticsEvent{
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	s.session.Events = append(s.session.Events, event)
}

func (s *Service) TrackAppStart() {
	s.TrackEvent("app_start", nil)
}

func (s *Service) TrackTabView(tab string) {
	s.TrackEvent("tab_view", map[string]any{"tab": tab})
	if s.session != nil {
		s.session.Stats.TabsViewCount[tab]++
	}
}

func (s *Service) TrackHelpView(section string) {
	s.TrackEvent("help_view", map[string]any{"section": section})
	if s.session != nil {
		s.session.Stats.HelpViewsBySection[section]++
	}
}

func (s *Service) TrackNetworkCheck() {
	s.TrackEvent("network_check", nil)
	if s.session != nil {
		s.session.Stats.NetworkCheckCount++
	}
}

func (s *Service) TrackConnectionLaunch(connType string) {
	s.TrackEvent("connection_launch", map[string]any{"type": connType})
	if s.session != nil {
		s.session.Stats.LaunchesByType[connType]++
		s.session.Stats.TotalLaunches++
	}
}

func (s *Service) TrackError(errMsg string) {
	s.TrackEvent("error", map[string]any{"message": errMsg})
	if s.session != nil {
		s.session.Stats.Errors = append(s.session.Stats.Errors, errMsg)
	}
}

func (s *Service) saveSession() {
	if s.session == nil {
		return
	}

	var existingData domain.AnalyticsData
	if data, err := os.ReadFile(s.dataFile); err == nil {
		_ = json.Unmarshal(data, &existingData)
	}

	existingData.Sessions = append(existingData.Sessions, *s.session)

	data, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		s.logInfo("analytics: failed to marshal session: %v", err)
		return
	}

	if err := os.WriteFile(s.dataFile, data, 0644); err != nil {
		s.logInfo("analytics: failed to save session: %v", err)
	}
}

func (s *Service) logInfo(format string, args ...any) {
	if s.logger != nil {
		s.logger.Infof(format, args...)
	}
}

func loadOrCreateClientID(dataDir string) string {
	clientIDFile := filepath.Join(dataDir, "client_id")
	if data, err := os.ReadFile(clientIDFile); err == nil {
		id := string(data)
		if len(id) == 36 {
			return id
		}
	}

	clientID := uuid.New().String()
	_ = os.WriteFile(clientIDFile, []byte(clientID), 0644)
	return clientID
}
