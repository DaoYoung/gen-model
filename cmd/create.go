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
	createCmd.Flags().BoolVarP(&CmdRequest.Gen.IsLowerCamelCaseJson, "isLowerCamelCaseJson", "i", true, "set IsLowerCamelCaseJson true/false")
	flagBindviper(createCmd, false, "searchTableName", "gen.searchTableName")
	flagBindviper(createCmd, false, "isLowerCamelCaseJson", "gen.isLowerCamelCaseJson")
	flagBindviper(createCmd, false, "modelSuffix", "gen.modelSuffix")
	createCmd.Flags().StringVarP(&CmdRequest.Gen.SourceType, "sourceType", "r", "self-table", "self-table: create struct by self table \nlocal-mapper: create struct by local mapper \ngen-table: create struct by stable \"gen_model_mapper\" table")
	createCmd.Flags().StringVarP(&CmdRequest.Gen.PersistType, "persistType", "y", "", "local-mapper: generate local struct mappers \ngen-table: generate mapper table \"gen_model_mapper\" ")
	flagBindviper(createCmd, false, "sourceType", "gen.sourceType")
	flagBindviper(createCmd, false, "persistType", "gen.persistType")
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
		return errors.New("tableName is empty")
	}
	if CmdRequest.Gen.OutPutPath == "" {
		return errors.New("outPutPath is empty")
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
