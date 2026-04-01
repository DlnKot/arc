package launchers

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestStripHostPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"192.168.1.1", "192.168.1.1"},
		{"server.domain.com", "server.domain.com"},
		{"RDP://192.168.1.1", "192.168.1.1"},
		{"rdp://server.domain.com", "server.domain.com"},
		{"HTTPS://server.domain.com/path", "server.domain.com/path"},
		{"http://server.domain.com/path", "server.domain.com/path"},
		{"mstsc://server", "server"},
		{"  rdp://192.168.1.1  ", "192.168.1.1"},
		{"RDP://DOMAIN\\user@server", "DOMAIN\\user@server"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := stripHostPrefix(tt.input)
			if got != tt.expected {
				t.Errorf("stripHostPrefix(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestStripDomainFromUsername(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ivanov", "ivanov"},
		{"MOSCOW\\ivanov", "ivanov"},
		{"REGIONS/ivanov", "ivanov"},
		{"user@domain.com", "user@domain.com"},
		{"  MOSCOW\\ivanov  ", "ivanov"},
		{"", ""},
		{"  ", ""},
		{"a\\b\\c", "c"},
		{"a/b/c", "c"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := stripDomainFromUsername(tt.input)
			if got != tt.expected {
				t.Errorf("stripDomainFromUsername(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestValidateHostOrURL(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
		wantHost string
	}{
		{"192.168.1.1", false, "192.168.1.1"},
		{"server.domain.com", false, "server.domain.com"},
		{"RDP://192.168.1.1", false, "192.168.1.1"},
		{"rdp://server", false, "server"},
		{"", true, ""},
		{"   ", true, ""},
		{"host\nwith\nnewlines", true, ""},
		{"host;with;semicolons", true, ""},
		{"host$with$dollars", true, ""},
		{"host|with|pipes", true, ""},
		{"host`with`backticks", true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			host, err := validateHostOrURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHostOrURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && host != tt.wantHost {
				t.Errorf("validateHostOrURL(%q) = %q, want %q", tt.input, host, tt.wantHost)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		want    string
	}{
		{"ivanov", false, "ivanov"},
		{"MOSCOW\\ivanov", false, "MOSCOW\\ivanov"},
		{"", false, ""},
		{"   ", false, ""},
		{"user@domain.com", false, "user@domain.com"},
		{"user\nwith\nnewlines", true, ""},
		{"user;with;semicolons", true, ""},
		{"user$with$dollars", true, ""},
		{"user`with`backticks", true, ""},
		{strings.Repeat("a", 257), true, ""},
		{strings.Repeat("a", 256), false, strings.Repeat("a", 256)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := validateUsername(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUsername(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("validateUsername(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeHTTPSURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://example.com/path?query=1#frag", "https://example.com/path"},
		{"http://example.com/path", "http://example.com/path"},
		{"example.com/path", "https://example.com/path"},
		{"example.com", "https://example.com"},
		{"", ""},
		{"  ", ""},
		{"HTTPS://EXAMPLE.COM", "https://EXAMPLE.COM"},
		{"https://example.com/", "https://example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeHTTPSURL(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeHTTPSURL(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNormalizeStorefrontAddress(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://store.company.com/Citrix/Store/discovery", "https://store.company.com/Citrix/Store"},
		{"https://store.company.com/Citrix/Store", "https://store.company.com/Citrix/Store"},
		{"https://store.company.com/Citrix/Store/", "https://store.company.com/Citrix/Store"},
		{"https://store.company.com/Citrix/Store?query=1", "https://store.company.com/Citrix/Store"},
		{"https://store.company.com/Citrix/Store#fragment", "https://store.company.com/Citrix/Store"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeStorefrontAddress(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeStorefrontAddress(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNormalizeStorefrontDiscoveryAddress(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://store.company.com/Citrix/Store", "https://store.company.com/Citrix/Store/discovery"},
		{"https://store.company.com/Citrix/Store/discovery", "https://store.company.com/Citrix/Store/discovery"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeStorefrontDiscoveryAddress(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeStorefrontDiscoveryAddress(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"flag1 value1 flag2 value2", []string{"flag1", "value1", "flag2", "value2"}},
		{`"flag with spaces" value`, []string{"flag with spaces", "value"}},
		{"", nil},
		{"  ", nil},
		{"  flag1  value1  ", []string{"flag1", "value1"}},
		{`flag1:"value with spaces"`, []string{"flag1:value with spaces"}},
		{"-flag1 value1\t-flag2 value2", []string{"-flag1", "value1", "-flag2", "value2"}},
		{"-flag\tvalue", []string{"-flag", "value"}},
		{"-flag\nvalue", []string{"-flag", "value"}},
		{"-flag\rvalue", []string{"-flag", "value"}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitArgs(tt.input)
			if len(got) == 0 && len(tt.expected) == 0 {
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("splitArgs(%q) len = %d, want %d; got %v", tt.input, len(got), len(tt.expected), got)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("splitArgs(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestBuildHorizonArgsWindows(t *testing.T) {
	settings := map[string]any{
		"appName":                      "Calculator",
		"desktopProtocol":              "PCoIP",
		"desktopLayout":                "fullscreen",
		"monitors":                     "1,2",
		"unattended":                   true,
		"nonInteractive":               true,
		"launchMinimized":              true,
		"loginAsCurrentUser":           true,
		"hideClientAfterLaunchSession": true,
		"useExisting":                  true,
		"singleAutoConnect":            true,
	}

	args := buildHorizonArgsWindows("server.domain.com", "MOSCOW\\ivanov", "Desktop Pool 1", settings)

	assertArg := func(args []string, prefix string) {
		for _, arg := range args {
			if strings.HasPrefix(arg, prefix) {
				return
			}
		}
		t.Errorf("expected arg with prefix %q in %v", prefix, args)
	}

	assertArg(args, "--serverURL=")
	assertArg(args, "--desktopName=Desktop Pool 1")
	assertArg(args, "--userName=ivanov")
	assertArg(args, "--appName=Calculator")
	assertArg(args, "--desktopProtocol=PCoIP")
	assertArg(args, "--desktopLayout=fullscreen")
	assertArg(args, "--monitors=1,2")
	assertArg(args, "--unattended")
	assertArg(args, "--nonInteractive")
	assertArg(args, "--launchMinimized")
	assertArg(args, "--loginAsCurrentUser=true")
	assertArg(args, "--hideClientAfterLaunchSession=true")
	assertArg(args, "--useExisting")
	assertArg(args, "--singleAutoConnect")
}

func TestBuildHorizonArgsWindowsDefaults(t *testing.T) {
	args := buildHorizonArgsWindows("server.domain.com", "ivanov", "", map[string]any{})

	for _, arg := range args {
		if strings.HasPrefix(arg, "--unattended") ||
			strings.HasPrefix(arg, "--nonInteractive") ||
			strings.HasPrefix(arg, "--launchMinimized") ||
			strings.HasPrefix(arg, "--loginAsCurrentUser") ||
			strings.HasPrefix(arg, "--useExisting") ||
			strings.HasPrefix(arg, "--singleAutoConnect") {
			t.Errorf("unexpected default arg %q in %v", arg, args)
		}
	}
}

func TestBuildHorizonArgsMac(t *testing.T) {
	settings := map[string]any{
		"appName":            "Calculator",
		"desktopProtocol":    "RDP",
		"desktopLayout":      "multimonitor",
		"monitors":           "1",
		"unattended":         true,
		"loginAsCurrentUser": true,
	}

	args := buildHorizonArgsMac("server.domain.com", "MOSCOW\\ivanov", "Desktop Pool 1", settings)

	assertArg := func(args []string, val string) bool {
		for _, arg := range args {
			if arg == val {
				return true
			}
		}
		return false
	}

	if !assertArg(args, "-serverURL") {
		t.Errorf("missing -serverURL in %v", args)
	}
	if !assertArg(args, "https://server.domain.com") {
		t.Errorf("missing normalized URL in %v", args)
	}
	if !assertArg(args, "-desktopName") {
		t.Errorf("missing -desktopName in %v", args)
	}
	if !assertArg(args, "Desktop Pool 1") {
		t.Errorf("missing desktop pool in %v", args)
	}
	if !assertArg(args, "-userName") {
		t.Errorf("missing -userName in %v", args)
	}
	if !assertArg(args, "ivanov") {
		t.Errorf("missing stripped username in %v", args)
	}
	if !assertArg(args, "-appName") {
		t.Errorf("missing -appName in %v", args)
	}
	if !assertArg(args, "Calculator") {
		t.Errorf("missing app name in %v", args)
	}
	if !assertArg(args, "-unattended") {
		t.Errorf("missing -unattended in %v", args)
	}
	if !assertArg(args, "-loginAsCurrentUser") {
		t.Errorf("missing -loginAsCurrentUser in %v", args)
	}
}

func TestBuildHorizonArgsMacCustomFlags(t *testing.T) {
	settings := map[string]any{
		"customFlags": "-extraFlag value",
	}
	args := buildHorizonArgsMac("server.com", "user", "", settings)

	hasCustom := false
	for _, arg := range args {
		if arg == "-extraFlag" || arg == "value" {
			hasCustom = true
			break
		}
	}
	if !hasCustom {
		t.Errorf("custom flags not found in args: %v", args)
	}
}

func TestBuildCitrixCreateAccountURL(t *testing.T) {
	url := buildCitrixCreateAccountURL("CorpStore", "https://store.company.com/discovery")
	if !strings.HasPrefix(url, "citrixreceiver://createaccount?") {
		t.Errorf("URL should start with scheme, got %q", url)
	}
	if !strings.Contains(url, "name=") {
		t.Errorf("URL should contain name param, got %q", url)
	}
	if !strings.Contains(url, "address=") {
		t.Errorf("URL should contain address param, got %q", url)
	}
	if !strings.Contains(url, "CorpStore") {
		t.Errorf("URL should contain encoded account name, got %q", url)
	}
}

func TestStringSetting(t *testing.T) {
	tests := []struct {
		source   map[string]any
		key      string
		def      string
		expected string
	}{
		{nil, "key", "default", "default"},
		{map[string]any{}, "key", "default", "default"},
		{map[string]any{"key": "value"}, "key", "default", "value"},
		{map[string]any{"key": "  value  "}, "key", "default", "value"},
		{map[string]any{"key": "value"}, "missing", "default", "default"},
		{map[string]any{"key": ""}, "key", "default", "default"},
		{map[string]any{"key": 123}, "key", "default", "123"},
	}
	for _, tt := range tests {
		got := stringSetting(tt.source, tt.key, tt.def)
		if got != tt.expected {
			t.Errorf("stringSetting(%v, %q, %q) = %q, want %q", tt.source, tt.key, tt.def, got, tt.expected)
		}
	}
}

func TestBoolSetting(t *testing.T) {
	tests := []struct {
		source   map[string]any
		key      string
		def      bool
		expected bool
	}{
		{nil, "key", false, false},
		{nil, "key", true, true},
		{map[string]any{}, "key", true, true},
		{map[string]any{"key": true}, "key", false, true},
		{map[string]any{"key": false}, "key", true, false},
		{map[string]any{"key": "true"}, "key", false, false},
		{map[string]any{"key": 1}, "key", false, false},
	}
	for _, tt := range tests {
		got := boolSetting(tt.source, tt.key, tt.def)
		if got != tt.expected {
			t.Errorf("boolSetting(%v, %q, %v) = %v, want %v", tt.source, tt.key, tt.def, got, tt.expected)
		}
	}
}

func TestFirstExisting(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "exists.txt")
	dir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Create(file); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		paths    []string
		expected string
	}{
		{[]string{file}, file},
		{[]string{dir}, dir},
		{[]string{filepath.Join(tmpDir, "missing"), file}, file},
		{[]string{filepath.Join(tmpDir, "missing1"), filepath.Join(tmpDir, "missing2")}, ""},
		{[]string{}, ""},
	}
	for _, tt := range tests {
		got := firstExisting(tt.paths...)
		if got != tt.expected {
			t.Errorf("firstExisting(%v) = %q, want %q", tt.paths, got, tt.expected)
		}
	}
}

func TestNestedMap(t *testing.T) {
	source := map[string]any{
		"top": map[string]any{
			"nested": map[string]any{
				"value": "found",
			},
		},
		"plain": "string",
	}

	got := nestedMap(source, "top")
	if got == nil || got["nested"] == nil {
		t.Errorf("nestedMap should return nested map, got %v", got)
	}

	got2 := nestedMap(source, "missing")
	if len(got2) != 0 {
		t.Errorf("nestedMap for missing key should return empty map, got %v", got2)
	}

	got3 := nestedMap(source, "plain")
	if len(got3) != 0 {
		t.Errorf("nestedMap for non-map value should return empty map, got %v", got3)
	}
}

func TestTernaryFunctions(t *testing.T) {
	if ternaryInt(true, 1, 2) != 1 {
		t.Error("ternaryInt(true) should return first value")
	}
	if ternaryInt(false, 1, 2) != 2 {
		t.Error("ternaryInt(false) should return second value")
	}
	if ternaryString(true, "a", "b") != "a" {
		t.Error("ternaryString(true) should return first value")
	}
	if ternaryString(false, "a", "b") != "b" {
		t.Error("ternaryString(false) should return second value")
	}
}

func TestAsString(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{"string", "string"},
		{"", ""},
		{nil, ""},
		{123, "123"},
		{12.5, "12.5"},
		{true, "true"},
		{false, "false"},
		{[]int{1, 2, 3}, "[1 2 3]"},
	}
	for _, tt := range tests {
		got := asString(tt.input)
		if got != tt.expected {
			t.Errorf("asString(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestPlatformNotSupported(t *testing.T) {
	err := platformNotSupported("TestFeature")
	msg := err.Error()
	if !strings.Contains(msg, "TestFeature") {
		t.Errorf("error message should contain feature name, got %q", msg)
	}
	if !strings.Contains(msg, runtime.GOOS) {
		t.Errorf("error message should contain platform, got %q", msg)
	}
}
