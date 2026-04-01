package launchers

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func launchVpn() error {
	if runtime.GOOS == "darwin" {
		bundle := firstExisting(
			"/Applications/Endpoint Security VPN.app",
			"/Applications/EndpointConnect.app",
			"/Applications/Endpoint Connect.app",
		)
		if bundle != "" {
			return startDetached("open", bundle)
		}
		return startDetached("open", "-a", "Endpoint Security VPN")
	}

	if runtime.GOOS == "windows" {
		programFiles := os.Getenv("ProgramFiles")
		programFiles86 := os.Getenv("ProgramFiles(x86)")
		programW6432 := os.Getenv("ProgramW6432")
		bases := []string{programFiles, programFiles86, programW6432}

		var candidates []string
		for _, base := range bases {
			if base == "" {
				continue
			}
			candidates = append(candidates,
				filepath.Join(base, "TrGUI", "TrGUI.exe"),
				filepath.Join(base, "TrGUI", "TrGUI", "TrGUI.exe"),
				filepath.Join(base, "CheckPoint", "Endpoint Connect", "TrGUI.exe"),
				filepath.Join(base, "CheckPoint", "Endpoint Security VPN", "TrGUI.exe"),
				filepath.Join(base, "Check Point", "Endpoint Connect", "TrGUI.exe"),
				filepath.Join(base, "Check Point", "Endpoint Security VPN", "TrGUI.exe"),
				filepath.Join(programFiles86, "CheckPoint", "Endpoint Connect", "TrGUI.exe"),
				filepath.Join(programFiles86, "CheckPoint", "Endpoint Security VPN", "TrGUI.exe"),
			)
		}

		exe := firstExisting(candidates...)
		if exe == "" {
			if resolved, err := exec.LookPath("TrGUI.exe"); err == nil {
				exe = resolved
			}
		}
		if exe == "" {
			return errors.New("VPN client not found: TrGUI.exe")
		}
		return startDetached(exe)
	}

	return platformNotSupported("VPN launcher")
}
