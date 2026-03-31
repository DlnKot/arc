package launchers

import "github.com/DlnKot/arc/internal/logging"

type Service struct {
	logger logging.Logger
}

func New(logger logging.Logger) *Service {
	setLogger(logger)
	return &Service{logger: logger}
}

func (s *Service) LaunchRdp(connection map[string]any, settings map[string]any) error {
	return launchRdp(connection, settings)
}

func (s *Service) LaunchHorizon(connection map[string]any, settings map[string]any) error {
	return launchHorizon(connection, settings)
}

func (s *Service) LaunchCitrix(connection map[string]any, settings map[string]any) error {
	return launchCitrix(connection, settings)
}

func (s *Service) LaunchVpn() error {
	return launchVpn()
}
