package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/DlnKot/arc/internal/config"
	"github.com/DlnKot/arc/internal/logging"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	githubOwner = "DlnKot"
	githubRepo  = "ARC"
)

var packageLogger logging.Logger

func setLogger(logger logging.Logger) {
	packageLogger = logger
}

func logInfo(format string, args ...any) {
	if packageLogger != nil {
		packageLogger.Infof(format, args...)
		return
	}
	fmt.Printf("[INFO] "+format+"\n", args...)
}

type Service struct {
	logger   logging.Logger
	ctx      context.Context
	mu       sync.Mutex
	status   Status
	progress int
	cancelCh chan struct{}
}

type UpdateInfo struct {
	Version     string `json:"version"`
	URL         string `json:"url"`
	SHA256      string `json:"sha256,omitempty"`
	Description string `json:"description,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
}

type Status struct {
	UpdateAvailable  bool   `json:"updateAvailable"`
	UpdateDownloaded bool   `json:"updateDownloaded"`
	Version          string `json:"version"`
	ReleaseURL       string `json:"releaseUrl,omitempty"`
	DownloadURL      string `json:"downloadUrl,omitempty"`
	InstallOnQuit    bool   `json:"installOnQuit"`
}

func New(logger logging.Logger) *Service {
	setLogger(logger)
	return &Service{
		cancelCh: make(chan struct{}, 1),
	}
}

func (s *Service) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *Service) CheckForUpdates(useGithub bool, internalServerURL string) (map[string]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentVersion := config.AppVersion
	logInfo("checking for updates, current version: %s", currentVersion)

	var latestVersion string
	var downloadURL string
	var releaseURL string

	if useGithub {
		v, url, err := checkGitHubUpdates(currentVersion)
		if err != nil {
			logInfo("github check failed: %v", err)
		} else {
			latestVersion = v
			downloadURL = url
			releaseURL = fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", githubOwner, githubRepo, v)
			if latestVersion != "" {
				logInfo("github check: update available %s", latestVersion)
			} else {
				logInfo("github check: version is current")
			}
		}
	}

	if latestVersion == "" {
		logInfo("falling back to internal server")
		v, url, err := checkInternalUpdates(currentVersion, internalServerURL)
		if err != nil {
			logInfo("internal server check failed: %v", err)
		} else {
			latestVersion = v
			downloadURL = url
			if latestVersion != "" {
				logInfo("internal server check: update available %s", latestVersion)
			} else {
				logInfo("internal server check: version is current")
			}
		}
	}

	if latestVersion == "" {
		s.status = Status{UpdateAvailable: false}
		return map[string]any{
			"checked":         true,
			"updateAvailable": false,
		}, nil
	}

	s.status = Status{
		UpdateAvailable: true,
		Version:         latestVersion,
		DownloadURL:     downloadURL,
		ReleaseURL:      releaseURL,
	}

	return map[string]any{
		"checked":         true,
		"updateAvailable": true,
		"version":         latestVersion,
		"downloadUrl":     downloadURL,
		"releaseUrl":      releaseURL,
	}, nil
}

func checkGitHubUpdates(currentVersion string) (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubOwner, githubRepo)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("github api returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}

	if err := json.Unmarshal(body, &release); err != nil {
		return "", "", err
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	if latestVersion == "" {
		return "", "", fmt.Errorf("empty tag_name")
	}

	if !isNewerVersion(currentVersion, latestVersion) {
		return "", "", nil
	}

	var downloadURL string
	runtimeOS := runtime.GOOS
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".exe") && runtimeOS == "windows" {
			downloadURL = asset.BrowserDownloadURL
			break
		}
		if strings.HasSuffix(asset.Name, ".dmg") && runtimeOS == "darwin" {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	return latestVersion, downloadURL, nil
}

func checkInternalUpdates(currentVersion string, serverURL string) (string, string, error) {
	if serverURL == "" {
		serverURL = "http://10.230.121.212"
	}

	url := strings.TrimRight(serverURL, "/") + "/update.json"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("internal server returned %d", resp.StatusCode)
	}

	var info UpdateInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", "", err
	}

	if info.Version == "" {
		return "", "", fmt.Errorf("empty version in response")
	}

	if !isNewerVersion(currentVersion, info.Version) {
		return "", "", nil
	}

	return info.Version, info.URL, nil
}

func isNewerVersion(current, latest string) bool {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	for i := 0; i < len(latestParts); i++ {
		curr := 0
		if i < len(currentParts) {
			fmt.Sscanf(currentParts[i], "%d", &curr)
		}
		var latestNum int
		fmt.Sscanf(latestParts[i], "%d", &latestNum)

		if latestNum > curr {
			return true
		}
		if latestNum < curr {
			return false
		}
	}
	return false
}

func (s *Service) DownloadUpdate() error {
	s.mu.Lock()
	if s.status.DownloadURL == "" {
		s.mu.Unlock()
		return fmt.Errorf("no update available")
	}
	downloadURL := s.status.DownloadURL
	s.progress = 0
	s.mu.Unlock()

	logInfo("downloading update from: %s", downloadURL)

	downloadPath := filepath.Join(os.TempDir(), "arc-update")
	if runtime.GOOS == "windows" {
		downloadPath += ".exe"
	} else if runtime.GOOS == "darwin" {
		downloadPath += ".dmg"
	}

	client := &http.Client{Timeout: 600 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		s.emitEvent("updater:error", "Download failed: "+err.Error())
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.emitEvent("updater:error", fmt.Sprintf("Download returned %d", resp.StatusCode))
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	file, err := os.Create(downloadPath)
	if err != nil {
		s.emitEvent("updater:error", "Create file failed: "+err.Error())
		return fmt.Errorf("create file failed: %w", err)
	}
	defer file.Close()

	totalSize := resp.ContentLength
	downloaded := int64(0)

	buffer := make([]byte, 32*1024)
	for {
		select {
		case <-s.cancelCh:
			file.Close()
			os.Remove(downloadPath)
			s.emitEvent("updater:error", "Download cancelled")
			return fmt.Errorf("download cancelled")
		default:
		}

		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, werr := file.Write(buffer[:n]); werr != nil {
				s.emitEvent("updater:error", "Write failed: "+werr.Error())
				return fmt.Errorf("write failed: %w", werr)
			}
			downloaded += int64(n)
			if totalSize > 0 {
				progress := int(float64(downloaded) / float64(totalSize) * 100)
				s.mu.Lock()
				s.progress = progress
				s.mu.Unlock()
				s.emitEvent("updater:progress", progress)
			}
		}
		if err != nil {
			break
		}
	}

	s.mu.Lock()
	s.status.UpdateDownloaded = true
	s.mu.Unlock()

	s.emitEvent("updater:downloaded", downloadPath)
	logInfo("update downloaded to: %s", downloadPath)

	return nil
}

func (s *Service) CancelDownload() error {
	select {
	case s.cancelCh <- struct{}{}:
		return nil
	default:
		return fmt.Errorf("no download in progress")
	}
}

func (s *Service) InstallNow() error {
	s.mu.Lock()
	if !s.status.UpdateDownloaded {
		s.mu.Unlock()
		return fmt.Errorf("update not downloaded")
	}
	s.mu.Unlock()

	logInfo("installing update now...")

	downloadPath := filepath.Join(os.TempDir(), "arc-update")
	if runtime.GOOS == "windows" {
		downloadPath += ".exe"

		cmd := exec.Command(downloadPath, "/S")
		cmd.Stdout = nil
		cmd.Stderr = nil

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start installer: %w", err)
		}

		logInfo("installer started, exiting app...")
		os.Exit(0)
	} else if runtime.GOOS == "darwin" {
		downloadPath += ".dmg"
		if err := runCommand("open", downloadPath); err != nil {
			return fmt.Errorf("failed to open dmg: %w", err)
		}
		os.Exit(0)
	}

	return fmt.Errorf("unsupported platform")
}

func (s *Service) InstallOnQuit() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.status.UpdateDownloaded {
		return fmt.Errorf("update not downloaded")
	}

	s.status.InstallOnQuit = true
	logInfo("will install on quit")

	return nil
}

func (s *Service) CheckAndInstallOnQuit() {
	s.mu.Lock()
	shouldInstall := s.status.InstallOnQuit && s.status.UpdateDownloaded
	s.mu.Unlock()

	if !shouldInstall {
		return
	}

	logInfo("installing on quit...")

	downloadPath := filepath.Join(os.TempDir(), "arc-update")
	if runtime.GOOS == "windows" {
		downloadPath += ".exe"
		_, err := os.Stat(downloadPath)
		if err == nil {
			_, _ = os.StartProcess(downloadPath, []string{}, &os.ProcAttr{})
		}
	} else if runtime.GOOS == "darwin" {
		downloadPath += ".dmg"
		_, err := os.Stat(downloadPath)
		if err == nil {
			_ = runCommand("open", downloadPath)
		}
	}
}

func (s *Service) GetStatus() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return map[string]any{
		"updateAvailable":  s.status.UpdateAvailable,
		"updateDownloaded": s.status.UpdateDownloaded,
		"version":          s.status.Version,
		"releaseUrl":       s.status.ReleaseURL,
		"downloadUrl":      s.status.DownloadURL,
		"installOnQuit":    s.status.InstallOnQuit,
		"progress":         s.progress,
	}
}

func (s *Service) emitEvent(name string, data any) {
	if s.ctx == nil {
		return
	}
	wailsruntime.EventsEmit(s.ctx, name, data)
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Start()
}
