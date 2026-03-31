package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/DlnKot/arc/internal/config"
	"github.com/DlnKot/arc/internal/logging"
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
	logger logging.Logger
	status Status
}

type UpdateInfo struct {
	Version     string `json:"version"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Status struct {
	UpdateAvailable  bool   `json:"updateAvailable"`
	UpdateDownloaded bool   `json:"updateDownloaded"`
	Version          string `json:"version"`
	ReleaseURL       string `json:"releaseUrl,omitempty"`
	DownloadURL      string `json:"downloadUrl,omitempty"`
}

func New(logger logging.Logger) *Service {
	setLogger(logger)
	return &Service{logger: logger}
}

func (s *Service) CheckForUpdates(useGithub bool, internalServerURL string) (map[string]any, error) {
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
		}
	}

	if latestVersion == "" {
		v, url, err := checkInternalUpdates(currentVersion, internalServerURL)
		if err != nil {
			logInfo("internal server check failed: %v", err)
		} else {
			latestVersion = v
			downloadURL = url
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
	if s.status.DownloadURL == "" {
		return fmt.Errorf("no update available")
	}

	logInfo("downloading update from: %s", s.status.DownloadURL)

	downloadPath := filepath.Join(os.TempDir(), "arc-update")
	if runtime.GOOS == "windows" {
		downloadPath += ".exe"
	} else if runtime.GOOS == "darwin" {
		downloadPath += ".dmg"
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(s.status.DownloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	file, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("create file failed: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	s.status.UpdateDownloaded = true
	logInfo("update downloaded to: %s", downloadPath)

	return nil
}

func (s *Service) InstallUpdate() error {
	if !s.status.UpdateDownloaded {
		return fmt.Errorf("update not downloaded")
	}

	logInfo("installing update...")

	downloadPath := filepath.Join(os.TempDir(), "arc-update")
	if runtime.GOOS == "windows" {
		downloadPath += ".exe"
		proc, err := os.StartProcess(downloadPath, []string{}, &os.ProcAttr{})
		if err != nil {
			return fmt.Errorf("failed to start installer: %w", err)
		}
		proc.Release()
		os.Exit(0)
	} else if runtime.GOOS == "darwin" {
		downloadPath += ".dmg"
		return fmt.Errorf("macOS update requires manual install: %s", downloadPath)
	}

	return fmt.Errorf("unsupported platform")
}

func (s *Service) GetStatus() map[string]any {
	return map[string]any{
		"updateAvailable":  s.status.UpdateAvailable,
		"updateDownloaded": s.status.UpdateDownloaded,
		"version":          s.status.Version,
		"releaseUrl":       s.status.ReleaseURL,
		"downloadUrl":      s.status.DownloadURL,
	}
}
