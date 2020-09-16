package cli

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
)

func newDeployCmd() *cobra.Command {
	var serverURL string
	var deployImage string
	var port string
	var deployCmd = &cobra.Command{
		Use:   "run",
		Short: "Create a new app",
		Long:  `Start serving & scaling the given container.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cl := gorequest.New()

			deployURL := fmt.Sprintf("https://%s/app", serverURL)
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
				log.Fatalf("Create failed: %s", body)
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
		"admin.wtfcncf.dev",
		"The URL to the admin server (without the 'http' prefix)",
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
