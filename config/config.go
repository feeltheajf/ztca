package config

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/feeltheajf/ztca/api"
	"github.com/feeltheajf/ztca/dto"
	"github.com/feeltheajf/ztca/pki"
	"github.com/feeltheajf/ztca/x/fs"
	"github.com/feeltheajf/ztca/x/log"
)

const (
	App         = "ztca"
	DefaultPath = App + ".yml"
)

var (
	DefaultDir = path.Join(fs.UserConfigDir(), App)
)

func init() {
	err := fs.Mkdir(DefaultDir)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	API *api.Config `yaml:"api"`
	DB  *dto.Config `yaml:"db"`
	Log *log.Config `yaml:"log"`
	PKI *pki.Config `yaml:"pki"`
}

func Load(path string) (*Config, error) {
	cfg := new(Config)
	return cfg, load(path, cfg)
}

func load(path string, i interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	b = []byte(os.ExpandEnv(string(b)))
	if err := yaml.UnmarshalStrict(b, i); err != nil {
		return err
	}

	return nil
}

func Path(filename string) string {
	return path.Join(DefaultDir, filename)
}
