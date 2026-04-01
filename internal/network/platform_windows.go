//go:build windows

package network

import (
	"os/exec"
	"syscall"
)

func applyPlatformPingAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
