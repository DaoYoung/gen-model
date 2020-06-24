/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/DaoYoung/gen-model/handler"
	"github.com/spf13/cobra"
)

// persistCmd represents the persist command
var persistCmd = &cobra.Command{
	Use:   "persist",
	Short: "generate local yaml config",
	Long: `manage vars into file`,
	Args:cobra.OnlyValidArgs,
	ValidArgs:[]string{"config"},
	Run: func(cmd *cobra.Command, args []string) {
		handler.GenConfigYaml(&cmdRequest)
	},
}

func init() {
	rootCmd.AddCommand(persistCmd)
	persistCmd.Flags().BoolP("force-over","f",false, "force over, if persist file exist")
	flagBindviper(persistCmd, false,"force-over","force-over")
}
