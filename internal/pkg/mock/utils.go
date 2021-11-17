package mock

import (
	"fmt"

	"github.com/puppetlabs/prm/pkg/prm"
)

type Utils struct {
	ExpectedPuppetVer   string
	ExpectedBackendType string
}

func (u *Utils) SetAndWriteConfig(k, v string) error {
	if k == prm.PuppetVerCfgKey && v == u.ExpectedPuppetVer || k == prm.BackendCfgKey && v == u.ExpectedBackendType {
		return nil
	}
	return fmt.Errorf(`mock.SetAndWriteConfig(): Unexpected args,
	Expected either:
	- Puppet Version: %s
	- Backend Type: %s
	Got:
	- Key: %s
	- Value: %s`, u.ExpectedPuppetVer, u.ExpectedBackendType, k, v)
}
