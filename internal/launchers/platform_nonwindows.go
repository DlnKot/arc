//go:build !windows

package launchers

import "os/exec"

func applyPlatformProcessAttrs(cmd *exec.Cmd) {
	_ = cmd
}
