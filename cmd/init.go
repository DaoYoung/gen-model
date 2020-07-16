package cmd

import (
	"github.com/DaoYoung/gen-model/handler"
	"github.com/spf13/cobra"
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
