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
	"time"

	"github.com/DlnKot/arc/internal/logging"
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

func startDetached(command string, args ...string) error {
	return startDetachedWithOptions(command, false, args...)
}

func startDetachedNoHide(command string, args ...string) error {
	return startDetachedWithOptions(command, true, args...)
}

func startDetachedWithOptions(command string, showWindow bool, args ...string) error {
	logInfo("launch: %s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	if showWindow {
		cmd.SysProcAttr = nil
	} else {
		applyPlatformProcessAttrs(cmd)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func firstExisting(paths ...string) string {
	for _, path := range paths {
		if fileExists(path) || dirExists(path) {
			return path
		}
	}
	return ""
}

func scheduleDelete(path string, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		_ = os.Remove(path)
	}()
}

func normalizeHTTPSURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	if !strings.HasPrefix(strings.ToLower(trimmed), "http://") && !strings.HasPrefix(strings.ToLower(trimmed), "https://") {
		trimmed = "https://" + trimmed
	}
	u, err := url.Parse(trimmed)
	if err != nil {
		return trimmed
	}
	u.RawQuery = ""
	u.Fragment = ""
	return strings.TrimRight(u.String(), "/")
}

func validateHostOrURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("host is required")
	}
	if len(trimmed) > 512 {
		return "", errors.New("host is too long")
	}
	if strings.ContainsAny(trimmed, "\n\r;$|`") {
		return "", errors.New("host contains invalid characters")
	}
	host := stripHostPrefix(trimmed)
	return host, nil
}

func stripHostPrefix(raw string) string {
	trimmed := strings.TrimSpace(raw)
	lower := strings.ToLower(trimmed)
	prefixes := []string{"rdp://", "https://", "http://", "mstsc://"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			return trimmed[len(prefix):]
		}
	}
	return trimmed
}

func validateUsername(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}
	if len(trimmed) > 256 {
		return "", errors.New("username is too long")
	}
	if strings.ContainsAny(trimmed, "\n\r;$|`") {
		return "", errors.New("username contains invalid characters")
	}
	return trimmed, nil
}

func writeTempFile(prefix string, suffix string, content string) (string, error) {
	file, err := os.CreateTemp("", prefix+"-*"+suffix)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		return "", err
	}
	return file.Name(), nil
}

func userHomeApp(pathParts ...string) string {
	home, _ := os.UserHomeDir()
	parts := append([]string{home}, pathParts...)
	return filepath.Join(parts...)
}

func boolSetting(source map[string]any, key string, def bool) bool {
	value, ok := source[key]
	if !ok {
		return def
	}
	b, ok := value.(bool)
	if !ok {
		return def
	}
	return b
}

func stringSetting(source map[string]any, key string, def string) string {
	value := strings.TrimSpace(asString(source[key]))
	if value == "" {
		return def
	}
	return value
}

func nestedMap(source map[string]any, key string) map[string]any {
	value, _ := source[key].(map[string]any)
	if value == nil {
		return map[string]any{}
	}
	return value
}

func platformNotSupported(feature string) error {
	return fmt.Errorf("%s is not supported on platform %s", feature, runtime.GOOS)
}

func splitArgs(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var args []string
	var current strings.Builder
	inQuotes := false

	for _, r := range raw {
		switch {
		case r == '"':
			inQuotes = !inQuotes
		case !inQuotes && (r == ' ' || r == '\t' || r == '\n' || r == '\r'):
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func asString(value any) string {
	if s, ok := value.(string); ok {
		return s
	}
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func ternaryInt(condition bool, ifTrue int, ifFalse int) int {
	if condition {
		return ifTrue
	}
	return ifFalse
}

func ternaryString(condition bool, ifTrue string, ifFalse string) string {
	if condition {
		return ifTrue
	}
	return ifFalse
}
