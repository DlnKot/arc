package launchers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func launchRdp(connection map[string]any, settings map[string]any) error {
	host, err := validateHostOrURL(asString(connection["host"]))
	if err != nil {
		return err
	}
	username, err := validateUsername(asString(connection["username"]))
	if err != nil {
		return err
	}

	rdpSettings := nestedMap(settings, "rdp")
	if len(rdpSettings) == 0 {
		rdpSettings = settings
	}

	resolution := stringSetting(rdpSettings, "resolution", "1920x1080")
	width, height := 1920, 1080
	fullscreen := resolution == "fullscreen" || boolSetting(rdpSettings, "startFullScreen", false)
	if !fullscreen && strings.Contains(resolution, "x") {
		_, _ = fmt.Sscanf(resolution, "%dx%d", &width, &height)
	}

	redirect := nestedMap(rdpSettings, "redirect")
	audio := nestedMap(rdpSettings, "audio")
	perf := nestedMap(rdpSettings, "performance")

	lines := []string{
		fmt.Sprintf("full address:s:%s", host),
		fmt.Sprintf("username:s:%s", username),
		fmt.Sprintf("screen mode id:i:%d", ternaryInt(fullscreen, 2, 1)),
		fmt.Sprintf("desktopwidth:i:%d", width),
		fmt.Sprintf("desktopheight:i:%d", height),
		fmt.Sprintf("session bpp:i:%s", stringSetting(rdpSettings, "colorDepth", "32")),
		fmt.Sprintf("winposstr:s:0,3,0,0,%d,%d", width, height),
		fmt.Sprintf("use multimon:i:%d", ternaryInt(boolSetting(rdpSettings, "multimon", false), 1, 0)),
		fmt.Sprintf("span monitors:i:%d", ternaryInt(boolSetting(rdpSettings, "span", false), 1, 0)),
		fmt.Sprintf("redirectclipboard:i:%d", ternaryInt(boolSetting(rdpSettings, "clipboard", true), 1, 0)),
		fmt.Sprintf("drivestoredirect:s:%s", ternaryString(boolSetting(rdpSettings, "driveMapping", false), "*", "")),
		fmt.Sprintf("redirectprinters:i:%d", ternaryInt(boolSetting(redirect, "printers", true), 1, 0)),
		fmt.Sprintf("redirectsmartcards:i:%d", ternaryInt(boolSetting(redirect, "smartcards", true), 1, 0)),
		fmt.Sprintf("redirectwebauthn:i:%d", ternaryInt(boolSetting(redirect, "webauthn", true), 1, 0)),
		fmt.Sprintf("audiomode:i:%d", ternaryInt(boolSetting(audio, "playback", true), 0, 2)),
		fmt.Sprintf("audiocapturemode:i:%d", ternaryInt(boolSetting(audio, "capture", false), 1, 0)),
		fmt.Sprintf("disable wallpaper:i:%d", ternaryInt(!boolSetting(perf, "wallpaper", true), 1, 0)),
		fmt.Sprintf("allow font smoothing:i:%d", ternaryInt(boolSetting(perf, "fontSmoothing", true), 1, 0)),
		fmt.Sprintf("allow desktop composition:i:%d", ternaryInt(boolSetting(perf, "desktopComposition", true), 1, 0)),
		fmt.Sprintf("disable full window drag:i:%d", ternaryInt(!boolSetting(perf, "fullWindowDrag", true), 1, 0)),
		fmt.Sprintf("disable menu anims:i:%d", ternaryInt(!boolSetting(perf, "menuAnimations", true), 1, 0)),
		fmt.Sprintf("prompt for credentials:i:%d", ternaryInt(boolSetting(rdpSettings, "promptCredentials", true), 1, 0)),
		fmt.Sprintf("administrative session:i:%d", ternaryInt(boolSetting(rdpSettings, "useAdminSession", false), 1, 0)),
	}
	if custom := strings.TrimSpace(stringSetting(rdpSettings, "customFlags", "")); custom != "" {
		lines = append(lines, splitArgs(custom)...)
	}

	rdpFile, err := writeTempFile("arc-rdp", ".rdp", strings.Join(lines, "\n"))
	if err != nil {
		return err
	}
	scheduleDelete(rdpFile, 15*time.Second)

	switch runtime.GOOS {
	case "windows":
		return startDetached("mstsc.exe", rdpFile)
	case "darwin":
		if err := startDetached("open", "-b", "com.microsoft.rdc.macos", rdpFile); err == nil {
			return nil
		}
		appPath := firstExisting(
			"/Applications/Windows App.app",
			"/Applications/Microsoft Remote Desktop.app",
			userHomeApp("Applications", "Windows App.app"),
			userHomeApp("Applications", "Microsoft Remote Desktop.app"),
		)
		if appPath != "" {
			return startDetached("open", "-a", filepath.Base(strings.TrimSuffix(appPath, ".app")), rdpFile)
		}
		return startDetached("open", rdpFile)
	default:
		return startDetached("xfreerdp", "/v:"+host)
	}
}
