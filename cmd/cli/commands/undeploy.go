package cli

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
)

func newUndeployCmd() *cobra.Command {
	var serverURL string
	var undeployCmd = &cobra.Command{
		Use:   "undeploy",
		Short: "Remove a deployment",
		Long:  `Remove a deployment by its name. This will delete all traces of the images running and immediately stop serving your app.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				log.Fatalf("You need to pass the deployment name")
			}
			cl := gorequest.New()

			deployURL := fmt.Sprintf("https://%s/admin/deploy?name=%s", serverURL, args[0])
			resp, body, errs := cl.Delete(deployURL).Send(nil).End()
			if len(errs) > 0 {
				var result error
				log.Printf("Error undeploying: %v", errs)
				return multierror.Append(result, errs...)
			}
			if resp.StatusCode != 200 {
				log.Fatalf("Undeploy failed: %s", body)
			}

			log.Printf("Undeployed!")
			return nil
		},
	}
	undeployCmd.Flags().StringVarP(
		&serverURL,
		"server-url",
		"s",
		"wtfcncf.dev",
		"The URL to the admin server (without the 'http' prefix",
	)
	return undeployCmd

}
