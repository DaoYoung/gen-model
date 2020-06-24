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
	"fmt"
	"github.com/DaoYoung/gen-model/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gen-model",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// 	Run: func(cmd *cobra.Command, args []string) { },
}
var cmdRequest handler.CmdRequest

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	dir, _ := os.Getwd()
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c","", "config file (default is "+dir+"/"+handler.YamlFile+".yaml)")
	rootCmd.PersistentFlags().StringVarP(&cmdRequest.Db.Host, "host","t", "localhost", "set DB host")
	rootCmd.PersistentFlags().StringVarP(&cmdRequest.Db.Database, "database", "d", "", "set your database")
	rootCmd.PersistentFlags().IntVarP(&cmdRequest.Db.Port, "port", "p", 3306, "set DB port")
	rootCmd.PersistentFlags().StringVarP(&cmdRequest.Db.Username, "username", "u", "root", "set DB login username")
	rootCmd.PersistentFlags().StringVarP(&cmdRequest.Db.Password, "password", "w", "", "set DB login password")
	rootCmd.PersistentFlags().StringVarP(&cmdRequest.Gen.OutPutPath, "outPutPath", "o", ".", "set your OutPutPath")
	rootCmd.PersistentFlags().BoolVarP(&cmdRequest.Gen.HasGormTag, "hasGormTag", "g", true, "gorm tag")
	rootCmd.PersistentFlags().BoolVarP(&cmdRequest.Gen.HasJsonTag, "hasJsonTag", "j", true, "gorm tag")
	rootCmd.PersistentFlags().BoolVarP(&cmdRequest.Gen.HasGureguNullPackage, "hasGureguNullPackage", "n", true, "have package: \"gopkg.in/guregu/null.v3\"")
	rootCmd.PersistentFlags().BoolP("force-cover","f",false, "force over, if persist file exist")
	flagBindviper(rootCmd, true,"force-cover","force-cover")
	flagBindviper(rootCmd, true,"host","mysql.host")
	flagBindviper(rootCmd, true,"database","mysql.database")
	flagBindviper(rootCmd, true,"port","mysql.port")
	flagBindviper(rootCmd, true,"username","mysql.username")
	flagBindviper(rootCmd, true,"password","mysql.password")
	flagBindviper(rootCmd, true,"outPutPath","gen.outPutPath")
	flagBindviper(rootCmd, true,"hasGormTag","gen.hasGormTag")
	flagBindviper(rootCmd, true,"hasJsonTag","gen.hasJsonTag")
	flagBindviper(rootCmd, true,"hasGureguNullPackage","gen.hasGureguNullPackage")
	handler.Welcome()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := homedir.Dir()
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gen-model" (without extension).
		viper.AddConfigPath(dir)
		viper.SetConfigName(handler.YamlFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func flagBindviper(cmd *cobra.Command, isPersistentFlag bool, flagKey, viperKey string) {
	var err error
	if isPersistentFlag {
		err = viper.BindPFlag(viperKey, cmd.PersistentFlags().Lookup(flagKey))
	}else{
		err = viper.BindPFlag(viperKey, cmd.Flags().Lookup(flagKey))
	}
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}
}


