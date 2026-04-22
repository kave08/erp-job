package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configPath string

var rootCMD = &cobra.Command{
	Use:          "erp-job",
	Short:        "ETL between erps service!",
	SilenceUsage: true,
}

func init() {
	rootCMD.PersistentFlags().StringVarP(&configPath, "config-path", "c", "env.yml", "path to config directory")

}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
