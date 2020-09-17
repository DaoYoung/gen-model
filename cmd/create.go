package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create model struct",
	Long:  `with mysql connect, generate model file`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		generateModel()

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&CmdRequest.Gen.SearchTableName, "searchTableName", "s", "", "set your searchTableName, support patten with '*'")
	createCmd.Flags().StringVarP(&CmdRequest.Gen.ModelSuffix, "modelSuffix", "m", "", "model suffix")
	createCmd.Flags().BoolVarP(&CmdRequest.Gen.JSONUcFirst, "jsonUcFirst", "j", true, "true/false")
	createCmd.Flags().StringVarP(&CmdRequest.Gen.Source, "source", "r", "self-table", "self-table: create struct by self table \nlocal-mapper: create struct by local mapper \ndb-mapper: create struct by stable \"gen_model_mapper\" table")
	createCmd.Flags().StringVarP(&CmdRequest.Gen.Persist, "persist", "y", "", "local-mapper: save mappers at local files \ndb-mapper: create db table: gen_model.struct_mappers, and save mappers in it ")
	flagBindviper(createCmd, false, "searchTableName", "gen.searchTableName")
	flagBindviper(createCmd, false, "jsonUcFirst", "gen.jsonUcFirst")
	flagBindviper(createCmd, false, "modelSuffix", "gen.modelSuffix")
	flagBindviper(createCmd, false, "source", "gen.source")
	flagBindviper(createCmd, false, "persist", "gen.persist")

	createCmd.Flags().Bool("debug", false, "true: print full message")
	flagBindviper(createCmd, false, "debug", "debug")
}

func validArgs() error {
	if viper.GetString("mysql.host") == "" {
		return errors.New("mysql.host is empty")
	}
	if viper.GetString("mysql.database") == "" {
		return errors.New("mysql.database is empty")
	}
	if viper.GetString("mysql.username") == "" {
		return errors.New("mysql.username is empty")
	}
	if viper.GetString("mysql.password") == "" {
		return errors.New("mysql.password is empty")
	}
	if CmdRequest.Gen.SearchTableName == "" {
		return errors.New("searchTableName is empty")
	}
	if CmdRequest.Gen.OutDir == "" {
		return errors.New("outDir is empty")
	}
	return nil
}

func generateModel() {
	CmdRequest.SetDataByViper()
	if viper.GetBool("debug") {
		log.Printf("%+v", CmdRequest)
	}
	if err := validArgs(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	CmdRequest.CreateModelStruct()
}
