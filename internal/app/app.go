package app

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/DlnKot/arc/internal/config"
	"github.com/DlnKot/arc/internal/domain"
	"github.com/DlnKot/arc/internal/launchers"
	"github.com/DlnKot/arc/internal/logging"
	"github.com/DlnKot/arc/internal/network"
	"github.com/DlnKot/arc/internal/store"
	"github.com/DlnKot/arc/internal/updater"
)

type App struct {
	ctx       context.Context
	store     *store.Service
	network   *network.Service
	launchers *launchers.Service
	updater   *updater.Service
	logger    logging.Logger
}

func New(storeSvc *store.Service, networkSvc *network.Service, launcherSvc *launchers.Service, updaterSvc *updater.Service, logger logging.Logger) *App {
	return &App{
		store:     storeSvc,
		network:   networkSvc,
		launchers: launcherSvc,
		updater:   updaterSvc,
		logger:    logger,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	if err := a.store.Load(); err != nil {
		if a.logger != nil {
			a.logger.Errorf("startup loadStore error: %v", err)
		}
	}

	a.updater.SetContext(ctx)

	go func() {
		time.Sleep(3 * time.Second)
		settings := a.store.GetSettings()
		if updates, ok := settings["updates"].(map[string]any); ok {
			if autoCheck, ok := updates["autoCheck"].(bool); ok && autoCheck {
				useGithub := false
				internalURL := "http://10.230.121.212"
				if v, ok := updates["useGithub"].(bool); ok {
					useGithub = v
				}
				if v, ok := updates["internalServerUrl"].(string); ok && v != "" {
					internalURL = v
				}
				a.updater.CheckForUpdates(useGithub, internalURL)
			}
		}
	}()
}

func (a *App) GetConnections() domain.Result[[]map[string]any] {
	return ok(a.store.GetConnections())
}

func (a *App) SaveConnection(connection map[string]any) domain.Result[bool] {
	if err := a.store.SaveConnection(connection); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) DeleteConnection(id string) domain.Result[bool] {
	if err := a.store.DeleteConnection(id); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) ResetDefaultConnections() domain.Result[[]map[string]any] {
	if err := a.store.ResetDefaultConnections(); err != nil {
		return fail[[]map[string]any](err)
	}
	return ok(a.store.GetConnections())
}

func (a *App) GetSettings() domain.Result[map[string]any] {
	return ok(a.store.GetSettings())
}

func (a *App) SaveSettings(settings map[string]any) domain.Result[bool] {
	if err := a.store.SaveSettings(settings); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) LaunchRdp(connection map[string]any, settings map[string]any) domain.Result[bool] {
	if err := a.launchers.LaunchRdp(connection, settings); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) LaunchHorizon(connection map[string]any, settings map[string]any) domain.Result[bool] {
	if err := a.launchers.LaunchHorizon(connection, settings); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) LaunchCitrix(connection map[string]any, settings map[string]any) domain.Result[bool] {
	if err := a.launchers.LaunchCitrix(connection, settings); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) LaunchVpn() domain.Result[bool] {
	if err := a.launchers.LaunchVpn(); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) NetworkGeo() domain.Result[map[string]any] {
	data, err := a.network.Geo()
	if err != nil {
		return fail[map[string]any](err)
	}
	return ok(data)
}

func (a *App) NetworkPing(host string, packets int) domain.Result[map[string]any] {
	data, err := a.network.Ping(host, packets)
	if err != nil {
		return fail[map[string]any](err)
	}
	return ok(data)
}

func (a *App) CheckForUpdates() domain.Result[map[string]any] {
	settings := a.store.GetSettings()
	useGithub := false
	if updates, ok := settings["updates"].(map[string]any); ok {
		if v, ok := updates["useGithub"].(bool); ok {
			useGithub = v
		}
	}
	internalServerURL := "http://10.230.121.212"
	if updates, ok := settings["updates"].(map[string]any); ok {
		if v, ok := updates["internalServerUrl"].(string); ok && v != "" {
			internalServerURL = v
		}
	}
	data, err := a.updater.CheckForUpdates(useGithub, internalServerURL)
	if err != nil {
		return fail[map[string]any](err)
	}
	return ok(data)
}

func (a *App) DownloadUpdate() domain.Result[bool] {
	if err := a.updater.DownloadUpdate(); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) CancelDownload() domain.Result[bool] {
	if err := a.updater.CancelDownload(); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) InstallNow() domain.Result[bool] {
	if err := a.updater.InstallNow(); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) InstallOnQuit() domain.Result[bool] {
	if err := a.updater.InstallOnQuit(); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (a *App) Shutdown(ctx context.Context) {
	a.updater.CheckAndInstallOnQuit()
}

func (a *App) GetUpdateStatus() domain.Result[map[string]any] {
	return ok(cloneMap(a.updater.GetStatus()))
}

func (a *App) GetVersion() domain.Result[string] {
	return ok(config.AppVersion)
}

func (a *App) GetPlatform() domain.Result[string] {
	return ok(runtime.GOOS)
}

func (a *App) Log(level string, message string) domain.Result[bool] {
	if a.logger != nil {
		switch strings.ToLower(strings.TrimSpace(level)) {
		case "warn":
			a.logger.Warnf("%s", message)
		case "error":
			a.logger.Errorf("%s", message)
		default:
			a.logger.Infof("%s", message)
		}
	}
	return ok(true)
}

func ok[T any](data T) domain.Result[T] {
	return domain.Result[T]{Success: true, Data: data}
}

func fail[T any](err error) domain.Result[T] {
	return domain.Result[T]{Success: false, Error: err.Error()}
}

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}
