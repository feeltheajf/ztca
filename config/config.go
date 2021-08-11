package config

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/feeltheajf/ztca/api"
	"github.com/feeltheajf/ztca/dto"
	"github.com/feeltheajf/ztca/log"
	"github.com/feeltheajf/ztca/pki"
)

const (
	App  = "ztca"
	File = App + ".yml"

	permissionsFile      = 0600
	permissionsDirectory = 0700
)

var (
	root string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	root = path.Join(home, App)
	if err := Mkdir(root); err != nil {
		panic(err)
	}
}

func Path(elem ...string) string {
	return path.Join(append([]string{root}, elem...)...)
}

func Mkdir(dir string) error {
	return os.MkdirAll(dir, permissionsDirectory)
}

func Write(file string, data []byte) error {
	return ioutil.WriteFile(file, data, permissionsFile)
}

type Config struct {
	API *api.Config `yaml:"api"`
	DB  *dto.Config `yaml:"db"`
	Log *log.Config `yaml:"log"`
	PKI *pki.Config `yaml:"pki"`
}

func Load(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	b = []byte(os.ExpandEnv(string(b)))
	if err := yaml.UnmarshalStrict(b, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
