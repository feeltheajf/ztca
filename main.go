package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/feeltheajf/ztca/api"
	"github.com/feeltheajf/ztca/config"
	"github.com/feeltheajf/ztca/dto"
	"github.com/feeltheajf/ztca/log"
	"github.com/feeltheajf/ztca/pki"
)

var cmd = &cobra.Command{
	Use: config.App,
	Run: serve,
}

var flags = struct {
	config string
}{}

func serve(cmd *cobra.Command, args []string) {
	cfg, err := config.Load(flags.config)
	fatal(err)

	fatal(log.Setup(cfg.Log))
	fatal(dto.Setup(cfg.DB))
	fatal(pki.Setup(cfg.PKI))
	fatal(api.Setup(cfg.API))

	fatal(api.Serve())
}

func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	cmd.PersistentFlags().StringVarP(&flags.config, "config", "c", config.File, "path to config file")

	if err := cmd.Execute(); err != nil {
		os.Exit(64)
	}
}
