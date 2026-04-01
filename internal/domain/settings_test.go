package domain

import (
	"testing"
)

func TestDefaultSettings(t *testing.T) {
	s := DefaultSettings()

	if s.Rdp.Resolution != "1920x1080" {
		t.Errorf("Rdp.Resolution = %q, want 1920x1080", s.Rdp.Resolution)
	}
	if s.Rdp.ColorDepth != "32" {
		t.Errorf("Rdp.ColorDepth = %q, want 32", s.Rdp.ColorDepth)
	}
	if s.Rdp.Clipboard != true {
		t.Error("Rdp.Clipboard should be true")
	}
	if s.Rdp.PromptCredentials != true {
		t.Error("Rdp.PromptCredentials should be true")
	}
	if s.Horizon.LoginAsCurrentUser != true {
		t.Error("Horizon.LoginAsCurrentUser should be true")
	}
	if s.Citrix.StoreAlreadyConfigured != false {
		t.Error("Citrix.StoreAlreadyConfigured should be false")
	}
	if s.Updates.UseGithub != true {
		t.Error("Updates.UseGithub should be true")
	}
	if s.NetworkCheck.LatencyThresholdMs != 100 {
		t.Errorf("NetworkCheck.LatencyThresholdMs = %d, want 100", s.NetworkCheck.LatencyThresholdMs)
	}
}

func TestSettingsToMap(t *testing.T) {
	s := DefaultSettings()
	m := s.ToMap()

	if _, ok := m["user"]; !ok {
		t.Error("ToMap missing user section")
	}
	if _, ok := m["rdp"]; !ok {
		t.Error("ToMap missing rdp section")
	}
	if _, ok := m["horizon"]; !ok {
		t.Error("ToMap missing horizon section")
	}
	if _, ok := m["citrix"]; !ok {
		t.Error("ToMap missing citrix section")
	}
	if _, ok := m["updates"]; !ok {
		t.Error("ToMap missing updates section")
	}

	citrix := m["citrix"].(map[string]any)
	if _, ok := citrix["storeAlreadyConfigured"]; !ok {
		t.Error("ToMap missing citrix.storeAlreadyConfigured")
	}
	if citrix["storeAlreadyConfigured"] != false {
		t.Errorf("citrix.storeAlreadyConfigured = %v, want false", citrix["storeAlreadyConfigured"])
	}
}

func TestSettingsFromMap(t *testing.T) {
	m := map[string]any{
		"user": map[string]any{
			"domain":   "MOSCOW",
			"username": "ivanov",
		},
		"rdp": map[string]any{
			"resolution": "1366x768",
			"clipboard":  false,
		},
		"citrix": map[string]any{
			"storeAlreadyConfigured": true,
			"resourceName":           "Desktop",
		},
		"updates": map[string]any{
			"useGithub": false,
		},
	}

	s := SettingsFromMap(m)

	if s.User.Domain != "MOSCOW" {
		t.Errorf("User.Domain = %q, want MOSCOW", s.User.Domain)
	}
	if s.User.Username != "ivanov" {
		t.Errorf("User.Username = %q, want ivanov", s.User.Username)
	}
	if s.Rdp.Resolution != "1366x768" {
		t.Errorf("Rdp.Resolution = %q, want 1366x768", s.Rdp.Resolution)
	}
	if s.Rdp.Clipboard != false {
		t.Error("Rdp.Clipboard should be false")
	}
	if s.Citrix.StoreAlreadyConfigured != true {
		t.Error("Citrix.StoreAlreadyConfigured should be true")
	}
	if s.Citrix.ResourceName != "Desktop" {
		t.Errorf("Citrix.ResourceName = %q, want Desktop", s.Citrix.ResourceName)
	}
	if s.Updates.UseGithub != false {
		t.Error("Updates.UseGithub should be false")
	}
}

func TestSettingsRoundTrip(t *testing.T) {
	original := DefaultSettings()
	original.User.Domain = "REGIONS"
	original.User.Username = "petrov"
	original.Rdp.Resolution = "1366x768"
	original.Rdp.CustomFlags = "-test flag"
	original.Horizon.AppName = "App1"
	original.Citrix.StoreAlreadyConfigured = true
	original.General.MinimizeToTray = true
	original.NetworkCheck.LatencyThresholdMs = 50
	original.Updates.UseGithub = false

	m := original.ToMap()
	restored := SettingsFromMap(m)

	if restored.User.Domain != original.User.Domain {
		t.Errorf("User.Domain round-trip: got %q, want %q", restored.User.Domain, original.User.Domain)
	}
	if restored.User.Username != original.User.Username {
		t.Errorf("User.Username round-trip: got %q, want %q", restored.User.Username, original.User.Username)
	}
	if restored.Rdp.Resolution != original.Rdp.Resolution {
		t.Errorf("Rdp.Resolution round-trip: got %q, want %q", restored.Rdp.Resolution, original.Rdp.Resolution)
	}
	if restored.Rdp.CustomFlags != original.Rdp.CustomFlags {
		t.Errorf("Rdp.CustomFlags round-trip: got %q, want %q", restored.Rdp.CustomFlags, original.Rdp.CustomFlags)
	}
	if restored.Horizon.AppName != original.Horizon.AppName {
		t.Errorf("Horizon.AppName round-trip: got %q, want %q", restored.Horizon.AppName, original.Horizon.AppName)
	}
	if restored.Citrix.StoreAlreadyConfigured != original.Citrix.StoreAlreadyConfigured {
		t.Errorf("Citrix.StoreAlreadyConfigured round-trip: got %v, want %v", restored.Citrix.StoreAlreadyConfigured, original.Citrix.StoreAlreadyConfigured)
	}
	if restored.General.MinimizeToTray != original.General.MinimizeToTray {
		t.Errorf("General.MinimizeToTray round-trip: got %v, want %v", restored.General.MinimizeToTray, original.General.MinimizeToTray)
	}
	if restored.NetworkCheck.LatencyThresholdMs != original.NetworkCheck.LatencyThresholdMs {
		t.Errorf("NetworkCheck.LatencyThresholdMs round-trip: got %d, want %d", restored.NetworkCheck.LatencyThresholdMs, original.NetworkCheck.LatencyThresholdMs)
	}
	if restored.Updates.UseGithub != original.Updates.UseGithub {
		t.Errorf("Updates.UseGithub round-trip: got %v, want %v", restored.Updates.UseGithub, original.Updates.UseGithub)
	}
}

func TestSettingsFromMapPartial(t *testing.T) {
	m := map[string]any{
		"user": map[string]any{
			"username": "ivanov",
		},
	}

	s := SettingsFromMap(m)

	if s.User.Username != "ivanov" {
		t.Errorf("User.Username = %q, want ivanov", s.User.Username)
	}
	if s.User.Domain != "" {
		t.Errorf("User.Domain = %q, want empty (default)", s.User.Domain)
	}
	if s.Rdp.Resolution != "1920x1080" {
		t.Errorf("Rdp.Resolution should use default, got %q", s.Rdp.Resolution)
	}
}
