package launchers

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func launchCitrix(connection map[string]any, settings map[string]any) error {
	citrixSettings := nestedMap(settings, "citrix")
	if len(citrixSettings) == 0 {
		citrixSettings = settings
	}

	storeURL := strings.TrimSpace(asString(connection["storeUrl"]))
	if storeURL == "" {
		storeURL = stringSetting(citrixSettings, "storeUrl", "")
	}
	storeURL = normalizeStorefrontAddress(storeURL)
	if storeURL == "" {
		return errors.New("Citrix Store URL is required")
	}

	resourceName := stringSetting(citrixSettings, "resourceName", "")
	customFlags := splitArgs(stringSetting(citrixSettings, "customFlags", ""))

	if runtime.GOOS == "windows" {
		exe := findCitrixExecutableWindows(stringSetting(citrixSettings, "customPath", ""))
		if exe == "" {
			return errors.New("Citrix Workspace not found")
		}

		storeAlreadyConfigured := boolSetting(citrixSettings, "storeAlreadyConfigured", false)
		var args []string
		if !storeAlreadyConfigured {
			args = []string{"-store", storeURL}
		}
		if resourceName != "" {
			args = append(args, "-launch", resourceName, "-quiet")
		}
		args = append(args, customFlags...)
		return startDetachedNoHide(exe, args...)
	}

	if runtime.GOOS == "darwin" {
		app := firstExisting(
			"/Applications/Citrix Workspace.app",
			"/Applications/Citrix Receiver.app",
			userHomeApp("Applications", "Citrix Workspace.app"),
		)
		if app == "" {
			return errors.New("Citrix Workspace not found on macOS")
		}

		accountName := stringSetting(citrixSettings, "accountName", stringSetting(connection, "name", "Store"))
		createURL := buildCitrixCreateAccountURL(accountName, normalizeStorefrontDiscoveryAddress(storeURL))
		_ = startDetached("open", createURL)

		args := []string{"-a", app}
		if resourceName != "" || len(customFlags) > 0 {
			args = append(args, "--args")
			if resourceName != "" {
				args = append(args, "-launch", resourceName)
			}
			args = append(args, customFlags...)
		}
		return startDetached("open", args...)
	}

	return platformNotSupported("Citrix launcher")
}

func findCitrixExecutableWindows(customPath string) string {
	if fileExists(customPath) {
		return customPath
	}

	programFiles := os.Getenv("ProgramFiles")
	programFiles86 := os.Getenv("ProgramFiles(x86)")
	programData := os.Getenv("ProgramData")
	localAppData := os.Getenv("LOCALAPPDATA")

	candidates := []string{
		filepath.Join(programFiles86, "Citrix", "ICA Client", "SelfServicePlugin", "SelfService.exe"),
		filepath.Join(programFiles, "Citrix", "ICA Client", "SelfServicePlugin", "SelfService.exe"),
		filepath.Join(programFiles86, "Citrix", "ICA Client", "selfservice.exe"),
		filepath.Join(programFiles, "Citrix", "ICA Client", "selfservice.exe"),
		filepath.Join(localAppData, "Citrix", "ICA Client", "SelfServicePlugin", "SelfService.exe"),
		filepath.Join(programData, "Citrix", "ICA Client", "SelfServicePlugin", "SelfService.exe"),
		filepath.Join(programFiles86, "Citrix", "ICA Client", "Citrix Workspace.exe"),
		filepath.Join(programFiles, "Citrix", "ICA Client", "Citrix Workspace.exe"),
		filepath.Join(programFiles86, "Citrix", "Workspace", "SelfService.exe"),
		filepath.Join(programFiles, "Citrix", "Workspace", "SelfService.exe"),
		filepath.Join(programFiles86, "Citrix", "Receiver", "SelfService.exe"),
		filepath.Join(programFiles, "Citrix", "Receiver", "SelfService.exe"),
	}
	if path := firstExisting(candidates...); path != "" {
		return path
	}
	if resolved, err := exec.LookPath("SelfService.exe"); err == nil {
		return resolved
	}
	if resolved, err := exec.LookPath("selfservice.exe"); err == nil {
		return resolved
	}
	if resolved, err := exec.LookPath("Citrix Workspace.exe"); err == nil {
		return resolved
	}
	return ""
}

func normalizeStorefrontAddress(raw string) string {
	normalized := normalizeHTTPSURL(raw)
	if normalized == "" {
		return ""
	}
	u, err := url.Parse(normalized)
	if err != nil {
		return strings.TrimRight(normalized, "/")
	}
	path := strings.TrimRight(u.Path, "/")
	if strings.HasSuffix(strings.ToLower(path), "/discovery") {
		path = path[:len(path)-len("/discovery")]
	}
	u.Path = path
	u.RawQuery = ""
	u.Fragment = ""
	return strings.TrimRight(u.String(), "/")
}

func normalizeStorefrontDiscoveryAddress(raw string) string {
	base := normalizeStorefrontAddress(raw)
	if base == "" {
		return ""
	}
	if strings.HasSuffix(strings.ToLower(base), "/discovery") {
		return base
	}
	return strings.TrimRight(base, "/") + "/discovery"
}

func buildCitrixCreateAccountURL(accountName string, addressURL string) string {
	return fmt.Sprintf("citrixreceiver://createaccount?name=%s&address=%s", url.QueryEscape(accountName), url.QueryEscape(addressURL))
}
