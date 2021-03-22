package cmd

import (
	"rest/config/srvs"
	"rest/pkg/gcore"
	"rest/insecure"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches the example webserver on  "+demoAddr,
	Run: func(cmd *cobra.Command, args []string) {
		gcore.MakeInsecure(insecure.Key,insecure.Cert)
		gcore.RunServe(srvs.ServerMap,demoAddr,port)
		//serve()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
