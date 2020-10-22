package cli

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newUndeployCmd() *cobra.Command {
	var serverURL string
	var acceptsHTTP bool
	var undeployCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove an app",
		Long:  `Remove an app completely. This will delete all resources associated with the app, including the running container and scaling configuration`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			if len(args) != 1 {
				log.Fatalf("You need to pass the deployment name")
			}

			serverProtocol := "https"
			if acceptsHTTP == true {
				serverProtocol = "http"
			}

			deployURL := fmt.Sprintf("%s://%s/app", serverProtocol, viper.GetViper().GetString("server_url"))
			if serverURL != "" {
				fmt.Printf("Overriding config file server URL for \"%s\"\n", serverURL)
				deployURL = fmt.Sprintf("%s://%s/app", serverProtocol, serverURL)
			}
			fmt.Println("Using server ", deployURL)

			cl := gorequest.New().Delete(deployURL)
			cl.QueryData.Add("name", appName)

			resp, body, errs := cl.Send(nil).End()
			if len(errs) > 0 {
				var result error
				log.Printf("Error deleting: %v", errs)
				return multierror.Append(result, errs...)
			}
			if resp.StatusCode != 200 {
				log.Fatalf("Undeploy failed: %s", body)
			}

			log.Printf("App %s deleted!", appName)
			return nil
		},
	}
	undeployCmd.Flags().StringVarP(
		&serverURL,
		"server-url",
		"s",
		"",
		"The URL to the admin server (without the 'http' prefix)",
	)

	undeployCmd.Flags().BoolVar(
		&acceptsHTTP,
		"use-http",
		false,
		"If set, the server will be called using HTTP instead of HTTPS",
	)

	return undeployCmd

}
