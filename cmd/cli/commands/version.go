package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cscaler",
	Long:  `All software has versions. This is cscaler's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cscaler version v0.1 -- HEAD")
	},
}
