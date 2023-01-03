package blipUtil

import (
	"time"
)

type IBlipMonitor interface {
	RenderComplete(escaper IBlipEscaper, name string, langType string, duration time.Duration, err error)
}
