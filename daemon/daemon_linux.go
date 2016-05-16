// +build darwin freebsd linux

package Daemon

import (
    "io"
    "github.com/VividCortex/godaemon"
)

func MakeDaemon() (io.Reader, io.Reader, error) {
    return godaemon.MakeDaemon(&godaemon.DaemonAttr{})
}