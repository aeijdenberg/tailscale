// Copyright (c) 2022 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tailssh

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
	"tailscale.com/types/logger"
)

func init() {
	maybeRunCommandWithLogin = maybeRunCommandWithLoginDarwin
}

func (ia *incubatorArgs) loginArgs() []string {
	return []string{ia.loginCmdPath, "-fp", "-h", ia.remoteIP, ia.localUser}
}

func setGroups(groupIDs []int) error {
	// darwin returns "invalid argument" if more than 16 groups are passed to syscall.Setgroups
	// some info can be found here:
	// https://opensource.apple.com/source/samba/samba-187.8/patches/support-darwin-initgroups-syscall.auto.html
	// this fix isn't great, as anyone reading this has probably just wasted hours figuring out why
	// some permissions thing isn't working, due to some arbitrary group ordering, but it at least allows
	// this to work for more things than it previously did.
	return syscall.Setgroups(groupIDs[:16])
}

// maybeRunCommandWithLoginDarwin is the darwin implementation of maybeRunCommandWithLogin
func maybeRunCommandWithLoginDarwin(logf logger.Logf, ia incubatorArgs) (bool, error) {
	if ia.isSFTP || ia.isShell || ia.cmdName == "" || ia.loginCmdPath == "" {
		return false, nil
	}

	// if we run it this way on Mac, this seems to solve various user ID issues
	return true, unix.Exec(ia.loginCmdPath, append(
		[]string{ia.loginCmdPath, "-f", "-lpq", "-h", ia.remoteIP, ia.localUser, ia.cmdName},
		ia.cmdArgs...), os.Environ())
}
