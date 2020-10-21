package cli

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newDeployCmd() *cobra.Command {
	var serverURL string
	var deployImage string
	var port string
	var acceptsHTTP bool
	var deployCmd = &cobra.Command{
		Use:   "run",
		Short: "Create a new app",
		Long:  `Start serving & scaling the given container.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cl := gorequest.New()

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

			resp, body, errs := cl.Post(deployURL).Send(map[string]string{
				"name":  name,
				"image": deployImage,
				"port":  port,
			}).End()

			if len(errs) > 0 {
				var result error
				log.Printf("Error creating: %v", errs)
				return multierror.Append(result, errs...)
			}

			if resp.StatusCode != 200 {
				log.Fatalf("Create failed: \"%s\" with status (%d)", body, resp.StatusCode)
			}

			log.Printf("Created %s (image %s)", name, deployImage)
			return nil
		},
	}

	flags := deployCmd.Flags()

	flags.StringVarP(
		&serverURL,
		"server-url",
		"s",
		"",
		"The URL to the admin server (without the 'http' prefix)",
	)

	flags.BoolVar(
		&acceptsHTTP,
		"use-http",
		false,
		"If set, the server will be called using HTTP instead of HTTPS",
	)

	flags.StringVarP(
		&port,
		"port",
		"p",
		"8080",
		"The port that the container will be listening on",
	)

	flags.StringVarP(
		&deployImage,
		"image",
		"i",
		"",
		"The container image to deploy",
	)

	return deployCmd

}
