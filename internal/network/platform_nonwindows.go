//go:build !windows

package network

import "os/exec"

func applyPlatformPingAttrs(cmd *exec.Cmd) {
	_ = cmd
}
