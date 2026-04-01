package domain

type User struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
}

type RdpAudio struct {
	Playback bool `json:"playback"`
	Capture  bool `json:"capture"`
}

type RdpRedirect struct {
	Printers   bool `json:"printers"`
	Smartcards bool `json:"smartcards"`
	Webauthn   bool `json:"webauthn"`
}

type RdpPerformance struct {
	Wallpaper          bool `json:"wallpaper"`
	FontSmoothing      bool `json:"fontSmoothing"`
	DesktopComposition bool `json:"desktopComposition"`
	FullWindowDrag     bool `json:"fullWindowDrag"`
	MenuAnimations     bool `json:"menuAnimations"`
}

type RdpSettings struct {
	Resolution        string         `json:"resolution"`
	ColorDepth        string         `json:"colorDepth"`
	Multimon          bool           `json:"multimon"`
	Span              bool           `json:"span"`
	Clipboard         bool           `json:"clipboard"`
	DriveMapping      bool           `json:"driveMapping"`
	UseAdminSession   bool           `json:"useAdminSession"`
	PromptCredentials bool           `json:"promptCredentials"`
	StartFullScreen   bool           `json:"startFullScreen"`
	Audio             RdpAudio       `json:"audio"`
	Redirect          RdpRedirect    `json:"redirect"`
	Performance       RdpPerformance `json:"performance"`
	CustomFlags       string         `json:"customFlags"`
}

type HorizonSettings struct {
	AppName                      string `json:"appName"`
	DesktopProtocol              string `json:"desktopProtocol"`
	DesktopLayout                string `json:"desktopLayout"`
	Monitors                     string `json:"monitors"`
	Unattended                   bool   `json:"unattended"`
	NonInteractive               bool   `json:"nonInteractive"`
	LaunchMinimized              bool   `json:"launchMinimized"`
	LoginAsCurrentUser           bool   `json:"loginAsCurrentUser"`
	HideClientAfterLaunchSession bool   `json:"hideClientAfterLaunchSession"`
	UseExisting                  bool   `json:"useExisting"`
	SingleAutoConnect            bool   `json:"singleAutoConnect"`
	CustomPath                   string `json:"customPath"`
	CustomFlags                  string `json:"customFlags"`
}

type CitrixSettings struct {
	AccountName            string `json:"accountName"`
	ResourceName           string `json:"resourceName"`
	StoreAlreadyConfigured bool   `json:"storeAlreadyConfigured"`
	CustomPath             string `json:"customPath"`
	CustomFlags            string `json:"customFlags"`
}

type GeneralSettings struct {
	MinimizeToTray bool `json:"minimizeToTray"`
	StartMinimized bool `json:"startMinimized"`
}

type NetworkCheckSettings struct {
	LatencyThresholdMs int `json:"latencyThresholdMs"`
}

type UpdatesSettings struct {
	UseGithub         bool   `json:"useGithub"`
	InternalServerURL string `json:"internalServerUrl"`
	AutoCheck         bool   `json:"autoCheck"`
	InstallOnQuit     bool   `json:"installOnQuit"`
}

type Settings struct {
	User         User                 `json:"user"`
	Rdp          RdpSettings          `json:"rdp"`
	Horizon      HorizonSettings      `json:"horizon"`
	Citrix       CitrixSettings       `json:"citrix"`
	General      GeneralSettings      `json:"general"`
	NetworkCheck NetworkCheckSettings `json:"networkCheck"`
	Updates      UpdatesSettings      `json:"updates"`
}

func DefaultSettings() Settings {
	return Settings{
		User: User{},
		Rdp: RdpSettings{
			Resolution:        "1920x1080",
			ColorDepth:        "32",
			Multimon:          false,
			Span:              false,
			Clipboard:         true,
			DriveMapping:      false,
			UseAdminSession:   false,
			PromptCredentials: true,
			StartFullScreen:   false,
			Audio:             RdpAudio{Playback: true, Capture: false},
			Redirect:          RdpRedirect{Printers: true, Smartcards: true, Webauthn: true},
			Performance:       RdpPerformance{Wallpaper: true, FontSmoothing: true, DesktopComposition: true, FullWindowDrag: true, MenuAnimations: true},
			CustomFlags:       "",
		},
		Horizon: HorizonSettings{
			AppName:                      "",
			DesktopProtocol:              "",
			DesktopLayout:                "",
			Monitors:                     "",
			Unattended:                   false,
			NonInteractive:               false,
			LaunchMinimized:              false,
			LoginAsCurrentUser:           true,
			HideClientAfterLaunchSession: false,
			UseExisting:                  false,
			SingleAutoConnect:            false,
			CustomPath:                   "",
			CustomFlags:                  "",
		},
		Citrix: CitrixSettings{
			AccountName:            "",
			ResourceName:           "",
			StoreAlreadyConfigured: false,
			CustomPath:             "",
			CustomFlags:            "",
		},
		General: GeneralSettings{
			MinimizeToTray: false,
			StartMinimized: false,
		},
		NetworkCheck: NetworkCheckSettings{
			LatencyThresholdMs: 100,
		},
		Updates: UpdatesSettings{
			UseGithub:         true,
			InternalServerURL: "http://10.230.121.212",
			AutoCheck:         true,
			InstallOnQuit:     false,
		},
	}
}

func (s Settings) ToMap() map[string]any {
	return map[string]any{
		"user": map[string]any{
			"domain":   s.User.Domain,
			"username": s.User.Username,
		},
		"rdp": map[string]any{
			"resolution":        s.Rdp.Resolution,
			"colorDepth":        s.Rdp.ColorDepth,
			"multimon":          s.Rdp.Multimon,
			"span":              s.Rdp.Span,
			"clipboard":         s.Rdp.Clipboard,
			"driveMapping":      s.Rdp.DriveMapping,
			"useAdminSession":   s.Rdp.UseAdminSession,
			"promptCredentials": s.Rdp.PromptCredentials,
			"startFullScreen":   s.Rdp.StartFullScreen,
			"audio": map[string]any{
				"playback": s.Rdp.Audio.Playback,
				"capture":  s.Rdp.Audio.Capture,
			},
			"redirect": map[string]any{
				"printers":   s.Rdp.Redirect.Printers,
				"smartcards": s.Rdp.Redirect.Smartcards,
				"webauthn":   s.Rdp.Redirect.Webauthn,
			},
			"performance": map[string]any{
				"wallpaper":          s.Rdp.Performance.Wallpaper,
				"fontSmoothing":      s.Rdp.Performance.FontSmoothing,
				"desktopComposition": s.Rdp.Performance.DesktopComposition,
				"fullWindowDrag":     s.Rdp.Performance.FullWindowDrag,
				"menuAnimations":     s.Rdp.Performance.MenuAnimations,
			},
			"customFlags": s.Rdp.CustomFlags,
		},
		"horizon": map[string]any{
			"appName":                      s.Horizon.AppName,
			"desktopProtocol":              s.Horizon.DesktopProtocol,
			"desktopLayout":                s.Horizon.DesktopLayout,
			"monitors":                     s.Horizon.Monitors,
			"unattended":                   s.Horizon.Unattended,
			"nonInteractive":               s.Horizon.NonInteractive,
			"launchMinimized":              s.Horizon.LaunchMinimized,
			"loginAsCurrentUser":           s.Horizon.LoginAsCurrentUser,
			"hideClientAfterLaunchSession": s.Horizon.HideClientAfterLaunchSession,
			"useExisting":                  s.Horizon.UseExisting,
			"singleAutoConnect":            s.Horizon.SingleAutoConnect,
			"customPath":                   s.Horizon.CustomPath,
			"customFlags":                  s.Horizon.CustomFlags,
		},
		"citrix": map[string]any{
			"accountName":            s.Citrix.AccountName,
			"resourceName":           s.Citrix.ResourceName,
			"storeAlreadyConfigured": s.Citrix.StoreAlreadyConfigured,
			"customPath":             s.Citrix.CustomPath,
			"customFlags":            s.Citrix.CustomFlags,
		},
		"general": map[string]any{
			"minimizeToTray": s.General.MinimizeToTray,
			"startMinimized": s.General.StartMinimized,
		},
		"networkCheck": map[string]any{
			"latencyThresholdMs": s.NetworkCheck.LatencyThresholdMs,
		},
		"updates": map[string]any{
			"useGithub":         s.Updates.UseGithub,
			"internalServerUrl": s.Updates.InternalServerURL,
			"autoCheck":         s.Updates.AutoCheck,
			"installOnQuit":     s.Updates.InstallOnQuit,
		},
	}
}

func SettingsFromMap(m map[string]any) Settings {
	s := DefaultSettings()

	if user, ok := m["user"].(map[string]any); ok {
		if v, ok := user["domain"].(string); ok {
			s.User.Domain = v
		}
		if v, ok := user["username"].(string); ok {
			s.User.Username = v
		}
	}

	if rdp, ok := m["rdp"].(map[string]any); ok {
		if v, ok := rdp["resolution"].(string); ok {
			s.Rdp.Resolution = v
		}
		if v, ok := rdp["colorDepth"].(string); ok {
			s.Rdp.ColorDepth = v
		}
		if v, ok := rdp["multimon"].(bool); ok {
			s.Rdp.Multimon = v
		}
		if v, ok := rdp["span"].(bool); ok {
			s.Rdp.Span = v
		}
		if v, ok := rdp["clipboard"].(bool); ok {
			s.Rdp.Clipboard = v
		}
		if v, ok := rdp["driveMapping"].(bool); ok {
			s.Rdp.DriveMapping = v
		}
		if v, ok := rdp["useAdminSession"].(bool); ok {
			s.Rdp.UseAdminSession = v
		}
		if v, ok := rdp["promptCredentials"].(bool); ok {
			s.Rdp.PromptCredentials = v
		}
		if v, ok := rdp["startFullScreen"].(bool); ok {
			s.Rdp.StartFullScreen = v
		}
		if v, ok := rdp["customFlags"].(string); ok {
			s.Rdp.CustomFlags = v
		}
		if audio, ok := rdp["audio"].(map[string]any); ok {
			if v, ok := audio["playback"].(bool); ok {
				s.Rdp.Audio.Playback = v
			}
			if v, ok := audio["capture"].(bool); ok {
				s.Rdp.Audio.Capture = v
			}
		}
		if redirect, ok := rdp["redirect"].(map[string]any); ok {
			if v, ok := redirect["printers"].(bool); ok {
				s.Rdp.Redirect.Printers = v
			}
			if v, ok := redirect["smartcards"].(bool); ok {
				s.Rdp.Redirect.Smartcards = v
			}
			if v, ok := redirect["webauthn"].(bool); ok {
				s.Rdp.Redirect.Webauthn = v
			}
		}
		if perf, ok := rdp["performance"].(map[string]any); ok {
			if v, ok := perf["wallpaper"].(bool); ok {
				s.Rdp.Performance.Wallpaper = v
			}
			if v, ok := perf["fontSmoothing"].(bool); ok {
				s.Rdp.Performance.FontSmoothing = v
			}
			if v, ok := perf["desktopComposition"].(bool); ok {
				s.Rdp.Performance.DesktopComposition = v
			}
			if v, ok := perf["fullWindowDrag"].(bool); ok {
				s.Rdp.Performance.FullWindowDrag = v
			}
			if v, ok := perf["menuAnimations"].(bool); ok {
				s.Rdp.Performance.MenuAnimations = v
			}
		}
	}

	if horizon, ok := m["horizon"].(map[string]any); ok {
		if v, ok := horizon["appName"].(string); ok {
			s.Horizon.AppName = v
		}
		if v, ok := horizon["desktopProtocol"].(string); ok {
			s.Horizon.DesktopProtocol = v
		}
		if v, ok := horizon["desktopLayout"].(string); ok {
			s.Horizon.DesktopLayout = v
		}
		if v, ok := horizon["monitors"].(string); ok {
			s.Horizon.Monitors = v
		}
		if v, ok := horizon["unattended"].(bool); ok {
			s.Horizon.Unattended = v
		}
		if v, ok := horizon["nonInteractive"].(bool); ok {
			s.Horizon.NonInteractive = v
		}
		if v, ok := horizon["launchMinimized"].(bool); ok {
			s.Horizon.LaunchMinimized = v
		}
		if v, ok := horizon["loginAsCurrentUser"].(bool); ok {
			s.Horizon.LoginAsCurrentUser = v
		}
		if v, ok := horizon["hideClientAfterLaunchSession"].(bool); ok {
			s.Horizon.HideClientAfterLaunchSession = v
		}
		if v, ok := horizon["useExisting"].(bool); ok {
			s.Horizon.UseExisting = v
		}
		if v, ok := horizon["singleAutoConnect"].(bool); ok {
			s.Horizon.SingleAutoConnect = v
		}
		if v, ok := horizon["customPath"].(string); ok {
			s.Horizon.CustomPath = v
		}
		if v, ok := horizon["customFlags"].(string); ok {
			s.Horizon.CustomFlags = v
		}
	}

	if citrix, ok := m["citrix"].(map[string]any); ok {
		if v, ok := citrix["accountName"].(string); ok {
			s.Citrix.AccountName = v
		}
		if v, ok := citrix["resourceName"].(string); ok {
			s.Citrix.ResourceName = v
		}
		if v, ok := citrix["storeAlreadyConfigured"].(bool); ok {
			s.Citrix.StoreAlreadyConfigured = v
		}
		if v, ok := citrix["customPath"].(string); ok {
			s.Citrix.CustomPath = v
		}
		if v, ok := citrix["customFlags"].(string); ok {
			s.Citrix.CustomFlags = v
		}
	}

	if general, ok := m["general"].(map[string]any); ok {
		if v, ok := general["minimizeToTray"].(bool); ok {
			s.General.MinimizeToTray = v
		}
		if v, ok := general["startMinimized"].(bool); ok {
			s.General.StartMinimized = v
		}
	}

	if nc, ok := m["networkCheck"].(map[string]any); ok {
		switch v := nc["latencyThresholdMs"].(type) {
		case float64:
			s.NetworkCheck.LatencyThresholdMs = int(v)
		case int:
			s.NetworkCheck.LatencyThresholdMs = v
		case int64:
			s.NetworkCheck.LatencyThresholdMs = int(v)
		}
	}

	if updates, ok := m["updates"].(map[string]any); ok {
		if v, ok := updates["useGithub"].(bool); ok {
			s.Updates.UseGithub = v
		}
		if v, ok := updates["internalServerUrl"].(string); ok {
			s.Updates.InternalServerURL = v
		}
		if v, ok := updates["autoCheck"].(bool); ok {
			s.Updates.AutoCheck = v
		}
		if v, ok := updates["installOnQuit"].(bool); ok {
			s.Updates.InstallOnQuit = v
		}
	}

	return s
}
