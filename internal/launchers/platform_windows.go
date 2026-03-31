//go:build windows

package launchers

import (
	"os/exec"
	"syscall"
)

func applyPlatformProcessAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
