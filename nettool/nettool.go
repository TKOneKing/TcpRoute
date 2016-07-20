package nettool

import (
	"fmt"
)

type SetLingerer interface {
	SetLinger(sec int) error
}
type SetNoDelayer interface {
	SetNoDelay(noDelay bool) error
}

type SetReadBufferer interface {
	SetReadBuffer(bytes int) error
}
type SetWriteBufferer interface {
	SetWriteBuffer(bytes int) error
}
type SetReadBuffer interface {
	SetReadBuffer(bytes int) error
}

//SetDeadline 包含 net.Conn

// 开关 Delay 算法
// noDelay = true 关闭
// false 性能更好，但是延迟高
func SetNoDelay(conn interface{}, noDelay bool) error {
	ccd, _ := conn.(SetNoDelayer)
	if ccd == nil {
		return fmt.Errorf("conn 未提供 SetNoDelay 方法。")
	}
	ccd.SetNoDelay(noDelay)
	return nil
}

func SetLinger(c interface{}, sec int) error {
	ccd, _ := c.(SetLingerer)
	if ccd == nil {
		return fmt.Errorf("conn 未提供 SetNoDelay 方法。")
	}
	ccd.SetLinger(sec)
	return nil

}
