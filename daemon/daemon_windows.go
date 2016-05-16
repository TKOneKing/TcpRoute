// +build windows

package Daemon

import (
    "io"
)

func MakeDaemon() (io.Reader, io.Reader, error) {
    panic("windows 不支持 daemon。")
}