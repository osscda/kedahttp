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
	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the image to the environment.",
		Long:  `Deploy the image to the environment.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cl := gorequest.New()

			deployURL := fmt.Sprintf("https://%s/admin/deploy", serverURL)
			_, _, errs := cl.Post(deployURL).Send(map[string]string{
				"name":  name,
				"image": deployImage,
			}).End()
			if len(errs) > 0 {
				var result error
				log.Printf("Error deploying: %v", errs)
				return multierror.Append(result, errs...)
			}

			log.Printf("Deployed image: %s", args[0])
			return nil
		},
	}

	flags := deployCmd.Flags()

	flags.StringVarP(
		&serverURL,
		"server-url",
		"s",
		"wtfcncf.dev",
		"The URL to the admin server (without the 'http' prefix",
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
