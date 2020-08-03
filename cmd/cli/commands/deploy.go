package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var platform string
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the image to the environment.",
	Long:  `Deploy the image to the environment.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := requestDeployment(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("deployed image: %s\n", args[0])
		return nil
	},
}

func init() {
	deployCmd.Flags().StringVarP(&platform, "platform", "p", "ACI", "which platform the image is deployed to")
	rootCmd.AddCommand(deployCmd)
}

func requestDeployment(image string) error {
	if image == "" {
		return fmt.Errorf("invalid image name")
	}
	if platform != "ACI" && platform != "VMSS" {
		return fmt.Errorf("invalid platform ACI or VMSS only")
	}

	return nil
}
