package pod

import (
	"git.parallelcoin.io/pod/lib/clog"
)

var log = clog.NewSubSystem("Pod", clog.Ninf)

func print(fmt string, items ...interface{}) clog.Fmt {
	return clog.Fmt{Fmt: fmt, Items: items}
}
