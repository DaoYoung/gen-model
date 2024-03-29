package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/victorien-a/gen-model/handler"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gen-model",
	Short: "generate model struct",
	Long:  "Features:\n  1. generate model struct\n  2. filter table columns with persistent mappers",
}

// CmdRequest request arguments manager
var CmdRequest handler.CmdRequest

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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is "+filepath.Join(dir, handler.YamlFile+handler.YamlExt)+")")
	rootCmd.PersistentFlags().BoolP("forceCover", "f", false, "force over, if persist file exist")

	rootCmd.PersistentFlags().StringVarP(&CmdRequest.Db.Host, "host", "t", "localhost", "set DB host")
	rootCmd.PersistentFlags().StringVarP(&CmdRequest.Db.Database, "database", "d", "", "set your database")
	rootCmd.PersistentFlags().IntVarP(&CmdRequest.Db.Port, "port", "p", 3306, "set DB port")
	rootCmd.PersistentFlags().StringVarP(&CmdRequest.Db.Username, "username", "u", "root", "set DB login username")
	rootCmd.PersistentFlags().BoolVarP(&CmdRequest.Db.UsePassword, "usePassword", "x", true, "set to use password")
	rootCmd.PersistentFlags().StringVarP(&CmdRequest.Db.Password, "password", "w", "", "set DB login password")
	rootCmd.PersistentFlags().StringVarP(&CmdRequest.Gen.OutDir, "outDir", "o", "./model/", "set your OutDir")
	rootCmd.PersistentFlags().BoolVarP(&CmdRequest.Gen.DumpAllTables, "dumpAllTables", "a", false, "set to dump all tables found in DB")

	flagBindviper(rootCmd, true, "forceCover", "forceCover")
	flagBindviper(rootCmd, true, "host", "mysql.host")
	flagBindviper(rootCmd, true, "database", "mysql.database")
	flagBindviper(rootCmd, true, "port", "mysql.port")
	flagBindviper(rootCmd, true, "username", "mysql.username")
	flagBindviper(rootCmd, true, "password", "mysql.password")
	flagBindviper(rootCmd, true, "usePassword", "mysql.usePassword")
	flagBindviper(rootCmd, true, "outDir", "gen.outDir")
	flagBindviper(rootCmd, true, "dumpAllTables", "gen.dumpAllTables")

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
	} else {
		err = viper.BindPFlag(viperKey, cmd.Flags().Lookup(flagKey))
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
