package cmd

import (
	"erp-job/config"
	"erp-job/logics"
	"erp-job/repository"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves the url shortner service",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	dbs := config.LoadConfig(configPath)

	repos := repository.NewRepository(dbs.SqlitConnection)
	ar := logics.NewAryan(repos)
	fr := logics.NewFararavand(repos, ar)
	fmt.Println(fr)
	e := echo.New()
	e.HideBanner = false

	e.Logger.Fatal(e.Start(":8080"))
}

func init() {
	rootCMD.AddCommand(serveCmd)
}
