package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/feeltheajf/ztca/api"
	"github.com/feeltheajf/ztca/dto"
	"github.com/feeltheajf/ztca/fs"
	"github.com/feeltheajf/ztca/log"
	"github.com/feeltheajf/ztca/pki"
)

const (
	App  = "ztca"
	File = App + ".yml"
)

type Config struct {
	API *api.Config `yaml:"api"`
	DB  *dto.Config `yaml:"db"`
	Log *log.Config `yaml:"log"`
	CA  *pki.Config `yaml:"ca"`
}

func Load(path string) (*Config, error) {
	raw, err := fs.Read(path)
	if err != nil {
		return nil, err
	}

	cfg := defaultConfig()
	if err := yaml.UnmarshalStrict([]byte(os.ExpandEnv(raw)), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		//
	}
}
