package config

import (
	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

// Config -- db
type Config struct {
	ReqQueueSize    int32
	RespQueueSize   int32
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int64
}

// Cfg -- config var
var Cfg Config

// LoadConfigFromFile load RelayConstConfig from toml
func LoadConfigFromFile(filename string) error {
	if _, err := toml.DecodeFile(filename, &Cfg); err != nil {
		glog.Errorln(err)
		return err
	} else {
		return nil
	}
}
