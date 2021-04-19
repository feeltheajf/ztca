package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/feeltheajf/ztca/api"
	"github.com/feeltheajf/ztca/config"
	"github.com/feeltheajf/ztca/dto"
	"github.com/feeltheajf/ztca/pki"
	"github.com/feeltheajf/ztca/x/log"
)

var cmd = &cobra.Command{
	Use: "ztca",
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
	cmd.PersistentFlags().StringVarP(&flags.config, "config", "c", "ztca.yml", "path to config file")

	if err := cmd.Execute(); err != nil {
		os.Exit(64)
	}
}
