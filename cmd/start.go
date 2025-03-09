package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/waynekn/tidytables/dbconn"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Tidy Tables",
	Long: `The start command initiates Tidy Tables and establishes a connection to the database.
This command requires specific flags to provide the necessary database connection details`,
	Run: func(cmd *cobra.Command, args []string) {

		port := getFlagValue(cmd, "port")
		user := getFlagValue(cmd, "user")
		password := getFlagValue(cmd, "password")
		dbName := getFlagValue(cmd, "name")
		host := getFlagValue(cmd, "host")

		db, err := dbconn.ConnectPostgres(host, port, user, password, dbName)

		if err != nil {
			log.SetFlags(0)
			log.Fatal(color.RedString(err.Error()))
		} else {
			log.SetFlags(0)
			log.Printf(color.GreenString("successfully connected to %v database"), dbName)

		}

		defer db.Close()

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().String("port", "5432", "port which the database is running on")

	startCmd.Flags().StringP("host", "H", "localhost", "the host of the database")

	startCmd.Flags().StringP("user", "U", "", "user to connect to the database as")
	startCmd.MarkFlagRequired("user")

	startCmd.Flags().StringP("password", "P", "", "database password")
	startCmd.MarkFlagRequired("password")

	startCmd.Flags().StringP("name", "N", "", "name of the database to connect to")
	startCmd.MarkFlagRequired("name")

	startCmd.Flags().String("db", "postgres", "name of the relational database management system")
	startCmd.MarkFlagRequired("db")
}

func getFlagValue(cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		logUnexpectedError()
	}
	return value
}

func logUnexpectedError() {
	log.Fatal(color.RedString("An unexpected error occured while getting flags."))
}
