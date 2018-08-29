// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"
	"github.com/elastic/beats/libbeat/common"
)

type Config struct {
	Modules []*common.Config `config:"modules"`
	Period time.Duration `config:"period"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}
