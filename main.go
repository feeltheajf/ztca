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

var (
	cfg *config.Config

	cmd = &cobra.Command{
		Use: config.App,
		Run: serve,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			cfg, err = config.Load(flags.config)
			fatal(err)
		},
	}

	flags = struct {
		config string
	}{}
)

func serve(cmd *cobra.Command, args []string) {
	fatal(log.Setup(cfg.Log))
	fatal(dto.Setup(cfg.DB))
	fatal(pki.Setup(cfg.CA))
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
