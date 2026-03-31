package launchers

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func launchHorizon(connection map[string]any, settings map[string]any) error {
	host, err := validateHostOrURL(asString(connection["host"]))
	if err != nil {
		return err
	}
	username, err := validateUsername(asString(connection["username"]))
	if err != nil {
		return err
	}

	horizonSettings := nestedMap(settings, "horizon")
	if len(horizonSettings) == 0 {
		horizonSettings = settings
	}

	if runtime.GOOS == "windows" {
		exe := findHorizonExecutableWindows(stringSetting(horizonSettings, "customPath", ""))
		if exe == "" {
			return errors.New("VMware Horizon Client not found")
		}
		return startDetached(exe, buildHorizonArgsWindows(host, username, asString(connection["desktopPool"]), horizonSettings)...)
	}

	if runtime.GOOS == "darwin" {
		appPath := firstExisting(
			"/Applications/VMware Horizon Client.app",
			"/Applications/VMware Horizon.app",
			userHomeApp("Applications", "VMware Horizon Client.app"),
		)
		if appPath == "" {
			return errors.New("VMware Horizon Client not found")
		}

		cli := firstExisting(
			filepath.Join(appPath, "Contents", "Launchers", "viewclient-macosx-arm64"),
			filepath.Join(appPath, "Contents", "Launchers", "viewclient-macosx"),
		)
		args := buildHorizonArgsMac(host, username, asString(connection["desktopPool"]), horizonSettings)
		if cli != "" {
			return startDetached(cli, args...)
		}
		return startDetached("open", append([]string{"-a", appPath, "--args"}, args...)...)
	}

	return platformNotSupported("Horizon launcher")
}

func buildHorizonArgsWindows(host string, username string, desktopPool string, settings map[string]any) []string {
	args := []string{"--serverURL=" + normalizeHTTPSURL(host)}
	if desktopPool != "" {
		args = append(args, "--desktopName="+desktopPool)
	}
	if user := stripDomainFromUsername(username); user != "" {
		args = append(args, "--userName="+user)
	}
	if value := stringSetting(settings, "appName", ""); value != "" {
		args = append(args, "--appName="+value)
	}
	if value := stringSetting(settings, "desktopProtocol", ""); value != "" {
		args = append(args, "--desktopProtocol="+value)
	}
	if value := stringSetting(settings, "desktopLayout", ""); value != "" {
		args = append(args, "--desktopLayout="+value)
	}
	if value := stringSetting(settings, "monitors", ""); value != "" {
		args = append(args, "--monitors="+value)
	}
	if boolSetting(settings, "unattended", false) {
		args = append(args, "--unattended")
	}
	if boolSetting(settings, "nonInteractive", false) {
		args = append(args, "--nonInteractive")
	}
	if boolSetting(settings, "launchMinimized", false) {
		args = append(args, "--launchMinimized")
	}
	if boolSetting(settings, "loginAsCurrentUser", false) {
		args = append(args, "--loginAsCurrentUser=true")
	}
	if boolSetting(settings, "hideClientAfterLaunchSession", false) {
		args = append(args, "--hideClientAfterLaunchSession=true")
	}
	if boolSetting(settings, "useExisting", false) {
		args = append(args, "--useExisting")
	}
	if boolSetting(settings, "singleAutoConnect", false) {
		args = append(args, "--singleAutoConnect")
	}
	args = append(args, splitArgs(stringSetting(settings, "customFlags", ""))...)
	return args
}

func buildHorizonArgsMac(host string, username string, desktopPool string, settings map[string]any) []string {
	args := []string{"-serverURL", normalizeHTTPSURL(host)}
	if desktopPool != "" {
		args = append(args, "-desktopName", desktopPool)
	}
	if user := stripDomainFromUsername(username); user != "" {
		args = append(args, "-userName", user)
	}
	if value := stringSetting(settings, "appName", ""); value != "" {
		args = append(args, "-appName", value)
	}
	if value := stringSetting(settings, "desktopProtocol", ""); value != "" {
		args = append(args, "-desktopProtocol", value)
	}
	if value := stringSetting(settings, "desktopLayout", ""); value != "" {
		args = append(args, "-desktopLayout", value)
	}
	if value := stringSetting(settings, "monitors", ""); value != "" {
		args = append(args, "-monitors", value)
	}
	if boolSetting(settings, "unattended", false) {
		args = append(args, "-unattended")
	}
	if boolSetting(settings, "loginAsCurrentUser", false) {
		args = append(args, "-loginAsCurrentUser")
	}
	args = append(args, splitArgs(stringSetting(settings, "customFlags", ""))...)
	return args
}

func stripDomainFromUsername(username string) string {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		return ""
	}
	if strings.Contains(trimmed, "\\") {
		parts := strings.Split(trimmed, "\\")
		return strings.TrimSpace(parts[len(parts)-1])
	}
	if strings.Contains(trimmed, "/") {
		parts := strings.Split(trimmed, "/")
		return strings.TrimSpace(parts[len(parts)-1])
	}
	return trimmed
}

func findHorizonExecutableWindows(customPath string) string {
	if fileExists(customPath) {
		return customPath
	}

	programFiles := os.Getenv("ProgramFiles")
	programFiles86 := os.Getenv("ProgramFiles(x86)")
	programData := os.Getenv("ProgramData")
	localAppData := os.Getenv("LOCALAPPDATA")

	candidates := []string{
		filepath.Join(programFiles, "VMware", "VMware Horizon Client", "bin", "vmware-view.exe"),
		filepath.Join(programFiles86, "VMware", "VMware Horizon Client", "bin", "vmware-view.exe"),
		filepath.Join(programFiles, "VMware", "VMware Horizon", "View", "Client", "bin", "vmware-view.exe"),
		filepath.Join(programFiles86, "VMware", "VMware Horizon", "View", "Client", "bin", "vmware-view.exe"),
		filepath.Join(programData, "VMware", "VMware Horizon", "View", "Client", "bin", "vmware-view.exe"),
		filepath.Join(localAppData, "VMware", "VMware Horizon Client", "bin", "vmware-view.exe"),
	}
	if path := firstExisting(candidates...); path != "" {
		return path
	}
	if resolved, err := exec.LookPath("vmware-view.exe"); err == nil {
		return resolved
	}
	return ""
}
