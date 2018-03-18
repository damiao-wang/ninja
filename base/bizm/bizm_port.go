package bizm

import (
	"fmt"

	"ninja/base/mconf"
	"gitlab.1dmy.com/ezbuy/base/misc/errors"
)

type Portor interface {
	SetPort(string) error
}

func SetPortByName(pkgName string, s Portor) error {
	data := mconf.GetListenAddrs()
	serviceKey := fmt.Sprintf("service %v", pkgName)
	if pkgName == "main" {
		serviceKey = ""
	}
	if v, ok := data[serviceKey]; ok {
		return s.SetPort(v)
	}
	return errors.Fmt(`port of "service:%v" is not found`, pkgName)

}