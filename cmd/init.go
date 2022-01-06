package cmd

import (
	"github.com/spf13/cobra"
	"github.com/victorien-a/gen-model/handler"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate local yaml config",
	Long:  `manage vars into file`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		handler.GenConfigYaml(&CmdRequest)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
