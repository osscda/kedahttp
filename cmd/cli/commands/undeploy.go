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
		Use:   "rmapp",
		Short: "Remove an app",
		Long:  `Remove an app completely. This will delete all resources associated with the app, including the running container and scaling configuration`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				log.Fatalf("You need to pass the deployment name")
			}
			cl := gorequest.New()

			deployURL := fmt.Sprintf("https://%s/app?name=%s", serverURL, args[0])
			resp, body, errs := cl.Delete(deployURL).Send(nil).End()
			if len(errs) > 0 {
				var result error
				log.Printf("Error deleting: %v", errs)
				return multierror.Append(result, errs...)
			}
			if resp.StatusCode != 200 {
				log.Fatalf("Undeploy failed: %s", body)
			}

			log.Printf("App %s deleted!", args[0])
			return nil
		},
	}
	undeployCmd.Flags().StringVarP(
		&serverURL,
		"server-url",
		"s",
		"admin.wtfcncf.dev",
		"The URL to the admin server (without the 'http' prefix)",
	)
	return undeployCmd

}
